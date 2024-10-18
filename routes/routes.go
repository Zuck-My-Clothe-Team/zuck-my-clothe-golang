package routes

import "zuck-my-clothe/zuck-my-clothe-backend/config"

func RoutesRegister(routeRegister *config.RoutesRegister) {
	UserRoutes(routeRegister)
	AuthRoutes(routeRegister)
	BranchRoutes(routeRegister)
}
