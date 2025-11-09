package routes

import (
	"github.com/CRS-Project/crs-backend/internal/api/controller"
	"github.com/CRS-Project/crs-backend/internal/middleware"
	"github.com/gin-gonic/gin"
)

func AreaOfConcern(app *gin.Engine, areaOfConcerncontroller controller.AreaOfConcernController, middleware middleware.Middleware) {
	routes := app.Group("/api/v1/area-of-concern-group/:area_of_concern_group_id/area-of-concern")
	{
		routes.POST("", middleware.Authenticate(), areaOfConcerncontroller.Create)
		routes.GET("", middleware.Authenticate(), areaOfConcerncontroller.GetAll)
		routes.GET("/:area_of_concern_id", middleware.Authenticate(), areaOfConcerncontroller.GetById)
		routes.PUT("/:area_of_concern_id", middleware.Authenticate(), areaOfConcerncontroller.Update)
		routes.DELETE("/:area_of_concern_id", middleware.Authenticate(), areaOfConcerncontroller.Delete)
	}
}
