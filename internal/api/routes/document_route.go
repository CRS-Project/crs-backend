package routes

import (
	"github.com/CRS-Project/crs-backend/internal/api/controller"
	"github.com/CRS-Project/crs-backend/internal/middleware"
	"github.com/gin-gonic/gin"
)

func Document(app *gin.Engine, documentcontroller controller.DocumentController, middleware middleware.Middleware) {
	routes := app.Group("/api/v1/document")
	{
		routes.GET("", middleware.Authenticate(), documentcontroller.GetAll)
		routes.GET("/:id", middleware.Authenticate(), documentcontroller.GetByID)
		routes.POST("", middleware.Authenticate(), documentcontroller.CreateDocument)
		routes.PUT("", middleware.Authenticate(), documentcontroller.Update)
		routes.DELETE("/:id", middleware.Authenticate(), documentcontroller.DeleteDocument)
	}
}
