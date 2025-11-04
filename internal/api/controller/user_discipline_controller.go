package controller

import (
	"github.com/CRS-Project/crs-backend/internal/api/service"
	"github.com/CRS-Project/crs-backend/internal/pkg/meta"
	"github.com/CRS-Project/crs-backend/internal/pkg/response"
	"github.com/gin-gonic/gin"
)

type (
	UserDisciplineController interface {
		GetAll(ctx *gin.Context)
	}

	userDisciplineController struct {
		userDisciplineService service.UserDisciplineService
	}
)

func NewUserDiscipline(userDisciplineService service.UserDisciplineService) UserDisciplineController {
	return &userDisciplineController{userDisciplineService}
}

func (c *userDisciplineController) GetAll(ctx *gin.Context) {
	userDiscipline, err := c.userDisciplineService.GetAll(ctx, meta.New(ctx))
	if err != nil {
		response.NewFailed("failed get all user discipline", err).Send(ctx)
		return
	}

	response.NewSuccess("success get all user discipline", userDiscipline).Send(ctx)
}
