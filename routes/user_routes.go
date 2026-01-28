package routes

import (
	"analog-be/controller"
	"analog-be/interceptor"

	"github.com/NARUBROWN/spine"
	"github.com/NARUBROWN/spine/pkg/route"
)

func RegisterUserRoutes(app spine.App) {
	app.Route("GET", "/users/search/list", (*controller.UserController).Search)
	app.Route("GET", "/users/:id", (*controller.UserController).Get, route.WithInterceptors(&interceptor.AuthInterceptor{}))

	app.Route("POST", "/users", (*controller.UserController).Create)
	app.Route("PUT", "/users", (*controller.UserController).Update, route.WithInterceptors(&interceptor.AuthInterceptor{}))
	app.Route("DELETE", "/users", (*controller.UserController).Delete, route.WithInterceptors(&interceptor.AuthInterceptor{}))
}
