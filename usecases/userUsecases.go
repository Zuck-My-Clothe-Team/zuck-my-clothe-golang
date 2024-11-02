package usecases

import (
	"errors"
	"zuck-my-clothe/zuck-my-clothe-backend/model"
	"zuck-my-clothe/zuck-my-clothe-backend/utils"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type userUsecases struct {
	repository model.UserRepository
}

func CreateNewUserUsecases(repository model.UserRepository) model.UserUsecases {
	return &userUsecases{repository: repository}
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

func (repo *userUsecases) DeleteUser(userID string) error {
	err := repo.repository.DeleteUser(userID)
	return err
}
