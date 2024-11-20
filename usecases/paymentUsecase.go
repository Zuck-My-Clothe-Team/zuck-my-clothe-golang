package usecases

import (
	"errors"
	"time"
	"zuck-my-clothe/zuck-my-clothe-backend/model"

	"github.com/google/uuid"
)

type paymentUsecase struct {
	paymentRepository model.PaymentRepository
}

func CreateNewPaymentUsecase(paymentRepository model.PaymentRepository) model.PaymentUsecase {
	return &paymentUsecase{paymentRepository: paymentRepository}
}

func (u *paymentUsecase) CreatePayment(newPayment model.Payments) (*model.Payments, error) {
	data := model.Payments{
		PaymentID:      uuid.New().String(),
		Amount:         newPayment.Amount,
		Payment_Status: "Pending",
		DueDate:        time.Now().UTC().Add(time.Minute * 10),
		CreatedAt:      time.Now().UTC(),
	}
	response, err := u.paymentRepository.CreatePayment(data)
	if err != nil {
		return nil, err
	}
	return response, nil
}

func (u *paymentUsecase) FindByPaymentID(paymentID string) (*model.Payments, error) {
	data, err := u.paymentRepository.FindByPaymentID(paymentID)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (u *paymentUsecase) UpdatePaymentStatus(paymentID string, status model.PaymentStatus) (*model.Payments, error) {
	var response *model.Payments = nil
	var err error

	response, err = u.paymentRepository.FindByPaymentID(paymentID)
	if err != nil {
		return nil, err
	}
	if response.Payment_Status != "Pending" {
		return nil, errors.New("err: cannot update payment that not in \"Pending\" State")
	}
	err = u.paymentRepository.UpdatePaymentStatus(paymentID, status)
	if err != nil {
		return nil, err
	}
	response, err = u.paymentRepository.FindByPaymentID(paymentID)
	if err != nil {
		return nil, err
	}
	return response, nil
}
