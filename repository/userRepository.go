package repository

import (
	"zuck-my-clothe/zuck-my-clothe-backend/model"
	"zuck-my-clothe/zuck-my-clothe-backend/platform"
)

type userRepository struct {
	db *platform.Postgres
}

func CreatenewUserRepository(db *platform.Postgres) model.UserRepository {
	return &userRepository{db: db}
}

func (repo *userRepository) CreateUser(newUser model.Users) error {
	returnValue := repo.db.Create(newUser)
	return returnValue.Error
}
