package routes

import (
	"zuck-my-clothe/zuck-my-clothe-backend/config"
	"zuck-my-clothe/zuck-my-clothe-backend/controller"
	"zuck-my-clothe/zuck-my-clothe-backend/middleware"
	"zuck-my-clothe/zuck-my-clothe-backend/repository"
	"zuck-my-clothe/zuck-my-clothe-backend/usecases"
)

func AuthRoutes(routeRegister *config.RoutesRegister) {
	userRepository := repository.CreatenewUserRepository(routeRegister.DbConnection)
	employeeContractRepository := repository.CreateNewEmployeeContractRepository(routeRegister.DbConnection)

	userUsecase := usecases.CreateNewUserUsecases(userRepository, employeeContractRepository)
	authRepository := repository.CreateNewAuthenticationRepository(routeRegister.DbConnection)
	authUsecases := usecases.CreateNewAuthenUsecase(authRepository, userRepository)
	authController := controller.CreateNewAuthenController(authUsecases, userUsecase, routeRegister.Config)

	application := routeRegister.Application

	authGroup := application.Group("/auth")
	authGroup.Post("signin", authController.SignIn)
	authGroup.Get("me", middleware.AuthRequire, authController.Me)
	authGroup.Post("google/callback", authController.GoogleCallback)

}
