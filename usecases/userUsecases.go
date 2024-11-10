package usecases

import (
	"errors"
	"fmt"
	"time"
	"zuck-my-clothe/zuck-my-clothe-backend/model"
	"zuck-my-clothe/zuck-my-clothe-backend/repository"
	"zuck-my-clothe/zuck-my-clothe-backend/utils"
	validatorboi "zuck-my-clothe/zuck-my-clothe-backend/validator"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type UserUsecases interface {
	CreateUser(newUser model.Users) error
	FindUserByEmail(email string) (*model.Users, error)
	FindUserByGoogleID(googleID string) (*model.Users, error)
	GetAll() ([]model.UserBranch, error)
	GetUserById(userID string) (*model.Users, error)
	GetBranchEmployee(branchId string) ([]model.UserContract, error)
	GetAllManager() ([]model.Users, error)
	DeleteUser(userID string) (*model.Users, error)
	UpdateUser(userID string, newUser model.UserUpdateDTO, role string) error
	UpdateUserPassword(userID string, newUser model.UserUpdatePasswordDTO) error
}

type userUsecases struct {
	repository                 repository.UserRepository
	employeeContractRepository repository.EmployeeContractRepository
	branchRepository           repository.BranchReopository
}

func CreateNewUserUsecases(
	repository repository.UserRepository,
	employeeContractRepository repository.EmployeeContractRepository,
	branchRepository repository.BranchReopository,
) UserUsecases {
	return &userUsecases{
		repository:                 repository,
		employeeContractRepository: employeeContractRepository,
		branchRepository:           branchRepository,
	}
}

func (repo *userUsecases) CreateUser(newUser model.Users) error {
	newUser.UserID = uuid.New().String()
	if utils.CheckStraoPling(newUser.FirstName) ||
		utils.CheckStraoPling(newUser.Email) ||
		(utils.CheckStraoPling(newUser.Password) == utils.CheckStraoPling(newUser.GoogleID)) {
		return errors.New("null detected on one or more essential field(s)")
	}
	buffer, hashErr := bcrypt.GenerateFromPassword([]byte(newUser.Password), bcrypt.DefaultCost)
	if hashErr != nil {
		return hashErr
	}

	newUser.Password = string(buffer)
	numRow, err := repo.repository.UndeleteUser(newUser)
	if numRow == 1 && err == nil {
		return nil
	}
	newUser.CreateAt = time.Now()
	newUser.UpdateAt = time.Now()
	repoError := repo.repository.CreateUser(newUser)
	return repoError
}

func (repo *userUsecases) FindUserByEmail(email string) (*model.Users, error) {
	user, err := repo.repository.FindUserByEmail(email)
	return user, err
}

func (repo *userUsecases) FindUserByGoogleID(googleID string) (*model.Users, error) {
	user, err := repo.repository.FindUserByGoogleID(googleID)
	return user, err
}

func (repo *userUsecases) GetAll() ([]model.UserBranch, error) {
	users, err := repo.repository.GetAll()

	var allUsers []model.UserBranch

	for _, user := range users {
		var userBranch model.UserBranch
		if user.Role == model.BranchManager {
			branch, err := repo.branchRepository.GetByBranchOwner(user.UserID)
			if err != nil {
				if err.Error() == "record not found" {
					branch = &[]model.Branch{}
				}
			}

			userBranch = model.UserBranch{
				UserID:          user.UserID,
				Email:           user.Email,
				Phone:           user.Phone,
				FirstName:       user.FirstName,
				LastName:        user.LastName,
				ProfileImageURL: user.ProfileImageURL,
				Role:            user.Role,
				Branch:          *branch,
			}

		} else if user.Role == model.Employee {
			contract, err := repo.employeeContractRepository.GetByUserID(user.UserID)
			if err != nil {
				if err.Error() == "record not found" {
					contract = &[]model.EmployeeContract{}
				} else {
					return nil, err
				}
			}

			var branch *model.Branch = nil

			if len(*contract) != 0 {
				branch, err = repo.branchRepository.GetByBranchID((*contract)[0].BranchID)
				if err != nil {
					if err.Error() != "record not found" {
						return nil, err
					}
					branch = nil
				}
			}

			branches := []model.Branch{}
			if branch != nil {
				branches = append(branches, *branch)
			}

			userBranch = model.UserBranch{
				UserID:          user.UserID,
				Email:           user.Email,
				Phone:           user.Phone,
				FirstName:       user.FirstName,
				LastName:        user.LastName,
				ProfileImageURL: user.ProfileImageURL,
				Role:            user.Role,
				Branch:          branches,
			}
		} else {
			userBranch = model.UserBranch{
				UserID:          user.UserID,
				Email:           user.Email,
				Phone:           user.Phone,
				FirstName:       user.FirstName,
				LastName:        user.LastName,
				ProfileImageURL: user.ProfileImageURL,
				Role:            user.Role,
				Branch:          []model.Branch{},
			}
		}
		allUsers = append(allUsers, userBranch)
	}

	return allUsers, err
}

func (repo *userUsecases) GetUserById(userID string) (*model.Users, error) {
	user, err := repo.repository.FindUserByUserID(userID)
	return user, err
}

func (repo *userUsecases) GetAllManager() ([]model.Users, error) {
	users, err := repo.repository.GetAllManager()
	return users, err
}

func (repo *userUsecases) DeleteUser(userID string) (*model.Users, error) {
	deletedUser, err := repo.repository.DeleteUser(userID)
	return deletedUser, err
}

func (repo *userUsecases) GetBranchEmployee(branchId string) ([]model.UserContract, error) {
	employeeContracts, err := repo.employeeContractRepository.GetByBranchID(branchId)
	if err != nil {
		return nil, err
	}

	userContractMap := make(map[string]model.UserContract)
	for _, contract := range *employeeContracts {
		user, err := repo.repository.FindUserByUserID(contract.UserID)
		if err != nil {
			return nil, err
		}

		if existingContract, exists := userContractMap[user.UserID]; exists {
			existingContract.Contracts = append(existingContract.Contracts, contract)
			userContractMap[user.UserID] = existingContract
		} else {
			userContractMap[user.UserID] = model.UserContract{
				UserID:          user.UserID,
				Email:           user.Email,
				Phone:           user.Phone,
				FirstName:       user.FirstName,
				LastName:        user.LastName,
				ProfileImageURL: user.ProfileImageURL,
				Role:            user.Role,
				Contracts:       []model.EmployeeContract{contract},
			}
		}
	}

	var userContracts []model.UserContract
	for _, userContract := range userContractMap {
		userContracts = append(userContracts, userContract)
	}

	return userContracts, nil
}

func (repo *userUsecases) UpdateUser(userID string, body model.UserUpdateDTO, role string) error {

	if err := validatorboi.Validate(body); err != nil {
		return err
	}
	existingUser, err := repo.repository.FindUserByUserID(userID)

	var newRole model.Roles

	if err != nil {
		return err
	}

	fmt.Println(existingUser.Role + " " + body.Role)
	if body.Role != existingUser.Role && (role == string(model.Employee) || role == string(model.BranchManager) || role == string(model.Client)) {
		return fmt.Errorf("role cannot be changed")
	} else {
		newRole = body.Role
	}

	updatedUser := model.Users{
		UserID:          existingUser.UserID,
		Email:           existingUser.Email,
		Phone:           body.Phone,
		FirstName:       body.FirstName,
		LastName:        body.LastName,
		ProfileImageURL: existingUser.ProfileImageURL,
		Role:            newRole,
		Password:        existingUser.Password,
		CreateAt:        existingUser.CreateAt,
		UpdateAt:        time.Now(),
	}

	err = repo.repository.UpdateUser(updatedUser)

	if err != nil {
		return err
	}

	return nil
}

func (repo *userUsecases) UpdateUserPassword(userID string, body model.UserUpdatePasswordDTO) error {

	if err := validatorboi.Validate(body); err != nil {
		return err
	}
	existingUser, err := repo.repository.FindUserByUserID(userID)

	if err != nil {
		return err
	}

	newPassword, hashErr := bcrypt.GenerateFromPassword([]byte(body.Password), bcrypt.DefaultCost)

	if hashErr != nil {
		return hashErr
	}

	updatedUser := model.Users{
		UserID:          existingUser.UserID,
		Email:           existingUser.Email,
		Phone:           existingUser.Phone,
		FirstName:       existingUser.FirstName,
		LastName:        existingUser.LastName,
		ProfileImageURL: existingUser.ProfileImageURL,
		Role:            existingUser.Role,
		Password:        string(newPassword),
		CreateAt:        existingUser.CreateAt,
		UpdateAt:        time.Now(),
	}

	err = repo.repository.UpdateUser(updatedUser)

	if err != nil {
		return err
	}

	return nil
}
