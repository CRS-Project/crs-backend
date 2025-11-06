package controller

import (
	"github.com/CRS-Project/crs-backend/internal/api/service"
	"github.com/CRS-Project/crs-backend/internal/dto"
	myerror "github.com/CRS-Project/crs-backend/internal/pkg/error"
	myjwt "github.com/CRS-Project/crs-backend/internal/pkg/jwt"
	"github.com/CRS-Project/crs-backend/internal/pkg/response"
	"github.com/CRS-Project/crs-backend/internal/utils"
	"github.com/gin-gonic/gin"
)

type (
	AuthController interface {
		Login(ctx *gin.Context)
		ForgetPassword(ctx *gin.Context)
		ChangePassword(ctx *gin.Context)
		Me(ctx *gin.Context)
	}

	authController struct {
		authService service.AuthService
	}
)

func NewAuth(authService service.AuthService) AuthController {
	return &authController{
		authService: authService,
	}
}

func (c *authController) Login(ctx *gin.Context) {
	var req dto.LoginRequest
	if err := ctx.ShouldBind(&req); err != nil {
		response.NewFailed("failed get data from body", myerror.ErrBodyRequest).Send(ctx)
		return
	}

	result, err := c.authService.Login(ctx.Request.Context(), req)
	if err != nil {
		response.NewFailed("failed login", err).Send(ctx)
		return
	}

	response.NewSuccess("success login", result).Send(ctx)
}

func (c *authController) ForgetPassword(ctx *gin.Context) {
	var req dto.ForgetPasswordRequest

	if err := ctx.ShouldBind(&req); err != nil {
		response.NewFailed("failed get data from body", err).Send(ctx)
		return
	}

	if err := c.authService.ForgetPassword(ctx, req); err != nil {
		response.NewFailed("failed forget password", err).Send(ctx)
		return
	}

	response.NewSuccess("success forget password", nil).Send(ctx)
}

func (c *authController) ChangePassword(ctx *gin.Context) {
	var req dto.ChangePasswordRequest
	if err := ctx.ShouldBind(&req); err != nil {
		response.NewFailed("failed get data from body", err).Send(ctx)
		return
	}

	token := ctx.Query("token")
	if token == "" {
		response.NewFailed("failed change password", myerror.ErrBodyRequest).Send(ctx)
		return
	}

	claims, err := myjwt.GetPayloadInsideToken(token)
	if err != nil {
		response.NewFailed("failed change password", err).Send(ctx)
		return
	}

	req.Email = claims["email"]
	if err := c.authService.ChangePassword(ctx, req); err != nil {
		response.NewFailed("failed change password", err).Send(ctx)
		return
	}

	response.NewSuccess("success change password", nil).Send(ctx)
}

func (c *authController) Me(ctx *gin.Context) {
	userId, err := utils.GetUserIdFromCtx(ctx)
	if err != nil {
		response.NewFailed("failed get user id", err).Send(ctx)
		return
	}

	res, err := c.authService.GetMe(ctx.Request.Context(), userId)
	if err != nil {
		response.NewFailed("failed get me", err).Send(ctx)
		return
	}

	response.NewSuccess("success get me", res).Send(ctx)
}
