package repository

import (
	"errors"
	"time"
	"zuck-my-clothe/zuck-my-clothe-backend/model"
	"zuck-my-clothe/zuck-my-clothe-backend/platform"

	"gorm.io/gorm"
)

type UserRepository interface {
	CreateUser(newUser model.Users) (*model.Users, error)
	FindUserByUserID(userID string) (*model.Users, error)
	FindUserByEmail(email string) (*model.Users, error)
	FindUserByGoogleID(googleID string) (*model.Users, error)
	GetAll() ([]model.Users, error)
	GetAllManager() ([]model.Users, error)
	GetUserByBranchID(branchID string) ([]model.Users, error)
	DeleteUser(userID string) (*model.Users, error)
	UndeleteUser(newUser model.Users) (int64, error)
	UpdateUser(newUser model.Users) error
}

type userRepository struct {
	db *platform.Postgres
}

func CreatenewUserRepository(db *platform.Postgres) UserRepository {
	return &userRepository{db: db}
}

func (repo *userRepository) CreateUser(newUser model.Users) (*model.Users, error) {
	if err := repo.db.Create(&newUser).Error; err != nil {
		return nil, err
	}
	return &newUser, nil
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

func (repo *userRepository) GetUserByBranchID(branchID string) ([]model.Users, error) {
	users := make([]model.Users, 0)
	// dbTx := repo.db.Find(&users)

	result := repo.db.Raw(`
	SELECT DISTINCT U.*
	FROM "Users" U
	LEFT JOIN "OrderHeaders" OH ON U.user_id = OH.user_id
	WHERE OH.branch_id = $1
	`, branchID).Scan(&users)

	if result.Error != nil {
		return nil, result.Error
	}

	for i := range users {
		users[i].Password = ""
	}

	return users, nil
}

func (repo *userRepository) GetAllManager() ([]model.Users, error) {
	var users []model.Users
	dbTx := repo.db.Find(&users, "role = ?", model.BranchManager)
	if dbTx.Error != nil {
		return nil, dbTx.Error
	}

	for i := range users {
		users[i].Password = ""
	}
	return users, nil
}

func (repo *userRepository) DeleteUser(userID string) (*model.Users, error) {
	deletedUser := new(model.Users)
	returnValue := repo.db.Model(&model.Users{}).Where("user_id = ?", userID).Delete(deletedUser)
	if returnValue.RowsAffected == 0 {
		return nil, gorm.ErrRecordNotFound
	}
	returnValue = repo.db.Unscoped().Where("user_id = ?", userID).First(deletedUser)
	return deletedUser, returnValue.Error
}

func (repo *userRepository) UndeleteUser(newUser model.Users) (int64, error) {
	deletedUser := new(model.Users)
	dbResponse := repo.db.Unscoped().First(deletedUser, "email = ?", newUser.Email)
	if dbResponse.Error != nil {
		return 0, dbResponse.Error
	}
	if !deletedUser.DeleteAt.Valid {
		return 0, errors.New("cannot undelete user")
	}
	deletedUser.DeleteAt.Valid = false
	deletedUser.UpdateAt = time.Now().UTC()
	dbResponse = repo.db.Table("Users").Unscoped().Where("user_id = ?", deletedUser.UserID).Updates(&model.Users{Phone: newUser.Phone,
		FirstName:       newUser.FirstName,
		LastName:        newUser.LastName,
		Password:        newUser.Password,
		ProfileImageURL: newUser.ProfileImageURL,
		Role:            newUser.Role,
		UpdateAt:        deletedUser.UpdateAt,
		DeleteAt:        deletedUser.DeleteAt})

	return dbResponse.RowsAffected, dbResponse.Error
}

func (repo *userRepository) UpdateUser(newUser model.Users) error {
	returnValue := repo.db.Model(&model.Users{}).Where("user_id = ?", newUser.UserID).Updates(newUser)
	return returnValue.Error
}
