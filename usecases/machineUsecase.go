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

func toMachineDetail(machine *model.Machine) interface{} {
	result := model.MachineDetail{
		MachineSerial: machine.MachineSerial,
		BranchID:      machine.BranchID,
		MachineType:   machine.MachineType,
		Weight:        int16(machine.Weight),
		IsActive:      machine.IsActive,
	}
	return result
}

func (u *machineUsecase) AddMachine(new_machine *model.AddMachineDTO) (*model.Machine, error) {
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

	return machine, nil
}

func (u *machineUsecase) GetAll() (*[]model.Machine, error) {
	var machines *[]model.Machine

	machines, err := u.machineRepository.GetAll()

	if err != nil {
		return nil, err
	}

	return machines, err
}

func (u *machineUsecase) GetByMachineSerial(machine_serial string, isAdminView bool) (*interface{}, error) {
	machine, err := u.machineRepository.GetByMachineSerial(machine_serial)

	if err != nil {
		return nil, err
	}

	var result interface{} = machine

	if !isAdminView {
		result = toMachineDetail(machine)
	}

	return &result, err
}

func (u *machineUsecase) GetByBranchID(branch_id string, isAdminView bool) (*[]interface{}, error) {
	machines, err := u.machineRepository.GetByBranchID(branch_id)

	if err != nil {
		return nil, err
	}

	var result []interface{}

	if len(*machines) == 0 {
		result = []interface{}{}
		return &result, err
	}

	if isAdminView {
		for _, machine := range *machines {
			result = append(result, machine)
		}
	} else {
		for _, machine := range *machines {
			result = append(result, toMachineDetail(&machine))
		}
	}

	return &result, err
}

func (u *machineUsecase) UpdateActive(machine_serial string, set_active bool, updated_by string) (*model.Machine, error) {

	updated_machine, err := u.machineRepository.UpdateActive(machine_serial, set_active, updated_by)

	if err != nil {
		return nil, err
	}

	return updated_machine, err
}

func (u *machineUsecase) SoftDelete(machine_serial string, deleted_by string) (*model.Machine, error) {
	deleted_machine, err := u.machineRepository.SoftDelete(machine_serial, deleted_by)

	if err != nil {
		return nil, err
	}

	return deleted_machine, err
}
