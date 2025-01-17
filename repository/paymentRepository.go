package repository

import (
	"time"
	"zuck-my-clothe/zuck-my-clothe-backend/model"
	"zuck-my-clothe/zuck-my-clothe-backend/platform"
)

type paymentReopository struct {
	db *platform.Postgres
}

func CreateNewPaymentRepository(db *platform.Postgres) model.PaymentRepository {
	return &paymentReopository{db: db}
}

func (u *paymentReopository) CreatePayment(newPayment model.Payments) (*model.Payments, error) {
	createdPayment := new(model.Payments)
	dbTx := u.db.Create(newPayment)
	if dbTx.Error != nil {
		return nil, dbTx.Error
	}
	if dbTx := u.db.First(createdPayment, "payment_id = ?", newPayment.PaymentID); dbTx.Error != nil {
		return nil, dbTx.Error
	}
	return createdPayment, nil
}

func (u *paymentReopository) FindByPaymentID(paymentID string) (*model.Payments, error) {
	data := new(model.Payments)
	dbTx := u.db.First(data, "payment_id = ?", paymentID)
	if dbTx.Error != nil {
		return nil, dbTx.Error
	}
	return data, nil
}

func (u *paymentReopository) UpdatePaymentStatus(paymentID string, status model.PaymentStatus) error {
	var payment *model.Payments
	dbTx := u.db.Debug().Raw(`
	UPDATE "Payments"
	SET "payment_status" = $1
	WHERE "payment_id" = $2 AND "payment_status" = 'Pending' AND "due_date" > $3;`, status, paymentID, time.Now().UTC()).Scan(payment)
	return dbTx.Error
}

func (u *paymentReopository) CleanupExpiredPayment() error {
	var list []model.Payments
	dbTx := u.db.Raw(`
	UPDATE "Payments"
	SET payment_status = 'Expired'
	WHERE due_date < $1 AND payment_status = 'Pending'`, time.Now().UTC()).Scan(&list)
	return dbTx.Error
}
