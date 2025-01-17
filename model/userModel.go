package model

import (
	"time"

	"gorm.io/gorm"
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
	UserID          string         `json:"user_id" gorm:"column:user_id"`
	GoogleID        string         `json:"google_id" gorm:"column:google_id"`
	Email           string         `json:"email" gorm:"column:email"`
	Phone           string         `json:"phone" gorm:"column:phone"`
	FirstName       string         `json:"firstname" gorm:"column:firstname"`
	LastName        string         `json:"lastname" gorm:"column:lastname"`
	Password        string         `json:"password,omitempty" gorm:"column:password"`
	ProfileImageURL string         `json:"profile_image_url" gorm:"column:profile_image_url"`
	Role            Roles          `json:"role" gorm:"column:role"`
	CreateAt        time.Time      `json:"created_at" gorm:"column:created_at"`
	UpdateAt        time.Time      `json:"updated_at" gorm:"column:updated_at"`
	DeleteAt        gorm.DeletedAt `gorm:"column:deleted_at" swaggertype:"string" example:"null"`
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

type UserDetailDTO struct {
	UserID          string `json:"user_id"`
	GoogleID        string `json:"google_id"`
	Email           string `json:"email"`
	Phone           string `json:"phone"`
	FirstName       string `json:"firstname"`
	LastName        string `json:"lastname"`
	ProfileImageURL string `json:"profile_image_url"`
	Role            Roles  `json:"role"`
}

type UserDTO struct {
	Email           string                `json:"email" validate:"required"`
	Password        string                `json:"password" validate:"required"`
	Phone           string                `json:"phone"`
	FirstName       string                `json:"firstname" validate:"required"`
	LastName        string                `json:"lastname" validate:"required"`
	ProfileImageURL string                `json:"profile_image_url"`
	Role            Roles                 `json:"role" validate:"required,userRoles"`
	Contracts       []EmployeeContractDTO `json:"contracts,omitempty"`
}

type UserUpdateDTO struct {
	Phone     string `json:"phone"`
	FirstName string `json:"firstname" validate:"required"`
	LastName  string `json:"lastname" validate:"required"`
	Role      Roles  `json:"role" validate:"required,userRoles"`
}

type UserUpdatePasswordDTO struct {
	Password string `json:"password" validate:"required"`
}

type UserContract struct {
	UserID          string             `json:"user_id"`
	Email           string             `json:"email"`
	Phone           string             `json:"phone"`
	FirstName       string             `json:"firstname"`
	LastName        string             `json:"lastname"`
	ProfileImageURL string             `json:"profile_image_url"`
	Role            Roles              `json:"role"`
	Contracts       []EmployeeContract `json:"contracts"`
}

type UserBranch struct {
	UserID          string   `json:"user_id"`
	Email           string   `json:"email"`
	Phone           string   `json:"phone"`
	FirstName       string   `json:"firstname"`
	LastName        string   `json:"lastname"`
	ProfileImageURL string   `json:"profile_image_url"`
	Role            Roles    `json:"role"`
	Branch          []Branch `json:"branch"`
}
