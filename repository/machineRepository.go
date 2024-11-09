package repository

import (
	"zuck-my-clothe/zuck-my-clothe-backend/model"
	"zuck-my-clothe/zuck-my-clothe-backend/platform"

	"gorm.io/gorm"
)

type MachineRepository interface {
	SoftDelete(machineSerial string, deletedBy string) (*model.Machine, error)
	UpdateActive(branchId string, setActive bool, updatedBy string) (*model.Machine, error)
	UpdateLabel(branchId string, label string, updatedBy string) (*model.Machine, error)
	GetByBranchID(branchId string) (*[]model.Machine, error)
	GetAll() (*[]model.Machine, error)
	AddMachine(newMachine *model.Machine) error
	GetByMachineSerial(machineSerial string) (*model.Machine, error)
	GetAvailableMachine(branchID string) (*[]model.MachineInBranch, error)
}

type machineRepository struct {
	db *platform.Postgres
}

func CreateMachineRepository(db *platform.Postgres) MachineRepository {
	return &machineRepository{db: db}
}

func (u *machineRepository) GetAll() (*[]model.Machine, error) {
	machineList := new([]model.Machine)
	result := u.db.Find(machineList)

	if result.Error != nil {
		return nil, result.Error
	}

	return machineList, result.Error
}

func (u *machineRepository) AddMachine(newMachine *model.Machine) error {
	result := u.db.Create(newMachine)

	return result.Error
}

func (u *machineRepository) SoftDelete(machineSerial string, deletedBy string) (*model.Machine, error) {
	deletedMachine := new(model.Machine)

	result := u.db.Where("machine_serial = ?", machineSerial).Delete(&deletedMachine)

	if result.RowsAffected == 0 {
		return nil, gorm.ErrRecordNotFound
	}

	if result.Error != nil {
		return nil, result.Error
	}

	queryErr := u.db.Unscoped().Model(&model.Machine{}).
		Where("machine_serial = ?", machineSerial).
		Update("machine_label", nil).
		Update("deleted_by", deletedBy).
		First(&deletedMachine)

	if queryErr.Error != nil {
		return nil, queryErr.Error
	}

	return deletedMachine, result.Error
}

func (u *machineRepository) UpdateActive(machineSerial string, isActive bool, updatedBy string) (*model.Machine, error) {
	updatedMachine := new(model.Machine)
	result := u.db.Model(&model.Machine{}).
		Where("machine_serial = ?", machineSerial).
		Update("is_active", isActive).
		Update("updated_by", updatedBy).
		Find(updatedMachine)

	if result.RowsAffected == 0 {
		return nil, gorm.ErrRecordNotFound
	}

	if result.Error != nil {
		return nil, result.Error
	}

	return updatedMachine, result.Error
}

func (u *machineRepository) UpdateLabel(machineSerial string, label string, updatedBy string) (*model.Machine, error) {
	updatedMachine := new(model.Machine)
	result := u.db.Model(&model.Machine{}).
		Where("machine_serial = ?", machineSerial).
		Update("machine_label", label).
		Update("updated_by", updatedBy).
		Find(updatedMachine)

	if result.RowsAffected == 0 {
		return nil, gorm.ErrRecordNotFound
	}

	if result.Error != nil {
		return nil, result.Error
	}

	return updatedMachine, result.Error
}

func (u *machineRepository) GetByMachineSerial(machineSerial string) (*model.Machine, error) {
	machine := new(model.Machine)

	result := u.db.Where("machine_serial = ?", machineSerial).First(machine)

	if result.RowsAffected == 0 {
		return nil, gorm.ErrRecordNotFound
	}

	if result.Error != nil {
		return nil, result.Error
	}

	return machine, nil
}

func (u *machineRepository) GetByBranchID(branchID string) (*[]model.Machine, error) {
	machine := new([]model.Machine)

	result := u.db.Where("branch_id = ?", branchID).Find(&machine)

	if result.Error != nil {
		return nil, result.Error
	}

	return machine, nil
}

func (u *machineRepository) GetAvailableMachine(branchID string) (*[]model.MachineInBranch, error) {
	machines := new([]model.MachineInBranch)

	result := u.db.Table("\"Machines\" AS M").
		Select(`M.machine_serial, OD.finished_at, M.weight, M.machine_label, M.machine_type,
			CASE 
					WHEN OD.order_status = 'Processing' THEN FALSE ELSE TRUE 
			END AS is_available`).
		Joins("LEFT JOIN \"OrderDetails\" AS OD ON M.machine_serial = OD.machine_serial").
		Where("M.branch_id = ? AND M.is_active = TRUE", branchID).
		Scan(&machines)

	if result.Error != nil {
		return nil, result.Error
	}

	return machines, nil
}
