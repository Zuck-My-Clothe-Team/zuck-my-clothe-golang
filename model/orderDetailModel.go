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
	Waiting      OrderStatus = "Waiting"
	Processing   OrderStatus = "Processing"
	Completed    OrderStatus = "Completed"
	Canceled     OrderStatus = "Canceled"
	OrderExpired OrderStatus = "Expired"
)

type ServiceType string

const (
	Washing  ServiceType = "Washing"
	Drying   ServiceType = "Drying"
	Pickup   ServiceType = "Pickup"
	Delivery ServiceType = "Delivery"
	Agents   ServiceType = "Agents"
)

const (
	SevenKGMachinePrice     int = 50
	FourteenKGMachinePrice  int = 100
	TwentyOneKGMachinePrice int = 150
	DeliveryPrice           int = 20
	PickupPrice             int = 20
	AgentsPrice             int = 20
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
	Weight        int16       `json:"weight"`
	ServiceType   ServiceType `json:"service_type"`
}

type UpdateOrder struct {
	OrderBasketID string      `json:"order_basket_id" validate:"required"`
	MachineSerial *string     `json:"machine_serial"`
	OrderStatus   OrderStatus `json:"order_status" validate:"orderStatus"`
	FinishedAt    *time.Time  `json:"finished_at"`
	UpdatedBy     string
}
