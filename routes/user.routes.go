package routes

import (
	"zuck-my-clothe/zuck-my-clothe-backend/config"
	"zuck-my-clothe/zuck-my-clothe-backend/controller"
	"zuck-my-clothe/zuck-my-clothe-backend/middleware"
	"zuck-my-clothe/zuck-my-clothe-backend/repository"
	"zuck-my-clothe/zuck-my-clothe-backend/usecases"
)

func UserRoutes(routeRegister *config.RoutesRegister) {

	userRepository := repository.CreatenewUserRepository(routeRegister.DbConnection)
	employeeContractRepository := repository.CreateNewEmployeeContractRepository(routeRegister.DbConnection)
	branchRepository := repository.CreateNewBranchRepository(routeRegister.DbConnection)
	machineRepository := repository.CreateMachineRepository(routeRegister.DbConnection)

	userUsecases := usecases.CreateNewUserUsecases(userRepository, employeeContractRepository, branchRepository)
	branchUseCase := usecases.CreateNewBranchUsecase(branchRepository, machineRepository)
	employeeContractUseCase := usecases.CreateNewEmployeeContractUsecase(employeeContractRepository, userRepository)
	userController := controller.CreateNewUserController(userUsecases, routeRegister.Config, employeeContractUseCase, branchUseCase)

	application := routeRegister.Application

	userGroup := application.Group("/users")
	userGroup.Post("/", userController.CreateUser)
	userGroup.Get("/all", middleware.AuthRequire, middleware.IsSuperAdmin, userController.GetAll)
	userGroup.Get("/branch/:branch_id", middleware.AuthRequire, middleware.IsBranchManager, userController.GetBranchEmployee)
	userGroup.Delete("/branch/:branch_id/:id", middleware.AuthRequire, middleware.IsBranchManager, userController.DeleteEmployeeFromBranch)
	userGroup.Get("/manager/all", middleware.AuthRequire, middleware.IsSuperAdmin, userController.GetAllManager)
	userGroup.Get("/:id", middleware.AuthRequire, userController.GetUserById)
	userGroup.Patch("/:id", middleware.AuthRequire, userController.UpdateUser)
	userGroup.Patch("/:id/password", middleware.AuthRequire, userController.UpdateUserPassword)
	userGroup.Delete("/:id", middleware.AuthRequire, middleware.IsSuperAdmin, userController.DeleteUser)

}
