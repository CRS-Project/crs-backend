package controller

import (
	"net/http"

	"github.com/CRS-Project/crs-backend/internal/api/service"
	"github.com/CRS-Project/crs-backend/internal/dto"
	myerror "github.com/CRS-Project/crs-backend/internal/pkg/error"
	"github.com/CRS-Project/crs-backend/internal/pkg/response"
	"github.com/gin-gonic/gin"
)

type (
	UserController interface {
		Create(ctx *gin.Context)
		GetById(ctx *gin.Context)
	}

	userController struct {
		userService service.UserService
	}
)

func NewUser(userService service.UserService) UserController {
	return &userController{
		userService: userService,
	}
}

func (c *userController) Create(ctx *gin.Context) {
	var req dto.CreateUserRequest

	if err := ctx.ShouldBind(&req); err != nil {
		response.NewFailed("failed get data from body", myerror.New(err.Error(), http.StatusBadRequest)).Send(ctx)
		return
	}

	user, err := c.userService.Create(ctx, req)
	if err != nil {
		response.NewFailed("failed create account", err).Send(ctx)
		return
	}

	response.NewSuccess("success create account", user).Send(ctx)
}

func (c *userController) GetById(ctx *gin.Context) {
	userId := ctx.Param("id")
	result, err := c.userService.GetById(ctx.Request.Context(), userId)
	if err != nil {
		response.NewFailed("failed get detail user", err).Send(ctx)
		return
	}

	response.NewSuccess("success get detail user", result).Send(ctx)
}
