package model

import (
	"time"

	"gorm.io/gorm"
)

func (Machine) TableName() string {
	return "Machines"
}

type MachineType string

const (
	Washer MachineType = "Washer"
	Dryer  MachineType = "Dryer"
)

type Machine struct {
	MachineSerial string         `json:"machine_serial" gorm:"column:machine_serial;primaryKey"`
	MachineLabel  string         `json:"machine_label" gorm:"column:machine_label"`
	BranchID      string         `json:"branch_id" gorm:"column:branch_id"`
	MachineType   MachineType    `json:"machine_type" gorm:"column:machine_type"`
	IsActive      bool           `json:"is_active" gorm:"column:is_active"`
	Weight        int16          `json:"weight" gorm:"column:weight"`
	CreatedAt     time.Time      `json:"created_at" gorm:"column:created_at"`
	CreatedBy     *string        `json:"created_by" gorm:"column:created_by"`
	UpdatedAt     time.Time      `json:"updated_at" gorm:"column:updated_at"`
	UpdatedBy     *string        `json:"updated_by" gorm:"column:updated_by"`
	DeletedAt     gorm.DeletedAt `json:"deleted_at" gorm:"column:deleted_at;index" swaggertype:"string" example:"null"`
	DeletedBy     *string        `json:"deleted_by" gorm:"column:deleted_by"`
}

type AddMachineDTO struct {
	MachineSerial string      `json:"machine_serial" validate:"required"`
	MachineLabel  int         `json:"machine_label" validate:"required"`
	BranchID      string      `json:"branch_id" validate:"required"`
	MachineType   MachineType `json:"machine_type" validate:"required,machineType"`
	Weight        int         `json:"weight" validate:"required,gte=0"`
	CreatedBy     string
}

type MachineDetail struct {
	MachineSerial string      `json:"machine_serial"`
	MachineLabel  string      `json:"machine_label"`
	BranchID      string      `json:"branch_id"`
	MachineType   MachineType `json:"machine_type"`
	IsActive      bool        `json:"is_active"`
	Weight        int16       `json:"weight"`
}
