package model

import (
	"time"
)

func (Branch) TableName() string {
	return "Branches"
}

type Branch struct {
	BranchID     string    `json:"branch_id" gorm:"column:branch_id"`
	BranchName   string    `json:"branch_name" gorm:"column:branch_name"`
	BranchDetail string    `json:"branch_detail" gorm:"column:branch_detail"`
	BranchLat    float64   `json:"branch_lat" gorm:"column:branch_lat"`
	BranchLon    float64   `json:"branch_long" gorm:"column:branch_long"`
	OwnerUserID  string    `json:"owner_user_id" gorm:"column:owner_user_id"`
	CreatedAt    time.Time `json:"created_at" gorm:"column:created_at"`
	CreatedBy    string    `json:"created_by" gorm:"column:created_by"`
	UpdatedAt    time.Time `gorm:"column:updated_at"`
	UpdatedBy    string    `json:"updated_by,omitempty" gorm:"column:updated_by"`
	DeletedAt    time.Time `gorm:"column:deleted_at"`
	DeletedBy    string    `json:"deleted_by,omitempty" gorm:"column:deleted_by"`
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
	CreateBranch(newBranch *Branch) (*Branch, error)
	GetAll() (*[]Branch, error)
	GetByBranchID(branchID string) (*Branch, error)
	GetByBranchOwner(ownerUserID string) (*[]Branch, error)
	UpdateBranch(branch *Branch, role string) (*Branch, error)
	DeleteBranch(branch *Branch) error
}
