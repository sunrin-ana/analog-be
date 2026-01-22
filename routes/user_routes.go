package routes

import (
	"analog-be/controller"
	"analog-be/interceptor"

	"github.com/NARUBROWN/spine"
	"github.com/NARUBROWN/spine/pkg/route"
)

func RegisterUserRoutes(app spine.App) {
	app.Route("GET", "/users/search", (*controller.UserController).SearchUser)
	app.Route("GET", "/users/:id", (*controller.UserController).GetUser, route.WithInterceptors(&interceptor.AuthInterceptor{}))

	app.Route("POST", "/users", (*controller.UserController).CreateUser)
	app.Route("PUT", "/users/:id", (*controller.UserController).UpdateUser, route.WithInterceptors(&interceptor.AuthInterceptor{}))
	app.Route("DELETE", "/users/:id", (*controller.UserController).DeleteUser, route.WithInterceptors(&interceptor.AuthInterceptor{}))
}
