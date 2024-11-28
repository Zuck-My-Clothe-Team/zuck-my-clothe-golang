package repository

import (
	"time"
	"zuck-my-clothe/zuck-my-clothe-backend/model"
	"zuck-my-clothe/zuck-my-clothe-backend/platform"

	"gorm.io/gorm"
)

type orderDetailRepository struct {
	db *platform.Postgres
}

type OrderDetailRepository interface {
	GetAll() (*[]model.OrderDetail, error)
	CreateOrderDetails(orderDetails *[]model.OrderDetail) (*[]model.OrderDetail, error)
	GetByHeaderID(orderHeaderID string, isAdminView bool) (*[]model.OrderDetail, error)
	UpdateStatus(order model.OrderDetail) (*model.OrderDetail, error)
	GetDetail(orderBasketID string) (*model.OrderDetail, error)
	DeleteByHeaderID(orderHeaderID string, deletedBy string) (*[]model.OrderDetail, error)
	CleanUpExpiredOrder() error
	CompleteZuckProcess() error
}

func CreateOrderDetailRepository(db *platform.Postgres) OrderDetailRepository {
	return &orderDetailRepository{db: db}
}

func (u *orderDetailRepository) GetAll() (*[]model.OrderDetail, error) {
	orderDetails := new([]model.OrderDetail)

	result := u.db.Find(orderDetails)

	if result.Error != nil {
		return nil, result.Error
	}

	return orderDetails, result.Error
}

func (u *orderDetailRepository) CreateOrderDetails(orderDetails *[]model.OrderDetail) (*[]model.OrderDetail, error) {
	result := u.db.CreateInBatches(orderDetails, len(*orderDetails))

	if result.Error != nil {
		return nil, result.Error
	}

	newOrder, err := u.GetByHeaderID((*orderDetails)[0].OrderHeaderID, false)

	if err != nil {
		return nil, err
	}

	return newOrder, err
}

func (u *orderDetailRepository) GetByHeaderID(orderHeaderID string, isAdminView bool) (*[]model.OrderDetail, error) {
	orders := new([]model.OrderDetail)

	var result *gorm.DB

	if isAdminView {
		result = u.db.Where("order_header_id = ?", orderHeaderID).Find(&orders)
	} else {
		result = u.db.Select("order_basket_id", "order_header_id", "machine_serial", "weight", "order_status", "service_type", "finished_at", "created_at", "created_by", "updated_at", "updated_by").Where("order_header_id = ?", orderHeaderID).Find(&orders)
	}

	if result.Error != nil {
		return nil, result.Error
	}

	return orders, result.Error
}

func (u *orderDetailRepository) GetDetail(orderBasketID string) (*model.OrderDetail, error) {
	orderDetail := new(model.OrderDetail)

	result := u.db.Where("order_basket_id = ?", orderBasketID).Find(&orderDetail)

	if result.Error != nil {
		return nil, result.Error
	}

	return orderDetail, result.Error
}

func (u *orderDetailRepository) UpdateStatus(order model.OrderDetail) (*model.OrderDetail, error) {
	updatedOrder := new(model.OrderDetail)
	result := u.db.
		Where("order_basket_id = ?", order.OrderBasketID).
		Updates(order).
		Find(&updatedOrder)

	if result.RowsAffected == 0 {
		return nil, gorm.ErrRecordNotFound
	}

	if result.Error != nil {
		return nil, result.Error
	}

	return updatedOrder, result.Error
}

func (u *orderDetailRepository) DeleteByHeaderID(orderHeaderID string, deletedBy string) (*[]model.OrderDetail, error) {
	orderDetails := new([]model.OrderDetail)

	result := u.db.Model(model.OrderDetail{}).
		Where("order_header_id = ?", orderHeaderID).
		Update("deleted_by", deletedBy).
		Find(&orderDetails)

	if result.Error != nil {
		return nil, result.Error
	}

	cascading, err := u.GetByHeaderID(orderHeaderID, true)

	if err != nil {
		return nil, err
	}

	result = u.db.
		Where("order_header_id = ?", orderHeaderID).
		Delete(cascading)

	if result.RowsAffected == 0 {
		return nil, gorm.ErrRecordNotFound
	}

	if result.Error != nil {
		return nil, result.Error
	}

	return cascading, result.Error
}

func (u *orderDetailRepository) CleanUpExpiredOrder() error {
	orderDetailDummy := new(model.OrderDetail)
	dbTx := u.db.Raw(`
	UPDATE "OrderDetails"
	SET order_status = 'Expired'
	WHERE order_basket_id IN (
		SELECT OD.order_basket_id
		FROM "OrderDetails" AS OD INNER JOIN "OrderHeaders" AS OH ON OD.order_header_id = OH.order_header_id
		INNER JOIN "Payments" AS PM ON PM.payment_id = OH.payment_id
		WHERE OD.order_status = 'Waiting' AND PM.payment_status = 'Expired' AND PM.due_date < $1);
	`, time.Now().UTC()).Scan(orderDetailDummy)
	return dbTx.Error
}

func (u *orderDetailRepository) CompleteZuckProcess() error {
	orderDetailDummy := new(model.OrderDetail)
	dbTx := u.db.Raw(`
	UPDATE "OrderDetails"
	SET order_status = 'Completed'
	WHERE finished_at < $1 AND order_status = 'Processing';	
	`, time.Now().UTC()).Scan(orderDetailDummy)
	return dbTx.Error
}
