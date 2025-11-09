package controller

import (
	"github.com/CRS-Project/crs-backend/internal/api/service"
	"github.com/CRS-Project/crs-backend/internal/dto"
	myerror "github.com/CRS-Project/crs-backend/internal/pkg/error"
	"github.com/CRS-Project/crs-backend/internal/pkg/meta"
	"github.com/CRS-Project/crs-backend/internal/pkg/response"
	"github.com/gin-gonic/gin"
)

type (
	UserController interface {
		Create(ctx *gin.Context)
		GetAll(ctx *gin.Context)
		GetById(ctx *gin.Context)
		Update(ctx *gin.Context)
		Delete(ctx *gin.Context)
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
		err = myerror.GetErrBodyRequest(err, dto.CreateUserRequest{})
		response.NewFailed("failed get data from body", err).Send(ctx)
		return
	}

	user, err := c.userService.Create(ctx, req)
	if err != nil {
		response.NewFailed("failed create account", err).Send(ctx)
		return
	}

	response.NewSuccess("success create account", user).Send(ctx)
}

func (c *userController) GetAll(ctx *gin.Context) {
	users, metaRes, err := c.userService.GetAll(ctx.Request.Context(), meta.New(ctx))
	if err != nil {
		response.NewFailed("failed get all users", err).Send(ctx)
		return
	}

	response.NewSuccess("success get all users", users, metaRes).Send(ctx)
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

func (c *userController) Update(ctx *gin.Context) {
	userId := ctx.Param("id")
	var req dto.UpdateUserRequest

	if err := ctx.ShouldBind(&req); err != nil {
		err = myerror.GetErrBodyRequest(err, dto.UpdateUserRequest{})
		response.NewFailed("failed get data from body", err).Send(ctx)
		return
	}

	result, err := c.userService.Update(ctx.Request.Context(), userId, req)
	if err != nil {
		response.NewFailed("failed update user", err).Send(ctx)
		return
	}

	response.NewSuccess("success update user", result).Send(ctx)
}

func (c *userController) Delete(ctx *gin.Context) {
	userId := ctx.Param("id")

	err := c.userService.Delete(ctx.Request.Context(), userId)
	if err != nil {
		response.NewFailed("failed delete user", err).Send(ctx)
		return
	}

	response.NewSuccess("success delete user", nil).Send(ctx)
}
