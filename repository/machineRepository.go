package repository

import (
	"zuck-my-clothe/zuck-my-clothe-backend/model"
	"zuck-my-clothe/zuck-my-clothe-backend/platform"

	"gorm.io/gorm"
)

type machineRepository struct {
	db *platform.Postgres
}

func CreateMachineRepository(db *platform.Postgres) model.MachineRepository {
	return &machineRepository{db: db}
}

func (u *machineRepository) GetAll() (*[]model.Machine, error) {
	machineList := new([]model.Machine)
	// result := u.db.Unscoped().Find(machineList)
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

func (u *machineRepository) SoftDelete(machine_serial string, deleted_by string) (*model.Machine, error) {
	deleted_machine := new(model.Machine)

	result := u.db.Where("machine_serial = ?", machine_serial).Delete(&deleted_machine)

	if result.RowsAffected == 0 {
		return nil, gorm.ErrRecordNotFound
	}

	if result.Error != nil {
		return nil, result.Error
	}

	queryErr := u.db.Unscoped().Model(&model.Machine{}).Where("machine_serial = ?", machine_serial).Update("deleted_by", deleted_by).First(&deleted_machine)

	if queryErr.Error != nil {
		return nil, queryErr.Error
	}

	return deleted_machine, result.Error
}

func (u *machineRepository) UpdateActive(machine_serial string, is_active bool, updated_by string) (*model.Machine, error) {
	updated_machine := new(model.Machine)
	result := u.db.Model(&model.Machine{}).
		Where("machine_serial = ?", machine_serial).
		Update("is_active", is_active).
		Update("updated_by", updated_by).
		Find(updated_machine)

	if result.RowsAffected == 0 {
		return nil, gorm.ErrRecordNotFound
	}

	if result.Error != nil {
		return nil, result.Error
	}

	return updated_machine, result.Error
}

func (u *machineRepository) GetByMachineSerial(machineSerial string) (*model.Machine, error) {
	machine := new(model.Machine)

	result := u.db.Where("machine_serial = ?", machineSerial).First(machine)

	if result.Error != nil {
		return nil, result.Error
	}

	return machine, nil
}

func (u *machineRepository) GetByBranchID(branch_id string) (*[]model.Machine, error) {
	machine := new([]model.Machine)

	result := u.db.Where("branch_id = ?", branch_id).Find(&machine)

	if result.Error != nil {
		return nil, result.Error
	}

	return machine, nil
}
