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

//	@Summary		Get machine details
//	@Description	Retrieve details of a machine
//	@Tags			Machine
//	@Produce		json
//	@Success		200	{object}	MachineResponse
//	@Router			/machine/{id} [get]
type Machine struct {
	MachineSerial string         `json:"machine_serial" gorm:"column:machine_serial;primaryKey"`
	BranchID      string         `json:"branch_id" gorm:"column:branch_id"`
	MachineType   MachineType    `json:"machine_type" gorm:"column:machine_type"`
	IsActive      bool           `json:"is_active" gorm:"column:is_active"`
	Weight        int16          `json:"weight" gorm:"column:weight"`
	CreatedAt     time.Time      `json:"created_at" gorm:"column:created_at"`
	CreatedBy     *string        `json:"created_by" gorm:"column:created_by"`
	UpdatedAt     time.Time      `json:"updated_at" gorm:"column:updated_at"`
	UpdatedBy     *string        `json:"updated_by" gorm:"column:updated_by"`
	DeletedAt     gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"column:deleted_at"`
	DeletedBy     *string        `json:"deleted_by" gorm:"column:deleted_by"`
}

type AddMachineDTO struct {
	MachineSerial string      `json:"machine_serial" validate:"required"`
	BranchID      string      `json:"branch_id" validate:"required"`
	MachineType   MachineType `json:"machine_type" validate:"required,machineType"`
	Weight        int         `json:"weight" validate:"required,gte=0"`
	CreatedBy     string
}

type MachineDetail struct {
	MachineSerial string      `json:"machine_serial" gorm:"column:machine_serial"`
	BranchID      string      `json:"branch_id" gorm:"column:branch_id"`
	MachineType   MachineType `json:"machine_type" gorm:"column:machine_type"`
	IsActive      bool        `json:"is_active" gorm:"column:is_active"`
	Weight        int16       `json:"weight" gorm:"column:weight"`
}

type MachineRepository interface {
	SoftDelete(machine_serial string, deleted_by string) (*Machine, error)
	UpdateActive(branch_id string, set_active bool, updated_by string) (*Machine, error)
	GetByBranchID(branch_id string) (*[]Machine, error)
	GetAll() (*[]Machine, error)
	AddMachine(newMachine *Machine) error
	GetByMachineSerial(machineSerial string) (*Machine, error)
}

type MachineUsecase interface {
	SoftDelete(machine_serial string, deleted_by string) (*MachineDetail, error)
	UpdateActive(branch_id string, set_active bool, updated_by string) (*MachineDetail, error)
	GetByBranchID(branch_id string) (*[]MachineDetail, error)
	GetAll() (*[]MachineDetail, error)
	AddMachine(newMachine *AddMachineDTO) (*MachineDetail, error)
	GetByMachineSerial(machineSerial string) (*MachineDetail, error)
}
