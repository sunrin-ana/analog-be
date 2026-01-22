package routes

import (
	"analog-be/controller"

	"github.com/NARUBROWN/spine"
)

func RegisterAuthRoutes(app spine.App) {
	app.Route("POST", "/auth/login/init", (*controller.AuthController).InitiateLogin)
	app.Route("GET", "/auth/login/callback", (*controller.AuthController).HandleLoginCallback)
	app.Route("POST", "/auth/signup/init", (*controller.AuthController).InitiateSignup)
	app.Route("GET", "/auth/signup/callback", (*controller.AuthController).HandleSignupCallback)

	app.Route("POST", "/auth/refresh", (*controller.AuthController).RefreshToken)

	app.Route("POST", "/auth/logout", (*controller.AuthController).Logout)
	app.Route("GET", "/auth/me", (*controller.AuthController).GetCurrentUser)
}
