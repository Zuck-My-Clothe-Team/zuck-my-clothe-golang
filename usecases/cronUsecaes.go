package usecases

import (
	"zuck-my-clothe/zuck-my-clothe-backend/model"
)

type KonCronUsecase interface {
	CleanupExpiredPayment() error
}

type cronUsecase struct {
	paymentRepo model.PaymentRepository
}

func CreateNewKonCronUsecase(paymentRepo model.PaymentRepository) KonCronUsecase {
	return &cronUsecase{paymentRepo: paymentRepo}
}

func (u *cronUsecase) CleanupExpiredPayment() error {
	return u.paymentRepo.CleanupExpiredPayment()
}
