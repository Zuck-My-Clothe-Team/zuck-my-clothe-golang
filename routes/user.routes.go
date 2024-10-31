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
	userUsecases := usecases.CreateNewUserUsecases(userRepository)
	userController := controller.CreateNewUserController(userUsecases, routeRegister.Config)

	application := routeRegister.Application

	userGroup := application.Group("/users")
	userGroup.Get("/all", middleware.AuthRequire, middleware.IsBranchManager, userController.GetAll)
	userGroup.Post("/register", userController.CreateUser)

}
