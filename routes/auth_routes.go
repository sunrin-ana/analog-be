package routes

import (
	"analog-be/controller"

	"github.com/NARUBROWN/spine"
)

func RegisterAuthRoutes(app spine.App) {
	app.Route("GET", "/auth/callback", (*controller.AuthController).HandleAuthCallback)

	app.Route("POST", "/auth/refresh", (*controller.AuthController).RefreshToken)
	app.Route("POST", "/auth/logout", (*controller.AuthController).Logout)
	app.Route("GET", "/auth/me", (*controller.AuthController).GetCurrentUser)
}
