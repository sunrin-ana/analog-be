package routes

import (
	"analog-be/controller"

	"github.com/NARUBROWN/spine"
)

func RegisterHealthRoutes(app spine.App) {
	app.Route("GET", "/health", (*controller.HealthController).Health)

	app.Route("GET", "/health/ready", (*controller.HealthController).Ready)

	app.Route("GET", "/health/live", (*controller.HealthController).Live)
}
