package routes

import (
	"analog-be/controller"
	"analog-be/interceptor"

	"github.com/NARUBROWN/spine"
	"github.com/NARUBROWN/spine/pkg/route"
)

func RegisterAuthRoutes(app spine.App) {
	app.Route("GET", "/auth/callback", (*controller.AuthController).HandleAuthCallback)

	app.Route("PUT", "/auth/token", (*controller.AuthController).RefreshToken)
	app.Route("DELETE", "/auth/token", (*controller.AuthController).Logout, route.WithInterceptors(&interceptor.AuthInterceptor{}))
}
