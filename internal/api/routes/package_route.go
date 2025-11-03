package routes

import (
	"github.com/CRS-Project/crs-backend/internal/api/controller"
	"github.com/CRS-Project/crs-backend/internal/middleware"
	"github.com/gin-gonic/gin"
)

func Package(app *gin.Engine, packagecontroller controller.PackageController, middleware middleware.Middleware) {
	routes := app.Group("api/v1/package")
	{
		routes.GET("/", packagecontroller.GetPackages)
		routes.POST("/add", packagecontroller.CreatePackage)
		routes.POST("/edit", packagecontroller.UpdatePackage)
		routes.DELETE("/delete", packagecontroller.DeletePackage)
	}
}
