package usecases

import (
	"fmt"
	"zuck-my-clothe/zuck-my-clothe-backend/model"
	"zuck-my-clothe/zuck-my-clothe-backend/repository"
	validatorboi "zuck-my-clothe/zuck-my-clothe-backend/validator"

	"github.com/google/uuid"
)

type EmployeeContractUsecases interface {
	CreateEmployeeContract(newEmployeeContract *model.EmployeeContractDTO) error
	SoftDelete(contract_id string, deleted_by string) (*model.EmployeeContract, error)
	GetByBranchID(branch_id string) (*[]model.EmployeeContract, error)
	GetByUserID(user_id string) (*[]model.EmployeeContract, error)
	GetAll() (*[]model.EmployeeContract, error)
}

type employeeContractUsecases struct {
	repository     repository.EmployeeContractRepository
	userRepository model.UserRepository
}

func CreateNewEmployeeContractUsecase(
	employeeContractRepository repository.EmployeeContractRepository,
	userRepository model.UserRepository,
) EmployeeContractUsecases {
	return &employeeContractUsecases{
		repository:     employeeContractRepository,
		userRepository: userRepository,
	}
}

func (usecase *employeeContractUsecases) CreateEmployeeContract(newEmployeeContract *model.EmployeeContractDTO) error {

	if err := validatorboi.Validate(newEmployeeContract); err != nil {
		return err
	}

	userRecord, err := usecase.userRepository.FindUserByUserID(newEmployeeContract.UserID)

	if err != nil {
		return err
	}

	if userRecord == nil {
		return fmt.Errorf("user not found")
	}

	if userRecord.Role != model.Employee {
		return fmt.Errorf("user is not an employee")
	}

	record, err := usecase.repository.GetByUserID(newEmployeeContract.UserID)
	if err != nil && err.Error() != "record not found" {
		return err
	}

	if record != nil {
		for _, contract := range *record {
			if contract.BranchID == newEmployeeContract.BranchID && contract.PositionId == newEmployeeContract.PositionId && contract.DeletedAt == (model.EmployeeContract{}).DeletedAt {
				return fmt.Errorf("contract already exist")
			}
		}
	}

	newContract := model.EmployeeContract{
		ContractId: uuid.New().String(),
		UserID:     newEmployeeContract.UserID,
		BranchID:   newEmployeeContract.BranchID,
		PositionId: newEmployeeContract.PositionId,
		CreatedBy:  newEmployeeContract.UserID,
		DeletedBy:  nil,
	}

	err = usecase.repository.CreateEmployeeContract(&newContract)
	if err != nil {
		return err
	}

	return nil
}

func (usecase *employeeContractUsecases) SoftDelete(contract_id string, deleted_by string) (*model.EmployeeContract, error) {
	deletedUser, err := usecase.repository.SoftDelete(contract_id, deleted_by)
	return deletedUser, err
}

func (usecase *employeeContractUsecases) GetByBranchID(branch_id string) (*[]model.EmployeeContract, error) {
	record, err := usecase.repository.GetByBranchID(branch_id)
	return record, err
}

func (usecase *employeeContractUsecases) GetByUserID(user_id string) (*[]model.EmployeeContract, error) {
	record, err := usecase.repository.GetByUserID(user_id)
	return record, err
}

func (usecase *employeeContractUsecases) GetAll() (*[]model.EmployeeContract, error) {
	record, err := usecase.repository.GetAll()
	return record, err
}
