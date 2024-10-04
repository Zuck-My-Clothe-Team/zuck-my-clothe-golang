package repository

import (
	"zuck-my-clothe/zuck-my-clothe-backend/model"
	"zuck-my-clothe/zuck-my-clothe-backend/platform"
)

type authenRepo struct {
	db *platform.Postgres
}

func CreateNewAuthenticationRepository(db *platform.Postgres) model.AuthenRepository {
	return &authenRepo{db: db}
}

func (s *authenRepo) SignIn(user *model.AuthenPayload) (*model.AuthenPayload, error) {
	payLoad := new(model.AuthenPayload)
	userStruct := new(model.Users)
	userStruct.Email = user.Email
	userStruct.Password = user.Password
	retVal := s.db.First(userStruct, "email = ?", userStruct.Email)

	if retVal.Error != nil {
		return nil, retVal.Error
	}
	return payLoad, nil
}
