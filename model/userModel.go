package model

import (
	"time"
)

type Roles string

func (Users) TableName() string {
	return "Users"
}

const (
	SuperAdmin    Roles = "SuperAdmin"
	BranchManager Roles = "BranchManager"
	Employee      Roles = "Employee"
	Client        Roles = "Client"
)

type Users struct {
	UserID          string    `json:"user_id" gorm:"column:user_id"`
	GoogleID        string    `json:"google_id" gorm:"column:google_id"`
	Email           string    `json:"email" gorm:"column:email"`
	Phone           string    `json:"phone" gorm:"column:phone"`
	FirstName       string    `json:"firstname" gorm:"column:firstname"`
	LastName        string    `json:"lastname" gorm:"column:lastname"`
	Password        string    `json:"password,omitempty" gorm:"column:password"`
	ProfileImageURL string    `json:"profile_image_url" gorm:"column:profile_image_url"`
	Role            Roles     `json:"role" gorm:"column:role"`
	CreateAt        time.Time `json:"created_at" gorm:"column:created_at"`
	UpdateAt        time.Time `json:"updated_at" gorm:"column:updated_at"`
	DeleteAt        time.Time `gorm:"column:deleted_at"`
}

type GoogleUser struct {
	GoogleID    string `json:"id" db:"google_id"`
	Email       string `json:"email" db:"email"`
	VerifyEmail bool   `json:"verified_email"`
	FullName    string `json:"name"`
	Name        string `json:"given_name" db:"name"`
	Surname     string `json:"family_name" db:"surname"`
	ImageUrl    string `json:"picture"`
}

type UserRepository interface {
	CreateUser(newUser Users) error
	FindUserByUserID(userID string) (*Users, error)
	FindUserByEmail(email string) (*Users, error)
	FindUserByGoogleID(googleID string) (*Users, error)
	//FindUserByUUID(uuid string) Users
}

type UserUsecases interface {
	CreateUser(newUser Users) error
	FindUserByEmail(email string) (*Users, error)
	FindUserByGoogleID(googleID string) (*Users, error)
}
