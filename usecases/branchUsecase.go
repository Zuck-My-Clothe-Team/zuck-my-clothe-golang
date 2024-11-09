package usecases

import (
	"errors"
	"zuck-my-clothe/zuck-my-clothe-backend/model"
	repo "zuck-my-clothe/zuck-my-clothe-backend/repository"
	"zuck-my-clothe/zuck-my-clothe-backend/utils"

	"github.com/google/uuid"
)

type branchUsecase struct {
	branchRepository  repo.BranchReopository
	machineRepository repo.MachineRepository
}

type BranchUsecase interface {
	CreateBranch(newBranch *model.CreateBranch, userID string) (*model.Branch, error)
	GetAll(isAdminView bool) (interface{}, error)
	GetClosestToMe(userLocation *model.UserGeoLocation) (*[]model.BranchDetail, error)
	GetByBranchID(branchID string, isAdminView bool) (*model.BranchDetail, error)
	GetByBranchOwner(ownerUserID string) (*[]model.Branch, error)
	UpdateBranch(branch *model.UpdateBranch, role string) (*model.Branch, error)
	DeleteBranch(branch *model.Branch) error
}

func CreateNewBranchUsecase(branchRepository repo.BranchReopository, machineRepository repo.MachineRepository) BranchUsecase {
	return &branchUsecase{
		branchRepository:  branchRepository,
		machineRepository: machineRepository,
	}
}

func toBranchDetail(branch *model.Branch, isAdminView bool) model.BranchDetail {
	res := model.BranchDetail{
		BranchID:     branch.BranchID,
		BranchName:   branch.BranchName,
		BranchDetail: branch.BranchDetail,
		BranchLat:    branch.BranchLat,
		BranchLon:    branch.BranchLon,
		OwnerUserID:  branch.OwnerUserID,
		CreatedAt:    &branch.CreatedAt,
		CreatedBy:    &branch.CreatedBy,
		UpdatedAt:    &branch.UpdatedAt,
		UpdatedBy:    &branch.UpdatedBy,
	}

	if !isAdminView {
		res.CreatedBy = nil
		res.UpdatedBy = nil
	}

	return res
}

func (u *branchUsecase) CreateBranch(newBranch *model.CreateBranch, userID string) (*model.Branch, error) {
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
			branchDetailList = append(branchDetailList, toBranchDetail(&branch, false))
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
		res = append(res, toBranchDetail(&branch, false))
	}

	sortedBranches := utils.SortBranchesByDistance(5, userLocation.BranchLat, userLocation.BranchLon, res)

	return &sortedBranches, err
}

func (u *branchUsecase) GetByBranchID(branchID string, isAdminView bool) (*model.BranchDetail, error) {
	branch, err := u.branchRepository.GetByBranchID(branchID)

	if err != nil {
		return nil, err
	}

	res := toBranchDetail(branch, isAdminView)

	// Get available machines
	machines, err := u.machineRepository.GetAvailableMachine(branchID)

	if err != nil {
		return nil, err
	}

	if len(*machines) != 0 {
		res.AvailableMachine = machines
	} else {
		res.AvailableMachine = &[]model.MachineInBranch{}
	}

	// Get review
	reviews, err := u.branchRepository.GetReviewsByBranchID(branchID)

	if err != nil {
		return nil, err
	}

	if len(*reviews) > 0 {
		res.UserReview = reviews

		// Calculate average star
		var averageStar float32 = 0

		for _, r := range *reviews {
			averageStar += float32(r.StarRating)
		}

		averageStar /= float32(len(*reviews))

		res.AverageStar = averageStar
	} else {
		res.UserReview = &[]model.UserReview{}
		res.AverageStar = 0
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

func (u *branchUsecase) UpdateBranch(branch *model.UpdateBranch, role string) (*model.Branch, error) {
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
