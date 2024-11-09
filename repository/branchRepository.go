package repository

import (
	"errors"
	"zuck-my-clothe/zuck-my-clothe-backend/model"
	"zuck-my-clothe/zuck-my-clothe-backend/platform"

	"gorm.io/gorm"
)

type BranchReopository interface {
	CreateBranch(newBranch *model.Branch) error
	GetAll() (*[]model.Branch, error)
	GetByBranchID(branchID string) (*model.Branch, error)
	GetByBranchOwner(ownerUserID string) (*[]model.Branch, error)
	GetReviewsByBranchID(branchID string) (*[]model.UserReview, error)
	UpdateBranch(branch *model.Branch) error
	ManagerUpdateBranch(branch *model.Branch) error
	DeleteBranch(branch *model.Branch) error
}

type branchReopository struct {
	db *platform.Postgres
}

func CreateNewBranchRepository(db *platform.Postgres) BranchReopository {
	return &branchReopository{db: db}
}

func (u *branchReopository) CreateBranch(newBranch *model.Branch) error {
	dbTx := u.db.Create(newBranch)
	return dbTx.Error
}

func (u *branchReopository) GetAll() (*[]model.Branch, error) {
	branchList := new([]model.Branch)
	dbTx := u.db.Find(branchList)

	if dbTx.Error != nil {
		return nil, dbTx.Error
	}

	return branchList, dbTx.Error
}

func (u *branchReopository) GetByBranchID(branchID string) (*model.Branch, error) {
	branch := new(model.Branch)
	dbTx := u.db.Where("branch_id = ?", branchID).First(branch)

	if dbTx.Error != nil {
		return nil, dbTx.Error
	}

	return branch, nil
}

func (u *branchReopository) GetByBranchOwner(ownerUserID string) (*[]model.Branch, error) {
	branch := new([]model.Branch)
	dbTx := u.db.Where("owner_user_id = ?", ownerUserID).Find(branch)

	if dbTx.RowsAffected == 0 {
		return nil, errors.New("record not found")
	}

	return branch, dbTx.Error
}

func (u *branchReopository) UpdateBranch(branch *model.Branch) error {
	dbTx := u.db.Where("branch_id = ?", branch.BranchID).Updates(&model.Branch{BranchName: branch.BranchName, BranchDetail: branch.BranchDetail, BranchLat: branch.BranchLat, BranchLon: branch.BranchLon, OwnerUserID: branch.OwnerUserID})
	return dbTx.Error
}

func (u *branchReopository) ManagerUpdateBranch(branch *model.Branch) error {
	dbTx := u.db.Where("branch_id = ?", branch.BranchID).Updates(&model.Branch{BranchName: branch.BranchName, BranchDetail: branch.BranchDetail, BranchLat: branch.BranchLat, BranchLon: branch.BranchLon})
	return dbTx.Error
}

func (u *branchReopository) DeleteBranch(branch *model.Branch) error {
	deleted_branch := new(model.Branch)

	queryErr := u.db.Unscoped().Model(&model.Branch{}).Where("branch_id = ?", branch.BranchID).Update("deleted_by", branch.DeletedBy).First(&deleted_branch)

	if queryErr.Error != nil {
		return queryErr.Error
	}

	result := u.db.Where("branch_id = ?", branch.BranchID).Delete(&deleted_branch)

	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	if result.Error != nil {
		return result.Error
	}

	return result.Error
}

func (u *branchReopository) GetReviewsByBranchID(branchID string) (*[]model.UserReview, error) {
	userReviews := new([]model.UserReview)

	result := u.db.Table("\"OrderHeaders\" as oh").
		Select("oh.order_header_id, oh.updated_at AS reviewed_at, oh.zuck_onsite, oh.star_rating, oh.review_comment, usr.firstname, usr.lastname, usr.profile_image_url").
		Joins("LEFT JOIN \"Users\" usr ON oh.user_id = usr.user_id").
		Where("oh.star_rating IS NOT NULL").
		Where("oh.branch_id = ?", branchID).
		Scan(&userReviews)

	if result.RowsAffected == 0 {
		return &[]model.UserReview{}, nil
	}

	if result.Error != nil {
		return nil, result.Error
	}

	return userReviews, nil
}
