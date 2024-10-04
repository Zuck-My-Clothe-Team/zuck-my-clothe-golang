package usecases

import (
	"errors"
	"zuck-my-clothe/zuck-my-clothe-backend/model"
	"zuck-my-clothe/zuck-my-clothe-backend/utils"

	"golang.org/x/crypto/bcrypt"
)

type authenUsecase struct {
	authenRepository model.AuthenRepository
	userRepository   model.UserRepository
}

func CreateNewAuthenUsecase(authenRepository model.AuthenRepository, userRepository model.UserRepository) model.AuthenRepository {
	return &authenUsecase{authenRepository: authenRepository,
		userRepository: userRepository}
}

func (s *authenUsecase) SignIn(user *model.AuthenPayload) (*model.AuthenPayload, error) {
	//Need to validate email
	if utils.CheckStraoPling(user.Email) ||
		utils.CheckStraoPling(user.Password) {
		return nil, errors.New("null detected on one or more essential field(s)")
	}
	dbResult, err := s.userRepository.FindUserByEmail(user.Email)
	if err != nil {
		return nil, err
	}
	//fmt.Println(dbResult)
	payLoad := new(model.AuthenPayload)
	if errBcrypt := bcrypt.CompareHashAndPassword([]byte(dbResult.Password), []byte(user.Password)); errBcrypt != nil {
		return nil, errors.New("password not match")
	}
	payLoad.UserId = dbResult.UserID
	payLoad.Email = dbResult.Email
	payLoad.Role = dbResult.Role
	payLoad.Password = ""
	return payLoad, nil
}