package repository

import (
	"zuck-my-clothe/zuck-my-clothe-backend/model"
	"zuck-my-clothe/zuck-my-clothe-backend/platform"

	"gorm.io/gorm"
)

type employeeContractRepository struct {
	db *platform.Postgres
}

type EmployeeContractRepository interface {
	CreateEmployeeContract(newEmployeeContract *model.EmployeeContract) error
	SoftDelete(contract_id string, deleted_by string) (*model.EmployeeContract, error)
	GetByBranchID(branch_id string) (*[]model.EmployeeContract, error)
	GetByUserID(user_id string) (*[]model.EmployeeContract, error)
	GetAll() (*[]model.EmployeeContract, error)
}

func CreateNewEmployeeContractRepository(db *platform.Postgres) EmployeeContractRepository {
	return &employeeContractRepository{db: db}
}

func (repo *employeeContractRepository) CreateEmployeeContract(newEmployeeContract *model.EmployeeContract) error {
	result := repo.db.Create(newEmployeeContract)
	return result.Error
}

// GetAll implements model.EmployeeContractRepository.
func (repo *employeeContractRepository) GetAll() (*[]model.EmployeeContract, error) {
	var contracts []model.EmployeeContract
	dbTx := repo.db.Find(&contracts)
	if dbTx.Error != nil {
		return nil, dbTx.Error
	}

	return &contracts, nil
}

func (repo *employeeContractRepository) GetByBranchID(branch_id string) (*[]model.EmployeeContract, error) {
	var contracts []model.EmployeeContract
	dbTx := repo.db.Where("branch_id = ?", branch_id).Find(&contracts)
	if dbTx.Error != nil {
		return nil, dbTx.Error
	}

	return &contracts, nil
}

func (repo *employeeContractRepository) GetByUserID(user_id string) (*[]model.EmployeeContract, error) {
	var contracts []model.EmployeeContract
	dbTx := repo.db.Where("user_id = ?", user_id).Find(&contracts)
	if dbTx.Error != nil {
		return nil, dbTx.Error
	}

	return &contracts, nil
}

func (repo *employeeContractRepository) SoftDelete(contract_id string, deleted_by string) (*model.EmployeeContract, error) {
	deleteContract := new(model.EmployeeContract)

	result := repo.db.Where("contract_id = ?", contract_id).Delete(&deleteContract)

	if result.RowsAffected == 0 {
		return nil, gorm.ErrRecordNotFound
	}

	if result.Error != nil {
		return nil, result.Error
	}

	queryErr := repo.db.Unscoped().Model(&model.EmployeeContract{}).Where("contract_id = ?", contract_id).Update("deleted_by", deleted_by).First(&deleteContract)

	if queryErr.Error != nil {
		return nil, queryErr.Error
	}

	return deleteContract, result.Error
}
