package usecases

import (
	_ "fmt"
	"zuck-my-clothe/zuck-my-clothe-backend/model"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	//"zuck-my-clothe/zuck-my-clothe-backend/repository"
)

type userUsecases struct {
	repository model.UserRepository
}

func CreateNewUserUsecases(repository model.UserRepository) model.UserUsecases {
	return &userUsecases{repository: repository}
}

func (repo *userUsecases) CreateUser(newUser model.Users) error {
	newUser.UserID = uuid.New().String()
	buffer, hashErr := bcrypt.GenerateFromPassword([]byte(newUser.Password), bcrypt.DefaultCost)
	if hashErr != nil {
		return hashErr
	}
	newUser.Password = string(buffer)
	repoError := repo.repository.CreateUser(newUser)
	return repoError
}
