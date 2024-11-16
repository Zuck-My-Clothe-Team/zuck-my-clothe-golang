package routes

import (
	"zuck-my-clothe/zuck-my-clothe-backend/config"
	"zuck-my-clothe/zuck-my-clothe-backend/controller"
	"zuck-my-clothe/zuck-my-clothe-backend/middleware"
	"zuck-my-clothe/zuck-my-clothe-backend/repository"
	"zuck-my-clothe/zuck-my-clothe-backend/usecases"
)

func MachineRoutes(routeRegister *config.RoutesRegister) {
	machineRepo := repository.CreateMachineRepository(routeRegister.DbConnection)
	machineUsecase := usecases.CreateMachineUsecase(machineRepo)
	machineController := controller.CreateMachineController(machineUsecase)

	application := routeRegister.Application

	machineGroup := application.Group("/machine", middleware.AuthRequire)

	machineGroup.Post("/add", middleware.IsBranchManager, machineController.AddMachine)
	machineGroup.Get("/all", middleware.IsSuperAdmin, machineController.GetAll)
	machineGroup.Get("/detail/:serial_id", machineController.GetByMachineSerial)
	machineGroup.Get("/available/branch/:branch_id", machineController.GetAvailableMachineInBranch)
	machineGroup.Get("/branch/:branch_id", machineController.GetByBranchID)
	machineGroup.Put("/update/:serial_id/set_active/:set_active", middleware.IsBranchManager, machineController.UpdateActive)
	machineGroup.Put("/update/:serial_id/set_label/:label", middleware.IsBranchManager, machineController.UpdateLabel)
	machineGroup.Delete("/delete/:serial_id", middleware.IsBranchManager, machineController.SoftDelete)
}
