package usecases

import (
	"errors"
	"time"
	"zuck-my-clothe/zuck-my-clothe-backend/model"
	"zuck-my-clothe/zuck-my-clothe-backend/utils"

	"github.com/google/uuid"
)

type branchUsecase struct {
	branchRepository model.BranchReopository
}

func CreateNewBranchUsecase(branchRepository model.BranchReopository) model.BranchUsecase {
	return &branchUsecase{branchRepository: branchRepository}
}

func (u *branchUsecase) CreateBranch(newBranch *model.Branch, userID string) (*model.Branch, error) {
	// validate that all essential feild not contain null value
	//List of essential field
	// BranchName
	// BranchDetail
	// BranchLat
	// BranchLon
	// OwnerUserID
	// CreatedBy
	if utils.CheckStraoPling(newBranch.BranchName) ||
		utils.CheckStraoPling(newBranch.BranchDetail) ||
		utils.CheckStraoPling(newBranch.OwnerUserID) ||
		(newBranch.BranchLat == 0.0 || newBranch.BranchLon == 0.0) {
		return nil, errors.New("null detected on one or more essential field(s)")
	}
	//Generate BranchID
	newBranch.BranchID = uuid.New().String()
	newBranch.CreatedBy = userID
	newBranch.CreatedAt = time.Now()
	newBranch.UpdatedBy = userID
	newBranch.UpdatedAt = time.Now()

	err := u.branchRepository.CreateBranch(newBranch)
	if err != nil {
		return nil, err
	}
	branch, err := u.branchRepository.GetByBranchID(newBranch.BranchID)
	if err != nil {
		return nil, err
	}
	return branch, nil
}

func (u *branchUsecase) GetAll() (*[]model.Branch, error) {
	branchList, err := u.branchRepository.GetAll()
	return branchList, err
}

func (u *branchUsecase) GetByBranchID(branchID string) (*model.Branch, error) {
	branch, err := u.branchRepository.GetByBranchID(branchID)
	branch.DeletedBy = ""
	branch.UpdatedBy = ""
	return branch, err
}

func (u *branchUsecase) GetByBranchOwner(owenerUserID string) (*[]model.Branch, error) {
	branch, err := u.branchRepository.GetByBranchOwner(owenerUserID)
	return branch, err
}

func (u *branchUsecase) UpdateBranch(branch *model.Branch, role string) (*model.Branch, error) {
	if !(utils.CheckStraoPling(branch.CreatedBy) &&
		utils.CheckStraoPling(branch.DeletedBy) &&
		(utils.CheckStraoPling(branch.OwnerUserID) || role == "SuperAdmin") &&
		branch.CreatedAt.IsZero() &&
		branch.UpdatedAt.IsZero() &&
		branch.DeletedAt.IsZero()) {
		return nil, errors.New("update error: you are trying to edit one or more un-editable field(s)")
	}

	if utils.CheckStraoPling(branch.BranchName) ||
		utils.CheckStraoPling(branch.BranchDetail) ||
		(utils.CheckStraoPling(branch.OwnerUserID) && role == "SuperAdmin") ||
		(branch.OwnerUserID == "c2f45bf6-6308-450f-9932-89632c40974c" && role == "SuperAdmin") ||
		(branch.BranchLat == 0.0 || branch.BranchLon == 0.0) {
		return nil, errors.New("null detected on one or more essential field(s)")
	}

	if role == "SuperAdmin" {
		if err := u.branchRepository.UpdateBranch(branch); err != nil {
			return nil, err
		}
	} else if role == "BranchManager" {
		if err := u.branchRepository.ManagerUpdateBranch(branch); err != nil {
			return nil, err
		}
	}

	response, err := u.GetByBranchID(branch.BranchID)
	return response, err
}

func (u *branchUsecase) DeleteBranch(branch *model.Branch) error {
	if utils.CheckStraoPling(branch.BranchID) ||
		utils.CheckStraoPling(branch.DeletedBy) {
		return errors.New("null detected on one or more essential field(s)")
	}
	err := u.branchRepository.DeleteBranch(branch)
	return err
}
