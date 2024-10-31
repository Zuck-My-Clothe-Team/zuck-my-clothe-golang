package usecases

import (
	"time"
	"zuck-my-clothe/zuck-my-clothe-backend/model"
)

type machineUsecase struct {
	machineRepository model.MachineRepository
}

func CreateMachineUsecase(machineRepository model.MachineRepository) model.MachineUsecase {
	return &machineUsecase{machineRepository: machineRepository}
}

func (u *machineUsecase) AddMachine(new_machine *model.AddMachineDTO) (*model.MachineDetail, error) {
	machine_data := model.Machine{
		MachineSerial: new_machine.MachineSerial,
		BranchID:      new_machine.BranchID,
		MachineType:   new_machine.MachineType,
		Weight:        int16(new_machine.Weight),
		IsActive:      false,
		CreatedAt:     time.Now(),
		CreatedBy:     &new_machine.CreatedBy,
		UpdatedAt:     time.Now(),
		UpdatedBy:     &new_machine.CreatedBy,
		DeletedBy:     nil,
	}

	err := u.machineRepository.AddMachine(&machine_data)
	if err != nil {
		return nil, err
	}

	machine, err := u.machineRepository.GetByMachineSerial(machine_data.MachineSerial)

	if err != nil {
		return nil, err
	}

	result := model.MachineDetail{
		MachineSerial: machine.MachineSerial,
		BranchID:      machine.BranchID,
		MachineType:   machine.MachineType,
		Weight:        int16(machine.Weight),
		IsActive:      false,
	}

	return &result, nil
}

func (u *machineUsecase) GetAll() (*[]model.MachineDetail, error) {
	var machines *[]model.Machine

	machines, err := u.machineRepository.GetAll()

	if err != nil {
		return nil, err
	}

	var result []model.MachineDetail

	for _, machine := range *machines {
		res := model.MachineDetail{
			MachineSerial: machine.MachineSerial,
			BranchID:      machine.BranchID,
			MachineType:   machine.MachineType,
			Weight:        int16(machine.Weight),
			IsActive:      machine.IsActive,
		}

		result = append(result, res)
	}

	return &result, err
}

func (u *machineUsecase) GetByMachineSerial(machine_serial string) (*model.MachineDetail, error) {
	machine, err := u.machineRepository.GetByMachineSerial(machine_serial)

	if err != nil {
		return nil, err
	}

	result := model.MachineDetail{
		MachineSerial: machine.MachineSerial,
		BranchID:      machine.BranchID,
		MachineType:   machine.MachineType,
		Weight:        int16(machine.Weight),
		IsActive:      false,
	}

	return &result, err
}

func (u *machineUsecase) GetByBranchID(branch_id string) (*[]model.MachineDetail, error) {
	machines, err := u.machineRepository.GetByBranchID(branch_id)

	if err != nil {
		return nil, err
	}

	var result []model.MachineDetail

	for _, machine := range *machines {
		res := model.MachineDetail{
			MachineSerial: machine.MachineSerial,
			BranchID:      machine.BranchID,
			MachineType:   machine.MachineType,
			Weight:        int16(machine.Weight),
			IsActive:      machine.IsActive,
		}

		result = append(result, res)
	}

	return &result, err
}

func (u *machineUsecase) UpdateActive(machine_serial string, set_active bool, updated_by string) (*model.MachineDetail, error) {

	updated_machine, err := u.machineRepository.UpdateActive(machine_serial, set_active, updated_by)

	if err != nil {
		return nil, err
	}

	result := model.MachineDetail{
		MachineSerial: updated_machine.MachineSerial,
		BranchID:      updated_machine.BranchID,
		MachineType:   updated_machine.MachineType,
		Weight:        int16(updated_machine.Weight),
		IsActive:      updated_machine.IsActive,
	}

	return &result, err
}

func (u *machineUsecase) SoftDelete(machine_serial string, deleted_by string) (*model.MachineDetail, error) {
	deleted_machine, err := u.machineRepository.SoftDelete(machine_serial, deleted_by)

	if err != nil {
		return nil, err
	}

	result := model.MachineDetail{
		MachineSerial: deleted_machine.MachineSerial,
		BranchID:      deleted_machine.BranchID,
		MachineType:   deleted_machine.MachineType,
		Weight:        int16(deleted_machine.Weight),
		IsActive:      deleted_machine.IsActive,
	}

	return &result, err
}
