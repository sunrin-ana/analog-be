package routes

import (
	"analog-be/controller"
	"github.com/NARUBROWN/spine"
)

func RegisterTopicRoutes(app spine.App) {
	app.Route("GET", "/topic", (*controller.TopicController).GetList)
}
