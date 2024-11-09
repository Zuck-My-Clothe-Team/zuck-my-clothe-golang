package model

import (
	"time"

	"gorm.io/gorm"
)

func (OrderDetail) TableName() string {
	return "OrderDetails"
}

type OrderStatus string

const (
	Waiting    OrderStatus = "Waiting"
	Processing OrderStatus = "Processing"
	Completed  OrderStatus = "Completed"
	Canceled   OrderStatus = "Canceled"
)

type ServiceType string

const (
	Washing  ServiceType = "Washing"
	Drying   ServiceType = "Drying"
	Delivery ServiceType = "Delivery"
)

type OrderDetail struct {
	OrderBasketID string          `json:"order_basket_id"`
	OrderHeaderID string          `json:"order_header_id"`
	MachineSerial *string         `json:"machine_serial"`
	Weight        int16           `json:"weight"`
	OrderStatus   OrderStatus     `json:"order_status"`
	ServiceType   ServiceType     `json:"service_type"`
	FinishedAt    *time.Time      `json:"finished_at"`
	CreatedAt     *time.Time      `json:"created_at,omitempty"`
	CreatedBy     *string         `json:"created_by,omitempty"`
	UpdatedAt     *time.Time      `json:"updated_at,omitempty"`
	UpdatedBy     *string         `json:"updated_by,omitempty" gorm:"column:updated_by"`
	DeletedAt     *gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"column:deleted_at;index" swaggertype:"string" example:"null"`
	DeletedBy     *string         `json:"deleted_by,omitempty" gorm:"column:deleted_by"`
}

type NewOrderDetail struct {
	MachineSerial *string     `json:"machine_serial"`
	Weight        int16       `json:"weight" validate:"required"`
	ServiceType   ServiceType `json:"service_type" validate:"required"`
}

type UpdateOrder struct {
	OrderBasketID string      `json:"order_basket_id" validate:"required"`
	MachineSerial *string     `json:"machine_serial"`
	OrderStatus   OrderStatus `json:"order_status" validate:"orderStatus"`
	FinishedAt    *time.Time  `json:"finished_at"`
	UpdatedBy     string
}
