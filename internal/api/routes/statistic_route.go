package routes

import (
	"github.com/CRS-Project/crs-backend/internal/api/controller"
	"github.com/CRS-Project/crs-backend/internal/middleware"
	"github.com/gin-gonic/gin"
)

func Statistic(app *gin.Engine, statisticcontroller controller.StatisticController, middleware middleware.Middleware) {
	routes := app.Group("/api/v1/statistic")
	{
		routes.GET("/aoc-comment-chart/:package_id", statisticcontroller.GetAOCAndCommentChart)
		routes.GET("/aoc-comment-card/:package_id", statisticcontroller.GetCommentCard)
		routes.GET("/comment-user-chart/:package_id", statisticcontroller.GetCommentUserChart)
	}
}
