package model

import (
	"time"
)

type Users struct {
	UserID   string    `json:"user_id" db:"user_id"`
	GoogleID string    `json:"google_id" db:"google_id"`
	Email    string    `json:"email" db:"email"`
	Name     string   `json:"name" db:"name"`
	Surname  string    `json:"surname" db:"surname"`
	Username string    `json:"username" db:"username"`
	Password string    `json:"password,omitempty" db:"password"`
	Role     int       `json:"role" db:"role"`
	CreateAt time.Time `json:"create_at" db:"create_at"`
	UpdateAt time.Time `json:"update_at" db:"update_at"`
	DeleteAt time.Time `db:"delete_at"`
}

type UserRepository interface {
	CreateUser(newUser Users) error
}

type UserUsecases interface {
	CreateUser(newUser Users) error
}
