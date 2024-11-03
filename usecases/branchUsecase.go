package usecases

import (
	"errors"
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

func toBranchDetail(branch *model.Branch) model.BranchDetail {
	res := model.BranchDetail{
		BranchID:     branch.BranchID,
		BranchName:   branch.BranchName,
		BranchDetail: branch.BranchDetail,
		BranchLat:    branch.BranchLat,
		BranchLon:    branch.BranchLon,
		OwnerUserID:  branch.OwnerUserID,
	}

	return res
}

func (u *branchUsecase) CreateBranch(newBranch *model.CreateBranchDTO, userID string) (*model.BranchDetail, error) {
	data := model.Branch{
		BranchID:     uuid.New().String(),
		BranchName:   newBranch.BranchName,
		BranchDetail: newBranch.BranchDetail,
		BranchLat:    newBranch.BranchLat,
		BranchLon:    newBranch.BranchLon,
		CreatedBy:    userID,
		OwnerUserID:  newBranch.OwnerUserID,
		UpdatedBy:    userID,
		DeletedBy:    nil,
	}

	err := u.branchRepository.CreateBranch(&data)
	if err != nil {
		return nil, err
	}

	branch, err := u.branchRepository.GetByBranchID(data.BranchID)
	if err != nil {
		return nil, err
	}

	res := toBranchDetail(branch)

	return &res, nil
}

func (u *branchUsecase) GetAll() (*[]model.BranchDetail, error) {
	branchList, err := u.branchRepository.GetAll()

	var result []model.BranchDetail

	for _, branch := range *branchList {
		result = append(result, toBranchDetail(&branch))
	}

	return &result, err
}

func (u *branchUsecase) GetByBranchID(branchID string) (*model.BranchDetail, error) {
	branch, err := u.branchRepository.GetByBranchID(branchID)

	if err != nil {
		return nil, err
	}

	res := toBranchDetail(branch)

	return &res, err
}

func (u *branchUsecase) GetByBranchOwner(owenerUserID string) (*[]model.BranchDetail, error) {
	branchList, err := u.branchRepository.GetByBranchOwner(owenerUserID)
	var result []model.BranchDetail

	for _, branch := range *branchList {
		result = append(result, toBranchDetail(&branch))
	}

	return &result, err
}

func (u *branchUsecase) UpdateBranch(branch *model.UpdateBranchDTO, role string) (*model.BranchDetail, error) {
	data := model.Branch{
		BranchID:     branch.BranchID,
		BranchName:   branch.BranchName,
		BranchDetail: branch.BranchDetail,
		BranchLat:    branch.BranchLat,
		BranchLon:    branch.BranchLon,
		OwnerUserID:  branch.OwnerUserID,
	}

	if role == "SuperAdmin" {
		if err := u.branchRepository.UpdateBranch(&data); err != nil {
			return nil, err
		}
	} else if role == "BranchManager" {
		if err := u.branchRepository.ManagerUpdateBranch(&data); err != nil {
			return nil, err
		}
	}

	response, err := u.branchRepository.GetByBranchID(branch.BranchID)
	res := toBranchDetail(response)

	return &res, err
}

func (u *branchUsecase) DeleteBranch(branch *model.Branch) error {
	if utils.CheckStraoPling(branch.BranchID) {
		return errors.New("null detected on one or more essential field(s)")
	}
	err := u.branchRepository.DeleteBranch(branch)
	return err
}
