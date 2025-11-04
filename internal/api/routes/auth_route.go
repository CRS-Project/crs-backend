package routes

import (
	"github.com/CRS-Project/crs-backend/internal/api/controller"
	"github.com/CRS-Project/crs-backend/internal/middleware"
	"github.com/gin-gonic/gin"
)

func Auth(app *gin.Engine, authcontroller controller.AuthController, middleware middleware.Middleware) {
	routes := app.Group("/api/v1/auth")
	{
		routes.POST("/login", authcontroller.Login)
		routes.POST("/forget", authcontroller.ForgetPassword)
		routes.POST("/change", authcontroller.ChangePassword)
		routes.GET("/me", middleware.Authenticate(), authcontroller.Me)
	}
}
