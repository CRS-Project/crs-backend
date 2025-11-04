package routes

import (
	"github.com/CRS-Project/crs-backend/internal/api/controller"
	"github.com/CRS-Project/crs-backend/internal/entity"
	"github.com/CRS-Project/crs-backend/internal/middleware"
	"github.com/gin-gonic/gin"
)

func User(app *gin.Engine, usercontroller controller.UserController, middleware middleware.Middleware) {
	routes := app.Group("/api/v1/user")
	{
		routes.POST("", middleware.Authenticate(), middleware.OnlyAllow(string(entity.RoleSuperAdmin)), usercontroller.Create)
		routes.GET("", middleware.Authenticate(), middleware.OnlyAllow(string(entity.RoleSuperAdmin)), usercontroller.GetAll)
		routes.GET("/:id", middleware.Authenticate(), middleware.OnlyAllow(string(entity.RoleSuperAdmin)), usercontroller.GetById)
		routes.PUT("/:id", middleware.Authenticate(), middleware.OnlyAllow(string(entity.RoleSuperAdmin)), usercontroller.Update)
		routes.DELETE("/:id", middleware.Authenticate(), middleware.OnlyAllow(string(entity.RoleSuperAdmin)), usercontroller.Delete)
	}
}
