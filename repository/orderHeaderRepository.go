package repository

import (
	"zuck-my-clothe/zuck-my-clothe-backend/model"
	"zuck-my-clothe/zuck-my-clothe-backend/platform"

	"gorm.io/gorm"
)

type orderHeaderRepository struct {
	db *platform.Postgres
}

type OrderHeaderRepository interface {
	CreateOrderHeader(orderHeader *model.OrderHeader) (*model.OrderHeader, error)
	GetAll() (*[]model.OrderHeader, error)
	GetByID(orderHeaderID string, isAdminView bool) (*model.OrderHeader, error)
	GetByBranchID(branchID string) (*[]model.OrderHeader, error)
	GetByUserID(userID string) (*[]model.OrderHeader, error)
	UpdateReview(order model.OrderHeader) (*model.OrderHeader, error)
	SoftDelete(orderHeaderID string, deletedBy string) (*model.OrderHeader, error)
}

func CreateOrderHeaderRepository(db *platform.Postgres) OrderHeaderRepository {
	return &orderHeaderRepository{db: db}
}

func (u *orderHeaderRepository) GetAll() (*[]model.OrderHeader, error) {
	orderHeaders := new([]model.OrderHeader)

	result := u.db.Find(orderHeaders)

	if result.Error != nil {
		return nil, result.Error
	}

	return orderHeaders, result.Error
}

func (u *orderHeaderRepository) CreateOrderHeader(orderHeader *model.OrderHeader) (*model.OrderHeader, error) {
	result := u.db.Create(orderHeader)

	if result.Error != nil {
		return nil, result.Error
	}

	newOrder, err := u.GetByID(orderHeader.OrderHeaderID, false)

	if err != nil {
		return nil, err
	}

	return newOrder, err
}

func (u *orderHeaderRepository) GetByID(orderHeaderID string, isAdminView bool) (*model.OrderHeader, error) {
	order := new(model.OrderHeader)

	var result *gorm.DB

	if isAdminView {
		result = u.db.Where("order_header_id = ?", orderHeaderID).Find(&order)
	} else {
		result = u.db.Select("order_header_id", "user_id", "branch_id", "order_note", "payment_id", "zuck_onsite", "delivery_address", "delivery_lat", "delivery_long", "created_at", "updated_at").Where("order_header_id = ?", orderHeaderID).Find(&order)
	}

	if result.Error != nil {
		return nil, result.Error
	}

	return order, result.Error
}

func (u *orderHeaderRepository) GetByBranchID(branchID string) (*[]model.OrderHeader, error) {
	order := new([]model.OrderHeader)

	result := u.db.Where("branch_id = ?", branchID).Find(&order)

	if result.Error != nil {
		return nil, result.Error
	}

	return order, result.Error
}

func (u *orderHeaderRepository) GetByUserID(userID string) (*[]model.OrderHeader, error) {
	order := new([]model.OrderHeader)

	result := u.db.Where("user_id = ?", userID).Find(&order)

	if result.Error != nil {
		return nil, result.Error
	}

	return order, result.Error
}

func (u *orderHeaderRepository) UpdateReview(order model.OrderHeader) (*model.OrderHeader, error) {
	updatedOrder := new(model.OrderHeader)

	result := u.db.
		Model(model.OrderHeader{}).
		Where("order_header_id = ?", order.OrderHeaderID).
		Updates(order).
		Find(updatedOrder)

	if result.RowsAffected == 0 {
		return nil, gorm.ErrRecordNotFound
	}

	if result.Error != nil {
		return nil, result.Error
	}

	return updatedOrder, result.Error
}

func (u *orderHeaderRepository) SoftDelete(orderHeaderID string, deletedBy string) (*model.OrderHeader, error) {
	order := new(model.OrderHeader)
	order.OrderHeaderID = orderHeaderID

	result := u.db.Model(model.OrderHeader{}).
		Where("order_header_id = ?", orderHeaderID).
		Update("deleted_by", deletedBy).
		Find(&order)

	if result.Error != nil {
		return nil, result.Error
	}

	result = u.db.Delete(&order)

	if result.RowsAffected == 0 {
		return nil, gorm.ErrRecordNotFound
	}

	if result.Error != nil {
		return nil, result.Error
	}

	return order, result.Error
}
