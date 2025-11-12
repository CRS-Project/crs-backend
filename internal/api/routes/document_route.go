package routes

import (
	"github.com/CRS-Project/crs-backend/internal/api/controller"
	"github.com/CRS-Project/crs-backend/internal/entity"
	"github.com/CRS-Project/crs-backend/internal/middleware"
	"github.com/gin-gonic/gin"
)

func Document(app *gin.Engine, documentcontroller controller.DocumentController, middleware middleware.Middleware) {
	routes := app.Group("/api/v1/document")
	{
		routes.POST("", middleware.Authenticate(), middleware.OnlyAllow(string(entity.RoleContractor), string(entity.RoleSuperAdmin)), documentcontroller.Create)
		routes.POST("/bulk/:package_id", middleware.Authenticate(), middleware.OnlyAllow(string(entity.RoleContractor), string(entity.RoleSuperAdmin)), documentcontroller.CreateBulk)
		routes.GET("", middleware.Authenticate(), documentcontroller.GetAll)
		routes.GET("/:document_id", middleware.Authenticate(), documentcontroller.GetByID)
		routes.PUT("/:document_id", middleware.Authenticate(), documentcontroller.Update)
		routes.DELETE("/:document_id", middleware.Authenticate(), documentcontroller.Delete)
	}
}
