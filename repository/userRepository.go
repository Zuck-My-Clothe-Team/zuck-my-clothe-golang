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

func (repo *userRepository) FindUserByUserID(userID string) (*model.Users, error) {
	result := new(model.Users)
	dbTx := repo.db.First(result, "user_id = ?", userID)
	if dbTx.Error != nil {
		return nil, dbTx.Error
	}
	return result, nil
}

func (repo *userRepository) FindUserByEmail(email string) (*model.Users, error) {
	result := new(model.Users)
	dbTx := repo.db.First(result, "email = ?", email)
	if dbTx.Error != nil {
		return nil, dbTx.Error
	}
	return result, nil
}

func (repo *userRepository) FindUserByGoogleID(googleID string) (*model.Users, error) {
	result := new(model.Users)
	dbTx := repo.db.First(result, "google_id = ?", googleID)
	if dbTx.Error != nil {
		return nil, dbTx.Error
	}
	return result, nil
}

func (repo *userRepository) GetAll() ([]model.Users, error) {
	var users []model.Users
	dbTx := repo.db.Find(&users)
	if dbTx.Error != nil {
		return nil, dbTx.Error
	}

	for i := range users {
		users[i].Password = ""
	}

	return users, nil
}
