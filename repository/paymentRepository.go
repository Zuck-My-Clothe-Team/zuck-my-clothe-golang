package repository

import (
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

func (u *paymentReopository) CleanupExpiredPayment() error {
	list := new(model.Payments)
	dbTx := u.db.Debug().Raw(`
	UPDATE "Payments"
	SET "Payments".payment_status = 'Cancel'
	WHERE "Payments".due_date < NOW() AND "Payments".payment_status = 'Pending';`).Scan(list)
	return dbTx.Error
}
