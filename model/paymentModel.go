package model

import (
	"time"

	"gorm.io/gorm"
)

func (Payments) TableName() string {
	return "Payments"
}

type PaymentStatus string

const (
	Pending PaymentStatus = "Pending"
	Paid    PaymentStatus = "Paid"
	Expired PaymentStatus = "Expired"
	Cancel  PaymentStatus = "Cancel"
)

type Payments struct {
	PaymentID      string         `json:"payment_id" db:"payment_id"`
	Amount         float64        `json:"amount" db:"amount" validate:"required"`
	Payment_Status PaymentStatus  `json:"payment_status" db:"payment_status"`
	DueDate        time.Time      `json:"due_date" db:"due_date"`
	CreatedAt      time.Time      `json:"created_at" db:"created_at"`
	DeletedAt      gorm.DeletedAt `json:"deleted_at,omitempty" db:"deleted_at" swaggertype:"string" example:"null"`
}

type PaymentRepository interface {
	CreatePayment(newPayment Payments) (*Payments, error)
	FindByPaymentID(paymentID string) (*Payments, error)
	UpdatePaymentStatus(paymentID string, status PaymentStatus) error
	CleanupExpiredPayment() error
}

type PaymentUsecase interface {
	CreatePayment(newPayment Payments) (*Payments, error)
	FindByPaymentID(paymentID string) (*Payments, error)
	UpdatePaymentStatus(paymentID string, status PaymentStatus) (*Payments, error)
}
