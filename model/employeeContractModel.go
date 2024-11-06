package model

import (
	"time"

	"gorm.io/gorm"
)

func (EmployeeContract) TableName() string {
	return "EmployeeContracts"
}

// create enum for position
type EmployeeContractPosition string

const (
	Worker  EmployeeContractPosition = "Worker"
	Deliver EmployeeContractPosition = "Deliver"
)

type EmployeeContract struct {
	ContractId string         `json:"contract_id" gorm:"column:contract_id;primaryKey"`
	UserID     string         `json:"user_id" gorm:"column:user_id"`
	BranchID   string         `json:"branch_id" gorm:"column:branch_id"`
	PositionId string         `json:"position_id" gorm:"column:position_id"`
	CreatedAt  time.Time      `json:"created_at" gorm:"column:created_at"`
	CreatedBy  string         `json:"created_by" gorm:"column:created_by"`
	DeletedAt  gorm.DeletedAt `json:"deleted_at" gorm:"column:deleted_at" swaggertype:"string" example:"null"`
	DeletedBy  *string        `json:"deleted_by" gorm:"column:deleted_by"`
}

type EmployeeContractDTO struct {
	UserID     string `json:"user_id"  validate:"required"`
	BranchID   string `json:"branch_id"  validate:"required"`
	PositionId string `json:"position_id"  validate:"required,employeeContractPosition"`
	CreatedBy  string `json:"created_by" validate:"required"`
}
