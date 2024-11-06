package routes

import (
	"zuck-my-clothe/zuck-my-clothe-backend/config"
	"zuck-my-clothe/zuck-my-clothe-backend/controller"
	"zuck-my-clothe/zuck-my-clothe-backend/middleware"
	"zuck-my-clothe/zuck-my-clothe-backend/repository"
	"zuck-my-clothe/zuck-my-clothe-backend/usecases"
)

func EmployeeContractRoutes(routeRegister *config.RoutesRegister) {

	employeeContractRepository := repository.CreateNewEmployeeContractRepository(routeRegister.DbConnection)
	userRepository := repository.CreatenewUserRepository(routeRegister.DbConnection)
	employeeContractUsecases := usecases.CreateNewEmployeeContractUsecase(employeeContractRepository, userRepository)
	employeeContractController := controller.CreateNewEmployeeContractController(employeeContractUsecases)

	application := routeRegister.Application
	employeeContractGroup := application.Group("/employee-contract", middleware.AuthRequire)

	employeeContractGroup.Get("/", middleware.IsSuperAdmin, employeeContractController.GetAll)
	employeeContractGroup.Post("/", middleware.IsBranchManager, employeeContractController.CreateEmployeeContract)
	employeeContractGroup.Delete("/:contract_id", middleware.IsBranchManager, employeeContractController.SoftDelete)
	employeeContractGroup.Get("/branch/:branch_id", middleware.IsBranchManager, employeeContractController.GetByBranchID)
	employeeContractGroup.Get("/user/:user_id", employeeContractController.GetByUserID)
}
