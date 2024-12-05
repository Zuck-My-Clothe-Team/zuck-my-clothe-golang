package usecases

import (
	"strconv"
	"time"
	"zuck-my-clothe/zuck-my-clothe-backend/model"
	"zuck-my-clothe/zuck-my-clothe-backend/repository"
)

type MachineUsecase interface {
	SoftDelete(machine_serial string, deleted_by string) (*model.Machine, error)
	UpdateActive(machine_serial string, set_active bool, updated_by string) (*model.Machine, error)
	UpdateLabel(machine_serial string, label int, updated_by string) (*model.Machine, error)
	GetByBranchID(branch_id string, isAdminView bool) (*[]interface{}, error)
	GetAll() (*[]model.Machine, error)
	AddMachine(newMachine *model.AddMachine) (*model.Machine, error)
	GetByMachineSerial(machineSerial string, isAdminView bool, withTime bool) (*interface{}, error)
	GetAvailableMachineInBranch(branchID string) (*[]model.MachineInBranch, error)
}

type machineUsecase struct {
	machineRepository repository.MachineRepository
}

func CreateMachineUsecase(machineRepository repository.MachineRepository) MachineUsecase {
	return &machineUsecase{machineRepository: machineRepository}
}

func toMachineDetail(machine *model.Machine) interface{} {
	result := model.MachineDetail{
		MachineSerial: machine.MachineSerial,
		MachineLabel:  machine.MachineLabel,
		BranchID:      machine.BranchID,
		MachineType:   machine.MachineType,
		Weight:        int16(machine.Weight),
		IsActive:      machine.IsActive,
	}
	return result
}

func (u *machineUsecase) AddMachine(new_machine *model.AddMachine) (*model.Machine, error) {
	var newMachineLabel string

	if new_machine.MachineType == model.Washer {
		newMachineLabel = "เครื่องซักที่ " + strconv.Itoa(new_machine.MachineLabel)
	} else {
		newMachineLabel = "เครื่องอบที่ " + strconv.Itoa(new_machine.MachineLabel)
	}

	machine_data := model.Machine{
		MachineSerial: new_machine.MachineSerial,
		MachineLabel:  newMachineLabel,
		BranchID:      new_machine.BranchID,
		MachineType:   new_machine.MachineType,
		Weight:        int16(new_machine.Weight),
		IsActive:      false,
		CreatedAt:     time.Now().UTC(),
		CreatedBy:     &new_machine.CreatedBy,
		UpdatedAt:     time.Now().UTC(),
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

func (u *machineUsecase) GetByMachineSerial(machineSerial string, isAdminView bool, withTime bool) (*interface{}, error) {
	var result interface{}

	if withTime {
		machineWithTime, err := u.machineRepository.GetWithTime(machineSerial)
		if err != nil {
			return nil, err
		}
		result = machineWithTime
	} else {
		machine, err := u.machineRepository.GetByMachineSerial(machineSerial)
		if err != nil {
			return nil, err
		}
		result = machine

		if !isAdminView {
			result = toMachineDetail(machine)
		}
	}

	return &result, nil
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

func (u *machineUsecase) GetAvailableMachineInBranch(branchID string) (*[]model.MachineInBranch, error) {
	response, err := u.machineRepository.GetAvailableMachine(branchID)
	if err != nil {
		return nil, err
	}
	return response, nil
}

func (u *machineUsecase) UpdateActive(machine_serial string, set_active bool, updated_by string) (*model.Machine, error) {

	updated_machine, err := u.machineRepository.UpdateActive(machine_serial, set_active, updated_by)

	if err != nil {
		return nil, err
	}

	return updated_machine, err
}

func (u *machineUsecase) UpdateLabel(machine_serial string, label int, updated_by string) (*model.Machine, error) {

	current_machine, err := u.machineRepository.GetByMachineSerial(machine_serial)
	if err != nil {
		return nil, err
	}
	var newMachineLabel string

	if current_machine.MachineType == model.Washer {
		newMachineLabel = "เครื่องซักที่ " + strconv.Itoa(label)
	} else {
		newMachineLabel = "เครื่องอบที่ " + strconv.Itoa(label)
	}

	updated_machine, err := u.machineRepository.UpdateLabel(machine_serial, newMachineLabel, updated_by)

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
