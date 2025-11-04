package routes

import (
	"github.com/CRS-Project/crs-backend/internal/api/controller"
	"github.com/CRS-Project/crs-backend/internal/middleware"
	"github.com/gin-gonic/gin"
)

func UserDiscipline(app *gin.Engine, userdisciplinecontroller controller.UserDisciplineController, middleware middleware.Middleware) {
	routes := app.Group("/api/v1/user-discipline")
	{
		routes.GET("", userdisciplinecontroller.GetAll)
	}
}
