package routes

import (
	"github.com/CRS-Project/crs-backend/internal/api/controller"
	"github.com/CRS-Project/crs-backend/internal/entity"
	"github.com/CRS-Project/crs-backend/internal/middleware"
	"github.com/gin-gonic/gin"
)

func Package(app *gin.Engine, packagecontroller controller.PackageController, middleware middleware.Middleware) {
	routes := app.Group("/api/v1/package")
	{
		routes.GET("", middleware.Authenticate(), middleware.OnlyAllow(string(entity.RoleSuperAdmin)), packagecontroller.GetAll)
		routes.GET("/:id", middleware.Authenticate(), middleware.OnlyAllow(string(entity.RoleSuperAdmin)), packagecontroller.GetByID)
		routes.POST("", middleware.Authenticate(), middleware.OnlyAllow(string(entity.RoleSuperAdmin)), packagecontroller.CreatePackage)
		routes.PUT("", middleware.Authenticate(), middleware.OnlyAllow(string(entity.RoleSuperAdmin)), packagecontroller.UpdatePackage)
		routes.DELETE("/:id", middleware.Authenticate(), middleware.OnlyAllow(string(entity.RoleSuperAdmin)), packagecontroller.DeletePackage)
	}
}
