package model

import (
	"time"

	"gorm.io/gorm"
)

func (Branch) TableName() string {
	return "Branches"
}

type Branch struct {
	BranchID     string         `json:"branch_id" gorm:"column:branch_id;primaryKey"`
	BranchName   string         `json:"branch_name" gorm:"column:branch_name"`
	BranchDetail string         `json:"branch_detail" gorm:"column:branch_detail"`
	BranchLat    float64        `json:"branch_lat" gorm:"column:branch_lat"`
	BranchLon    float64        `json:"branch_long" gorm:"column:branch_long"`
	OwnerUserID  string         `json:"owner_user_id" gorm:"column:owner_user_id"`
	CreatedAt    time.Time      `json:"created_at" gorm:"column:created_at"`
	CreatedBy    string         `json:"created_by" gorm:"column:created_by"`
	UpdatedAt    time.Time      `json:"updated_at" gorm:"column:updated_at"`
	UpdatedBy    string         `json:"updated_by" gorm:"column:updated_by"`
	DeletedAt    gorm.DeletedAt `json:"deleted_at" gorm:"column:deleted_at;index" swaggertype:"string" example:"null"`
	DeletedBy    *string        `json:"deleted_by" gorm:"column:deleted_by"`
}

type CreateBranch struct {
	BranchName   string  `json:"branch_name" validate:"required"`
	BranchDetail string  `json:"branch_detail" validate:"required"`
	BranchLat    float64 `json:"branch_lat" validate:"required"`
	BranchLon    float64 `json:"branch_long" validate:"required"`
	OwnerUserID  string  `json:"owner_user_id" validate:"required"`
}
type UpdateBranch struct {
	BranchID     string  `json:"branch_id" validate:"required"`
	BranchName   string  `json:"branch_name" validate:"required"`
	BranchDetail string  `json:"branch_detail" validate:"required"`
	BranchLat    float64 `json:"branch_lat" validate:"required"`
	BranchLon    float64 `json:"branch_long" validate:"required"`
	OwnerUserID  string  `json:"owner_user_id" validate:"required"`
}

type BranchDetail struct {
	BranchID         string             `json:"branch_id"`
	BranchName       string             `json:"branch_name"`
	BranchDetail     string             `json:"branch_detail"`
	BranchLat        float64            `json:"branch_lat"`
	BranchLon        float64            `json:"branch_long"`
	OwnerUserID      string             `json:"owner_user_id"`
	CreatedAt        *time.Time         `json:"created_at,omitempty"`
	CreatedBy        *string            `json:"created_by,omitempty"`
	UpdatedAt        *time.Time         `json:"updated_at,omitempty"`
	UpdatedBy        *string            `json:"updated_by,omitempty"`
	DeletedAt        *gorm.DeletedAt    `json:"deleted_at,omitempty" swaggertype:"string" example:"null"`
	DeletedBy        *string            `json:"deleted_by,omitempty"`
	Distance         float64            `json:"distance,omitempty"`
	AverageStar      float32            `json:"average_star"`
	UserReview       *[]UserReview      `json:"user_reviews,omitempty"`
	AvailableMachine *[]MachineInBranch `json:"machines"`
}

type UserReview struct {
	StarRating      int16   `json:"star_rating" gorm:"column:star_rating"`
	ReviewComment   *string `json:"review_comment" gorm:"column:review_comment"`
	FirstName       string  `json:"firstname" gorm:"column:firstname"`
	LastName        string  `json:"lastname" gorm:"column:lastname"`
	ProfileImageURL string  `json:"profile_image_url" gorm:"column:profile_image_url"`
}

type MachineInBranch struct {
	MachineSerial string      `json:"machine_serial" gorm:"column:machine_serial"`
	MachineLabel  string      `json:"machine_label" gorm:"column:machine_label"`
	MachineType   MachineType `json:"machine_type" gorm:"column:machine_type"`
	FinishedAt    *time.Time  `json:"finished_at" gorm:"column:finished_at"`
	IsAvailable   bool        `json:"is_available" gorm:"column:is_available"`
	Weight        int16       `json:"weight" gorm:"column:weight"`
}

type UserGeoLocation struct {
	BranchLat float64 `json:"user_lat" validate:"required"`
	BranchLon float64 `json:"user_lon" validate:"required"`
}
