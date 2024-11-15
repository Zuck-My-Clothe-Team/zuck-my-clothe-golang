package model

import (
	"time"

	"gorm.io/gorm"
)

func (OrderHeader) TableName() string {
	return "OrderHeaders"
}

type OrderHeader struct {
	OrderHeaderID   string         `json:"order_header_id" gorm:"column:order_header_id;primaryKey"`
	UserID          string         `json:"user_id" gorm:"column:user_id"`
	UserDetail      UserDetailDTO  `json:"user_detail" gorm:"-"`
	BranchID        string         `json:"branch_id" gorm:"column:branch_id"`
	OrderNote       *string        `json:"order_note" gorm:"column:order_note"`
	PaymentID       string         `json:"payment_id" gorm:"column:payment_id"`
	ZuckOnsite      bool           `json:"zuck_onsite" gorm:"column:zuck_onsite"`
	DeliveryAddress *string        `json:"delivery_address" gorm:"column:delivery_address"`
	DeliveryLat     *float64       `json:"delivery_lat" gorm:"column:delivery_lat"`
	DeliveryLong    *float64       `json:"delivery_long" gorm:"column:delivery_long"`
	StarRating      *int16         `json:"star_rating" gorm:"star_rating"`
	ReviewComment   *string        `json:"review_comment" gorm:"review_comment"`
	CreatedAt       time.Time      `json:"created_at" gorm:"column:created_at"`
	CreatedBy       string         `json:"created_by" gorm:"column:created_by"`
	UpdatedAt       time.Time      `json:"updated_at" gorm:"column:updated_at"`
	UpdatedBy       string         `json:"updated_by" gorm:"column:updated_by"`
	DeletedAt       gorm.DeletedAt `json:"deleted_at" gorm:"column:deleted_at;index" swaggertype:"string" example:"null"`
	DeletedBy       *string        `json:"deleted_by" gorm:"column:deleted_by"`
}

type NewOrder struct {
	UserID          string
	BranchID        string           `json:"branch_id" validate:"required"`
	OrderNote       *string          `json:"order_note"`
	ZuckOnsite      bool             `json:"zuck_onsite" validate:"requiredBool"`
	DeliveryAddress *string          `json:"delivery_address"`
	DeliveryLat     *float64         `json:"delivery_lat"`
	DeliveryLong    *float64         `json:"delivery_long"`
	OrderDetails    []NewOrderDetail `json:"order_details" validate:"required"`
}

type FullOrder struct {
	OrderHeaderID   string          `json:"order_header_id"`
	UserID          string          `json:"user_id"`
	UserDetail      UserDetailDTO   `json:"user_detail"`
	BranchID        string          `json:"branch_id"`
	OrderNote       *string         `json:"order_note"`
	PaymentID       string          `json:"payment_id"`
	ZuckOnsite      bool            `json:"zuck_onsite"`
	DeliveryAddress *string         `json:"delivery_address"`
	DeliveryLat     *float64        `json:"delivery_lat"`
	DeliveryLong    *float64        `json:"delivery_long"`
	StarRating      *int16          `json:"star_rating"`
	ReviewComment   *string         `json:"review_comment"`
	CreatedAt       *time.Time      `json:"created_at,omitempty"`
	CreatedBy       *string         `json:"created_by,omitempty"`
	UpdatedAt       *time.Time      `json:"updated_at,omitempty"`
	UpdatedBy       *string         `json:"updated_by,omitempty"`
	DeletedAt       *gorm.DeletedAt `json:"deleted_at,omitempty" swaggertype:"string" example:"null"`
	DeletedBy       *string         `json:"deleted_by,omitempty"`
	OrderDetails    []OrderDetail   `json:"order_details"`
}

type OrderReview struct {
	OrderHeaderID string `json:"order_header_id" validate:"required"`
	UserID        string
	StarRating    *int16  `json:"star_rating" validate:"required,gte=1,lte=5"`
	ReviewComment *string `json:"review_comment"`
}
