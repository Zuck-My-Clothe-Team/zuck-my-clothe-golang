package usecases

import (
	"zuck-my-clothe/zuck-my-clothe-backend/model"
	"zuck-my-clothe/zuck-my-clothe-backend/repository"
)

type KonCronUsecase interface {
	CleanupExpiredPayment() error
	CleanUpExpiredOrder() error
	CompleteZuckProcess() error
}

type cronUsecase struct {
	paymentRepo     model.PaymentRepository
	orderDetailRepo repository.OrderDetailRepository
}

func CreateNewKonCronUsecase(paymentRepo model.PaymentRepository, orderDetailRepo repository.OrderDetailRepository) KonCronUsecase {
	return &cronUsecase{paymentRepo: paymentRepo,
		orderDetailRepo: orderDetailRepo}
}

func (u *cronUsecase) CleanupExpiredPayment() error {
	return u.paymentRepo.CleanupExpiredPayment()
}

func (u *cronUsecase) CleanUpExpiredOrder() error {
	return u.orderDetailRepo.CleanUpExpiredOrder()
}

func (u *cronUsecase) CompleteZuckProcess() error {
	return u.orderDetailRepo.CompleteZuckProcess()
}
