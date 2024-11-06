package usecases

import (
	"errors"
	"zuck-my-clothe/zuck-my-clothe-backend/model"
	"zuck-my-clothe/zuck-my-clothe-backend/repository"
	"zuck-my-clothe/zuck-my-clothe-backend/utils"

	"github.com/google/uuid"
)

type BranchUsecase interface {
	CreateBranch(newBranch *model.CreateBranchDTO, userID string) (*model.Branch, error)
	GetAll(isAdminView bool) (interface{}, error)
	GetClosestToMe(userLocation *model.UserGeoLocation) (*[]model.BranchDetail, error)
	GetByBranchID(branchID string, isAdminView bool) (*interface{}, error)
	GetByBranchOwner(ownerUserID string) (*[]model.Branch, error)
	UpdateBranch(branch *model.UpdateBranchDTO, role string) (*model.Branch, error)
	DeleteBranch(branch *model.Branch) error
}

type branchUsecase struct {
	branchRepository repository.BranchReopository
}

func CreateNewBranchUsecase(branchRepository repository.BranchReopository) BranchUsecase {
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

func (u *branchUsecase) CreateBranch(newBranch *model.CreateBranchDTO, userID string) (*model.Branch, error) {
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

	return branch, nil
}

func (u *branchUsecase) GetAll(isAdminView bool) (interface{}, error) {
	branchList, err := u.branchRepository.GetAll()

	if err != nil {
		return nil, err
	}

	var res interface{} = *branchList

	if !isAdminView {
		var branchDetailList []model.BranchDetail

		for _, branch := range *branchList {
			branchDetailList = append(branchDetailList, toBranchDetail(&branch))
		}

		res = branchDetailList
	}

	return res, err
}

func (u *branchUsecase) GetClosestToMe(userLocation *model.UserGeoLocation) (*[]model.BranchDetail, error) {
	branchList, err := u.branchRepository.GetAll()
	if err != nil {
		return nil, err
	}
	var res []model.BranchDetail

	for _, branch := range *branchList {
		res = append(res, toBranchDetail(&branch))
	}

	sortedBranches := utils.SortBranchesByDistance(5, userLocation.BranchLat, userLocation.BranchLon, res)

	return &sortedBranches, err
}

func (u *branchUsecase) GetByBranchID(branchID string, isAdminView bool) (*interface{}, error) {
	branch, err := u.branchRepository.GetByBranchID(branchID)

	if err != nil {
		return nil, err
	}

	var res interface{} = branch

	if !isAdminView {
		res = toBranchDetail(branch)
	}

	return &res, err
}

func (u *branchUsecase) GetByBranchOwner(ownerUserID string) (*[]model.Branch, error) {
	branchList, err := u.branchRepository.GetByBranchOwner(ownerUserID)

	if err != nil {
		return nil, err
	}

	return branchList, err
}

func (u *branchUsecase) UpdateBranch(branch *model.UpdateBranchDTO, role string) (*model.Branch, error) {
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

	if err != nil {
		return nil, err
	}

	return response, err
}

func (u *branchUsecase) DeleteBranch(branch *model.Branch) error {
	if utils.CheckStraoPling(branch.BranchID) {
		return errors.New("null detected on one or more essential field(s)")
	}
	err := u.branchRepository.DeleteBranch(branch)
	return err
}
