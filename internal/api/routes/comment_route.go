package routes

import (
	"github.com/CRS-Project/crs-backend/internal/api/controller"
	"github.com/CRS-Project/crs-backend/internal/entity"
	"github.com/CRS-Project/crs-backend/internal/middleware"
	"github.com/gin-gonic/gin"
)

func Comment(app *gin.Engine, commentcontroller controller.CommentController, middleware middleware.Middleware) {
	routes := app.Group("/api/v1/area-of-concern-group/:area_of_concern_group_id/area-of-concern/:area_of_concern_id/comment")
	{
		routes.POST("", middleware.Authenticate(), middleware.OnlyAllow(string(entity.RoleSuperAdmin), string(entity.RoleReviewer)), commentcontroller.Create)
		routes.POST("/:comment_id/reply", middleware.Authenticate(), commentcontroller.ReplyId)

		routes.GET("", middleware.Authenticate(), commentcontroller.GetAllByAreaOfConcernId)
		routes.GET("/:comment_id", middleware.Authenticate(), commentcontroller.GetById)
		routes.GET("/:comment_id/reply", middleware.Authenticate(), commentcontroller.GetAllReplyByCommentId)

		routes.PUT("/:comment_id", middleware.Authenticate(), commentcontroller.Update)
		routes.DELETE("/:comment_id", middleware.Authenticate(), commentcontroller.Delete)
	}
}
