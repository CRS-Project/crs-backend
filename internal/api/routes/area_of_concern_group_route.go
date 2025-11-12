package routes

import (
	"github.com/CRS-Project/crs-backend/internal/api/controller"
	"github.com/CRS-Project/crs-backend/internal/middleware"
	"github.com/gin-gonic/gin"
)

func AreaOfConcernGroup(app *gin.Engine, areaOfConcernGroupcontroller controller.AreaOfConcernGroupController, middleware middleware.Middleware) {
	routes := app.Group("/api/v1/area-of-concern-group")
	{
		routes.POST("", middleware.Authenticate(), areaOfConcernGroupcontroller.Create)
		routes.GET("", middleware.Authenticate(), areaOfConcernGroupcontroller.GetAll)
		routes.GET("/statistic/:package_id", middleware.Authenticate(), areaOfConcernGroupcontroller.Statistic)
		routes.GET("/:area_of_concern_group_id/generate-pdf", middleware.Authenticate(), areaOfConcernGroupcontroller.GeneratePDF)
		routes.GET("/:area_of_concern_group_id", middleware.Authenticate(), areaOfConcernGroupcontroller.GetById)
		routes.PUT("/:area_of_concern_group_id", middleware.Authenticate(), areaOfConcernGroupcontroller.Update)
		routes.DELETE("/:area_of_concern_group_id", middleware.Authenticate(), areaOfConcernGroupcontroller.Delete)
	}
}
