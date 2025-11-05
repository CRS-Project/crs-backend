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
		routes.POST("", middleware.Authenticate(), documentcontroller.CreateDocument)
		routes.DELETE("", middleware.Authenticate(), documentcontroller.DeleteDocument)
	}
}
