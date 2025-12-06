package routes

import (
	"github.com/CRS-Project/crs-backend/internal/api/controller"
	"github.com/CRS-Project/crs-backend/internal/middleware"
	"github.com/gin-gonic/gin"
)

func DisciplineGroup(app *gin.Engine, areaOfConcernGroupcontroller controller.DisciplineGroupController, middleware middleware.Middleware) {
	routes := app.Group("/api/v1/discipline-group")
	{
		routes.POST("", middleware.Authenticate(), areaOfConcernGroupcontroller.Create)
		routes.GET("", middleware.Authenticate(), areaOfConcernGroupcontroller.GetAll)
		routes.GET("/statistic/:package_id", middleware.Authenticate(), areaOfConcernGroupcontroller.Statistic)
		routes.GET("/:discipline_group_id/generate-pdf", middleware.Authenticate(), areaOfConcernGroupcontroller.GeneratePDF)
		routes.GET("/:discipline_group_id/consolidator", middleware.Authenticate(), areaOfConcernGroupcontroller.GetAllConsolidator)
		routes.GET("/:discipline_group_id", middleware.Authenticate(), areaOfConcernGroupcontroller.GetById)
		routes.PUT("/:discipline_group_id", middleware.Authenticate(), areaOfConcernGroupcontroller.Update)
		routes.DELETE("/:discipline_group_id", middleware.Authenticate(), areaOfConcernGroupcontroller.Delete)
	}
}
