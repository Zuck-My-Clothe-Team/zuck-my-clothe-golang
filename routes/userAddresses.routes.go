package routes

import (
	"zuck-my-clothe/zuck-my-clothe-backend/config"
	"zuck-my-clothe/zuck-my-clothe-backend/controller"
	"zuck-my-clothe/zuck-my-clothe-backend/middleware"
	"zuck-my-clothe/zuck-my-clothe-backend/repository"
	"zuck-my-clothe/zuck-my-clothe-backend/usecases"
)

func UserAddressesRoutes(routeRegister *config.RoutesRegister) {
	userAddressesRepo := repository.CreateNewUserAddressesRepository(routeRegister.DbConnection)
	userAddressesUsecase := usecases.CreateNewUserAddressesUsecase(userAddressesRepo)
	userAddressesController := controller.CreateNewUserAddressesController(userAddressesUsecase)

	application := routeRegister.Application
	userAddressesGroup := application.Group("/address", middleware.AuthRequire)
	userAddressesGroup.Post("/add", userAddressesController.AddUserAddress)
	//userAddressesGroup.Get("/detail/aid/:addressID", userAddressesController.FindByAddressID)
	userAddressesGroup.Get("/detail/owner", userAddressesController.FindUserAddresByOwnerID)
	userAddressesGroup.Put("/update", userAddressesController.UpdateUserAddressData)
	userAddressesGroup.Delete("/delete/:addressID", userAddressesController.DeleteUserAddress)
}
