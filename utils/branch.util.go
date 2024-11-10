package utils

import "zuck-my-clothe/zuck-my-clothe-backend/model"

func ToBranchDetail(branch *model.Branch, isAdminView bool) model.BranchDetail {
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
