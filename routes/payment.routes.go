package routes

import (
	"zuck-my-clothe/zuck-my-clothe-backend/config"
	"zuck-my-clothe/zuck-my-clothe-backend/controller"
	"zuck-my-clothe/zuck-my-clothe-backend/middleware"
	"zuck-my-clothe/zuck-my-clothe-backend/repository"
	"zuck-my-clothe/zuck-my-clothe-backend/usecases"
)

func PaymentRoutes(routeRegister *config.RoutesRegister) {
	paymentRepo := repository.CreateNewPaymentRepository(routeRegister.DbConnection)
	paymentUsecase := usecases.CreateNewPaymentUsecase(paymentRepo)
	paymentController := controller.CreateNewPaymentController(paymentUsecase)

	application := routeRegister.Application
	paymentGroup := application.Group("/payment", middleware.AuthRequire)
	paymentGroup.Post("/add", paymentController.CreatePayment)
	paymentGroup.Get("/detail/:paymentID", paymentController.FindByPaymentID)
	paymentGroup.Put("/update/:paymentID/setstatus/:status", paymentController.UpdatePaymenstatus)
}
