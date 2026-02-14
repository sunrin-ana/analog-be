package routes

import (
	"analog-be/controller"

	"github.com/NARUBROWN/spine"
)

func RegisterFeedRoutes(app spine.App) {
	app.Route("GET", "/feed", (*controller.FeedController).GetFeed)
	app.Route("GET", "/sitemaps/:file", (*controller.FeedController).GetSitemap)
}
