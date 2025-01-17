package repository

import (
	"errors"
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
	GetByBranchID(branchID string, status string) (*[]model.OrderHeader, error)
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

	if result.RowsAffected == 0 {
		return &model.OrderHeader{}, nil
	}

	if result.Error != nil {
		return &model.OrderHeader{}, result.Error
	}

	return order, result.Error
}

func (u *orderHeaderRepository) GetByBranchID(branchID string, status string) (*[]model.OrderHeader, error) {
	order := new([]model.OrderHeader)

	var result *gorm.DB

	if status == "" {
		result = u.db.Raw(`
		SELECT OH.order_header_id, OH.user_id, OH.branch_id, OH.order_note, OH.payment_id, OH.zuck_onsite, OH.delivery_address, OH.delivery_lat,
		OH.delivery_long, OH.star_rating, OH.review_comment, OH.created_at, OH.created_by, OH.updated_at, OH.updated_by, OH.deleted_at, OH.deleted_by
		FROM "OrderHeaders" AS OH INNER JOIN "Payments" AS PM ON OH.payment_id = PM.payment_id
		WHERE OH.branch_id = $1;`, branchID).Scan(&order)
	} else {
		result = u.db.Raw(`
		SELECT OH.order_header_id, OH.user_id, OH.branch_id, OH.order_note, OH.payment_id, OH.zuck_onsite, OH.delivery_address, OH.delivery_lat,
		OH.delivery_long, OH.star_rating, OH.review_comment, OH.created_at, OH.created_by, OH.updated_at, OH.updated_by, OH.deleted_at, OH.deleted_by
		FROM "OrderHeaders" AS OH INNER JOIN "Payments" AS PM ON OH.payment_id = PM.payment_id
		WHERE OH.branch_id = $1 AND PM.payment_status = $2;`, branchID, status).Scan(&order)
	}

	if result == nil {
		return nil, errors.New("ERR: unable to initialize query")
	}

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

	if result.RowsAffected == 0 {
		return new([]model.OrderHeader), nil
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
