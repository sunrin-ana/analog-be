package routes

import (
	"analog-be/controller"
	"analog-be/interceptor"

	"github.com/NARUBROWN/spine"
	"github.com/NARUBROWN/spine/pkg/route"
)

// 여기서 'log' 란 article을 의미합니다.
func RegisterLogRoutes(app spine.App) {
	app.Route("GET", "/logs", (*controller.LogController).GetListOfLog)
	app.Route("GET", "/logs/search", (*controller.LogController).SearchLogs)
	app.Route("GET", "/logs/:id", (*controller.LogController).GetLog)

	app.Route("POST", "/logs", (*controller.LogController).CreateLog, route.WithInterceptors(&interceptor.AuthInterceptor{}))
	app.Route("PUT", "/logs/:id", (*controller.LogController).UpdateLog, route.WithInterceptors(&interceptor.AuthInterceptor{}))
	app.Route("DELETE", "/logs/:id", (*controller.LogController).DeleteLog, route.WithInterceptors(&interceptor.AuthInterceptor{}))

	app.Route("GET", "/logs/:id/", (*controller.LogController).FindAllCommentByLogID)
	app.Route("POST", "/logs/:id/comments", (*controller.LogController).CreateComment, route.WithInterceptors(&interceptor.AuthInterceptor{}))
	app.Route("PUT", "/logs/:id/comments/:commentId", (*controller.LogController).UpdateComment, route.WithInterceptors(&interceptor.AuthInterceptor{}))
	app.Route("DELETE", "/logs/:id/comments/:commentId", (*controller.LogController).DeleteComment, route.WithInterceptors(&interceptor.AuthInterceptor{}))
}
