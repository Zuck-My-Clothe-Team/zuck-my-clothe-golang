package routes

import (
	"zuck-my-clothe/zuck-my-clothe-backend/config"
	"zuck-my-clothe/zuck-my-clothe-backend/controller"
	"zuck-my-clothe/zuck-my-clothe-backend/middleware"
	"zuck-my-clothe/zuck-my-clothe-backend/repository"
	"zuck-my-clothe/zuck-my-clothe-backend/usecases"
)

func BranchRoutes(routeRegister *config.RoutesRegister) {
	branchRepo := repository.CreateNewBranchRepository(routeRegister.DbConnection)
	branchUsecase := usecases.CreateNewBranchUsecase(branchRepo)
	branchController := controller.CreateNewBranchController(branchUsecase)

	application := routeRegister.Application

	branchGroup := application.Group("/branch", middleware.AuthRequire)

	branchGroup.Post("/create", middleware.AuthRequire, middleware.IsSuperAdmin, branchController.CreateBranch)
	branchGroup.Get("/all", branchController.GetAll)
	branchGroup.Get("/id/:id", branchController.GetByBranchID)
	branchGroup.Get("/owns", middleware.AuthRequire, middleware.IsBranchManager, branchController.GetByBranchOwner)
	branchGroup.Put("/update", branchController.UpdateBranch)
	branchGroup.Delete("/delete/:id", middleware.AuthRequire, middleware.IsSuperAdmin, branchController.DeleteBranch)
}
