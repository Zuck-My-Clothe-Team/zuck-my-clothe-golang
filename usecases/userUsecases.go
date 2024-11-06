package usecases

import (
	"errors"
	"time"
	"zuck-my-clothe/zuck-my-clothe-backend/model"
	"zuck-my-clothe/zuck-my-clothe-backend/repository"
	"zuck-my-clothe/zuck-my-clothe-backend/utils"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type UserUsecases interface {
	CreateUser(newUser model.Users) error
	FindUserByEmail(email string) (*model.Users, error)
	FindUserByGoogleID(googleID string) (*model.Users, error)
	GetAll() ([]model.Users, error)
	GetBranchEmployee(branchId string) ([]model.UserContract, error)
	GetAllManager() ([]model.Users, error)
	DeleteUser(userID string) (*model.Users, error)
}

type userUsecases struct {
	repository                 repository.UserRepository
	employeeContractRepository repository.EmployeeContractRepository
}

func CreateNewUserUsecases(
	repository repository.UserRepository,
	employeeContractRepository repository.EmployeeContractRepository,
) UserUsecases {
	return &userUsecases{
		repository:                 repository,
		employeeContractRepository: employeeContractRepository,
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

func (repo *userUsecases) GetAll() ([]model.Users, error) {
	users, err := repo.repository.GetAll()
	return users, err
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
