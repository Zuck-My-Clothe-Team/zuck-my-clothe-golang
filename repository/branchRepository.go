package repository

import (
	"errors"
	"time"
	"zuck-my-clothe/zuck-my-clothe-backend/model"
	"zuck-my-clothe/zuck-my-clothe-backend/platform"
)

type branchReopository struct {
	db *platform.Postgres
}

func CreateNewBranchRepository(db *platform.Postgres) model.BranchReopository {
	return &branchReopository{db: db}
}

func (u *branchReopository) CreateBranch(newBranch *model.Branch) error {
	// retVal := new(model.Branch)
	dbTx := u.db.Create(newBranch)
	return dbTx.Error
}

func (u *branchReopository) GetAll() (*[]model.Branch, error) {
	branchList := new([]model.Branch)
	dbTx := u.db.Where("deleted_at = ? OR deleted_at >= ?", "0001-01-01T00:00:00Z", time.Now()).Find(branchList)
	return branchList, dbTx.Error
}

func (u *branchReopository) GetByBranchID(branchID string) (*model.Branch, error) {
	branch := new(model.Branch)
	// zeroTime, _ := time.Parse(time.RFC3339, "0001-01-01T00:00:00Z")
	dbTx := u.db.Where("branch_id = ? AND (deleted_at = ? OR deleted_at >= ?)", branchID, "0001-01-01T00:00:00Z", time.Now()).First(branch)
	if dbTx.Error != nil {
		return nil, dbTx.Error
	}
	return branch, nil
}

func (u *branchReopository) GetByBranchOwner(ownerUserID string) (*[]model.Branch, error) {
	branch := new([]model.Branch)
	dbTx := u.db.Table("Branches").Where("owner_user_id = ? AND (deleted_at = ? OR deleted_at >= ?)", ownerUserID, "0001-01-01T00:00:00Z", time.Now()).Find(branch)
	if dbTx.RowsAffected == 0 {
		return nil, errors.New("record not found")
	}
	return branch, dbTx.Error
}

func (u *branchReopository) UpdateBranch(branch *model.Branch) error {
	dbTx := u.db.Table("Branches").Where("branch_id = ? AND (deleted_at = ? OR deleted_at >= ?)", branch.BranchID, "0001-01-01T00:00:00Z", time.Now()).Updates(&model.Branch{BranchName: branch.BranchName, BranchDetail: branch.BranchDetail, BranchLat: branch.BranchLat, BranchLon: branch.BranchLon, OwnerUserID: branch.OwnerUserID})
	return dbTx.Error
}

func (u *branchReopository) ManagerUpdateBranch(branch *model.Branch) error {
	dbTx := u.db.Table("Branches").Where("branch_id = ? AND (deleted_at = ? OR deleted_at >= ?)", branch.BranchID, "0001-01-01T00:00:00Z", time.Now()).Updates(&model.Branch{BranchName: branch.BranchName, BranchDetail: branch.BranchDetail, BranchLat: branch.BranchLat, BranchLon: branch.BranchLon})
	return dbTx.Error
}

func (u *branchReopository) DeleteBranch(branch *model.Branch) error {
	dbTx := u.db.Table("Branches").Where("branch_id = ? AND (deleted_at = ? OR deleted_at >= ?)", branch.BranchID, "0001-01-01T00:00:00Z", time.Now()).Updates(&model.Branch{DeletedAt: time.Now(), DeletedBy: branch.DeletedBy})
	return dbTx.Error
}
