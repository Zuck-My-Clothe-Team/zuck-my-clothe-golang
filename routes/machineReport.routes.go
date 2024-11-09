package routes

import (
	"zuck-my-clothe/zuck-my-clothe-backend/config"
	"zuck-my-clothe/zuck-my-clothe-backend/controller"
	"zuck-my-clothe/zuck-my-clothe-backend/middleware"
	"zuck-my-clothe/zuck-my-clothe-backend/repository"
	"zuck-my-clothe/zuck-my-clothe-backend/usecases"
)

func MachineReportRoutes(routeRegister *config.RoutesRegister) {
	brachRepo := repository.CreateNewBranchRepository(routeRegister.DbConnection)
	contractRepo := repository.CreateNewEmployeeContractRepository(routeRegister.DbConnection)
	machineReportRepo := repository.CreateNewMachineReportRepository(routeRegister.DbConnection)
	machineRepo := repository.CreateMachineRepository(routeRegister.DbConnection)
	machineReportUsecase := usecases.CreateNewMachineReportUsecase(machineReportRepo, machineRepo, brachRepo, contractRepo)
	machineReportController := controller.CreateNewMachineReportController(machineReportUsecase)

	application := routeRegister.Application
	machineReportGroup := application.Group("/report", middleware.AuthRequire)
	machineReportGroup.Get("/", middleware.IsSuperAdmin, machineReportController.GetAll)
	machineReportGroup.Post("/add", machineReportController.CreateMachineReport)
	machineReportGroup.Get("/user", machineReportController.FindMachineReportByUserID)
	machineReportGroup.Put("/update", machineReportController.UpdateMachineReportStatus)
	machineReportGroup.Get("/branch/:branchID", middleware.IsBranchMember, machineReportController.FindMachineReportByBranch)
	machineReportGroup.Delete("/delete/:reportID", machineReportController.DeleteMachineReport)
}
