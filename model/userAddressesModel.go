package model

import (
	"time"

	"gorm.io/gorm"
)

func (UserAddresses) TableName() string {
	return "UserAddresses"
}

type UserAddresses struct {
	AddressID   string         `json:"address_id" gorm:"column:address_id;primaryKey"`
	UserID      string         `json:"user_id" gorm:"column:user_id"`
	Address     string         `json:"address" gorm:"column:address"`
	Province    string         `json:"province" gorm:"column:province"`
	District    string         `json:"district" gorm:"column:district"`
	SubDistrict string         `json:"subdistrict" gorm:"column:subdistrict"`
	Zipcode     string         `json:"zipcode" gorm:"column:zipcode"`
	Lat         float64        `json:"lat" gorm:"column:lat"`
	Long        float64        `json:"long" gorm:"column:long"`
	CreatedAt   time.Time      `json:"created_at" gorm:"column:created_at"`
	UpdatedAt   time.Time      `json:"updated_at" gorm:"column:updated_at"`
	DeletedAt   gorm.DeletedAt `json:"deleted_at" gorm:"column:deleted_at;index" swaggertype:"string" example:"null"`
}

type AddUserAddressDTO struct {
	AddressID   string  `json:"address_id" gorm:"column:address_id;primaryKey"`
	Address     string  `json:"address" validate:"required"`
	Province    string  `json:"province" validate:"required"`
	District    string  `json:"district" validate:"required"`
	SubDistrict string  `json:"subdistrict" validate:"required"`
	Zipcode     string  `json:"zipcode" validate:"required"`
	Lat         float64 `json:"lat" validate:"required"`
	Long        float64 `json:"long" validate:"required"`
}

type UpdateUserAddressDTO struct {
	Address     string    `gorm:"column:address"`
	Province    string    `gorm:"column:province"`
	District    string    `gorm:"column:district"`
	SubDistrict string    `gorm:"column:subdistrict"`
	Zipcode     string    `gorm:"column:zipcode"`
	Lat         float64   `gorm:"column:lat"`
	Long        float64   `gorm:"column:long"`
	UpdatedAt   time.Time `gorm:"column:updated_at"`
}

type UserAddressDetail struct {
	AddressID   string  `json:"address_id"`
	UserID      string  `json:"user_id"`
	Address     string  `json:"address"`
	Province    string  `json:"province"`
	District    string  `json:"district"`
	SubDistrict string  `json:"subdistrict"`
	Zipcode     string  `json:"zipcode"`
	Lat         float64 `json:"lat"`
	Long        float64 `json:"long"`
}

type UserAddressesRepository interface {
	AddUserAddress(newUserAddress *UserAddresses) error
	FindUserAddressByID(addressID string, userID string) (*UserAddresses, error)
	FindUserAddresByOwnerID(ownerID string) (*[]UserAddresses, error)
	UpdateUserAddressData(userID string, addressID string, updatedAddressData *UpdateUserAddressDTO) error
	DeleteUserAddress(userID string, addressID string) error
}

type UserAddressesUsecase interface {
	AddUserAddress(owenrID string, newUserAddress *AddUserAddressDTO) (*interface{}, error)
	FindUserAddressByID(addressID string, userID string) (*interface{}, error)
	FindUserAddresByOwnerID(ownerID string) (*[]interface{}, error)
	UpdateUserAddressData(userID string, updatedAddressData *AddUserAddressDTO) (*interface{}, error)
	DeleteUserAddress(userID string, addressID string) error
}
