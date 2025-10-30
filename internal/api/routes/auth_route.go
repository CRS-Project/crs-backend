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
		routes.POST("/register", authcontroller.Register)
		routes.GET("/verify", authcontroller.Verify)
		routes.GET("/me", middleware.Authenticate(), authcontroller.Me)

		routes.GET("/google/login", authcontroller.LoginWithGoogle)
		routes.GET("/google/callback", authcontroller.CallbackGoogle)
	}
}
