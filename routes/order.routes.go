package routes

import (
	"zuck-my-clothe/zuck-my-clothe-backend/config"
	"zuck-my-clothe/zuck-my-clothe-backend/controller"
	"zuck-my-clothe/zuck-my-clothe-backend/middleware"
	"zuck-my-clothe/zuck-my-clothe-backend/repository"
	"zuck-my-clothe/zuck-my-clothe-backend/usecases"
)

func OrderRoutes(routeRegister *config.RoutesRegister) {
	orderHeaderRepo := repository.CreateOrderHeaderRepository(routeRegister.DbConnection)
	orderDetailRepo := repository.CreateOrderDetailRepository(routeRegister.DbConnection)
	userRepo := repository.CreatenewUserRepository(routeRegister.DbConnection)

	paymentRepo := repository.CreateNewPaymentRepository(routeRegister.DbConnection)
	paymentUsecase := usecases.CreateNewPaymentUsecase(paymentRepo)

	machineRepo := repository.CreateMachineRepository(routeRegister.DbConnection)

	orderUsecase := usecases.CreateOrderUsecase(orderHeaderRepo, orderDetailRepo, userRepo, machineRepo, paymentUsecase)
	orderController := controller.CreateOrderController(orderUsecase)

	application := routeRegister.Application

	orderGroup := application.Group("/order", middleware.AuthRequire)

	orderGroup.Post("/new", orderController.CreateNewOrder)
	orderGroup.Get("/all", middleware.IsSuperAdmin, orderController.GetAll)
	orderGroup.Get("/branch/:branch_id", middleware.IsEmployee, orderController.GetByBranchID)
	orderGroup.Get("/:order_header_id/:option", orderController.GetByHeaderID)
	orderGroup.Get("/me", orderController.GetByUserID)

	orderGroup.Put("/review", orderController.UpdateReview)
	orderGroup.Put("/update", middleware.IsEmployee, orderController.UpdateStatus)
	orderGroup.Delete("/delete/:order_header_id", middleware.IsBranchManager, orderController.SoftDelete)
}
