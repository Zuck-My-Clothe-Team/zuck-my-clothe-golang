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

	branchGroup.Post("/create", middleware.IsSuperAdmin, branchController.CreateBranch)
	branchGroup.Get("/all", branchController.GetAll)
	branchGroup.Post("/closest-to-me", branchController.GetClosestToMe)
	branchGroup.Get("/owner", middleware.IsBranchManager, branchController.GetByBranchOwner)
	branchGroup.Get("/:id", branchController.GetByBranchID)

	branchGroup.Put("/update", middleware.IsBranchManager, branchController.UpdateBranch)
	branchGroup.Delete("/:id", middleware.IsSuperAdmin, branchController.DeleteBranch)
}
