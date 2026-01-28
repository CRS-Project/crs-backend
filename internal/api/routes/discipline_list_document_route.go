package routes

import (
	"github.com/CRS-Project/crs-backend/internal/api/controller"
	"github.com/CRS-Project/crs-backend/internal/middleware"
	"github.com/gin-gonic/gin"
)

func DisciplineListDocument(app *gin.Engine, areaOfConcerncontroller controller.DisciplineListDocumentController, middleware middleware.Middleware) {
	routes := app.Group("/api/v1/discipline-group/:discipline_group_id/discipline-list-document")
	{
		routes.POST("", middleware.Authenticate(), areaOfConcerncontroller.Create)
		routes.GET("", middleware.Authenticate(), areaOfConcerncontroller.GetAll)
		routes.GET("/:discipline_list_document_id", middleware.Authenticate(), areaOfConcerncontroller.GetById)
		routes.PUT("/:discipline_list_document_id", middleware.Authenticate(), areaOfConcerncontroller.Update)
		routes.DELETE("/:discipline_list_document_id", middleware.Authenticate(), areaOfConcerncontroller.Delete)
		routes.GET("/:discipline_list_document_id/generate-excel", middleware.Authenticate(), areaOfConcerncontroller.GenerateExcel)
	}
}
