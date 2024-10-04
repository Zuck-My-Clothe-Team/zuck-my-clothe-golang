package routes

import (
	"zuck-my-clothe/zuck-my-clothe-backend/config"
	"zuck-my-clothe/zuck-my-clothe-backend/controller"
	"zuck-my-clothe/zuck-my-clothe-backend/repository"
	"zuck-my-clothe/zuck-my-clothe-backend/usecases"
)

func AuthRoutes(routeRegister *config.RoutesRegister) {

	userRepository := repository.CreatenewUserRepository(routeRegister.DbConnection)
	userUsecase := usecases.CreateNewUserUsecases(userRepository)
	authRepository := repository.CreateNewAuthenticationRepository(routeRegister.DbConnection)
	authUsecases := usecases.CreateNewAuthenUsecase(authRepository, userRepository)
	authController := controller.CreateNewAuthenController(authUsecases,userUsecase, routeRegister.Config)

	application := routeRegister.Application

	authGroup := application.Group("/auth")
	authGroup.Post("signin", authController.SignIn)
	authGroup.Post("google/callback", authController.GoogleCallback)

}