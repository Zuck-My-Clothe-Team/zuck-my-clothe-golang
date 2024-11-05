package model

import (
	"time"

	"gorm.io/gorm"
)

func (Branch) TableName() string {
	return "Branches"
}

type Branch struct {
	BranchID     string         `json:"branch_id" gorm:"column:branch_id;primaryKey"`
	BranchName   string         `json:"branch_name" gorm:"column:branch_name"`
	BranchDetail string         `json:"branch_detail" gorm:"column:branch_detail"`
	BranchLat    float64        `json:"branch_lat" gorm:"column:branch_lat"`
	BranchLon    float64        `json:"branch_long" gorm:"column:branch_long"`
	OwnerUserID  string         `json:"owner_user_id" gorm:"column:owner_user_id"`
	CreatedAt    time.Time      `json:"created_at" gorm:"column:created_at"`
	CreatedBy    string         `json:"created_by" gorm:"column:created_by"`
	UpdatedAt    time.Time      `json:"updated_at" gorm:"column:updated_at"`
	UpdatedBy    string         `json:"updated_by" gorm:"column:updated_by"`
	DeletedAt    gorm.DeletedAt `json:"deleted_at" gorm:"column:deleted_at;index" swaggertype:"string" example:"null"`
	DeletedBy    *string        `json:"deleted_by" gorm:"column:deleted_by"`
}

type CreateBranchDTO struct {
	BranchName   string  `json:"branch_name" validate:"required"`
	BranchDetail string  `json:"branch_detail" validate:"required"`
	BranchLat    float64 `json:"branch_lat" validate:"required"`
	BranchLon    float64 `json:"branch_long" validate:"required"`
	OwnerUserID  string  `json:"owner_user_id" validate:"required"`
}
type UpdateBranchDTO struct {
	BranchID     string  `json:"branch_id" validate:"required"`
	BranchName   string  `json:"branch_name" validate:"required"`
	BranchDetail string  `json:"branch_detail" validate:"required"`
	BranchLat    float64 `json:"branch_lat" validate:"required"`
	BranchLon    float64 `json:"branch_long" validate:"required"`
	OwnerUserID  string  `json:"owner_user_id" validate:"required"`
}

type BranchDetail struct {
	BranchID     string  `json:"branch_id"`
	BranchName   string  `json:"branch_name"`
	BranchDetail string  `json:"branch_detail"`
	BranchLat    float64 `json:"branch_lat"`
	BranchLon    float64 `json:"branch_long"`
	OwnerUserID  string  `json:"owner_user_id"`
	Distance     float64 `json:"distance,omitempty"`
}

type UserGeoLocation struct {
	BranchLat float64 `json:"user_lat" validate:"required"`
	BranchLon float64 `json:"user_lon" validate:"required"`
}

type BranchReopository interface {
	CreateBranch(newBranch *Branch) error
	GetAll() (*[]Branch, error)
	GetByBranchID(branchID string) (*Branch, error)
	GetByBranchOwner(ownerUserID string) (*[]Branch, error)
	UpdateBranch(branch *Branch) error
	ManagerUpdateBranch(branch *Branch) error
	DeleteBranch(branch *Branch) error
}

type BranchUsecase interface {
	CreateBranch(newBranch *CreateBranchDTO, userID string) (*Branch, error)
	GetAll(isAdminView bool) (interface{}, error)
	GetClosestToMe(userLocation *UserGeoLocation) (*[]BranchDetail, error)
	GetByBranchID(branchID string, isAdminView bool) (*interface{}, error)
	GetByBranchOwner(ownerUserID string) (*[]Branch, error)
	UpdateBranch(branch *UpdateBranchDTO, role string) (*Branch, error)
	DeleteBranch(branch *Branch) error
}
