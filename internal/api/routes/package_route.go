package routes

import (
	"github.com/CRS-Project/crs-backend/internal/api/controller"
	"github.com/CRS-Project/crs-backend/internal/middleware"
	"github.com/gin-gonic/gin"
)

func Package(app *gin.Engine, packagecontroller controller.PackageController, middleware middleware.Middleware) {
	routes := app.Group("/api/v1/package")
	{
		routes.GET("", packagecontroller.GetAll)
		routes.GET("/:id", packagecontroller.GetByID)
		routes.GET("/:id/generate-pdf", middleware.Authenticate(), packagecontroller.GeneratePDF)
		routes.GET("/:id/generate-excel", middleware.Authenticate(), packagecontroller.GenerateExcel)
		routes.GET("/me", middleware.Authenticate(), packagecontroller.GetAllByUser)
		routes.POST("", packagecontroller.CreatePackage)
		routes.PUT("", packagecontroller.UpdatePackage)
		routes.DELETE("/:id", packagecontroller.DeletePackage)
	}
}
