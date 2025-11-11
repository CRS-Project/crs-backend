package controller

import (
	"net/http"

	"github.com/CRS-Project/crs-backend/internal/api/service"
	"github.com/CRS-Project/crs-backend/internal/dto"
	myerror "github.com/CRS-Project/crs-backend/internal/pkg/error"
	"github.com/CRS-Project/crs-backend/internal/pkg/meta"
	"github.com/CRS-Project/crs-backend/internal/pkg/response"
	"github.com/CRS-Project/crs-backend/internal/utils"
	"github.com/gin-gonic/gin"
)

type (
	PackageController interface {
		CreatePackage(ctx *gin.Context)
		GetAll(ctx *gin.Context)
		GetAllByUser(ctx *gin.Context)
		UpdatePackage(ctx *gin.Context)
		DeletePackage(ctx *gin.Context)
		GetByID(ctx *gin.Context)
	}

	packageController struct {
		packageService service.PackageService
	}
)

func NewPackage(packageService service.PackageService) PackageController {
	return &packageController{
		packageService: packageService,
	}
}

func (c *packageController) CreatePackage(ctx *gin.Context) {
	var req dto.CreatePackageRequest
	if err := ctx.ShouldBind(&req); err != nil {
		err = myerror.GetErrBodyRequest(err, dto.CreatePackageRequest{})
		response.NewFailed("failed get data from body", err).Send(ctx)
		return
	}

	res, err := c.packageService.CreatePackage(ctx.Request.Context(), req)
	if err != nil {
		response.NewFailed("failed to create package", err).Send(ctx)
		return
	}

	response.NewSuccess("success create package", res).Send(ctx)
}

func (c *packageController) GetAll(ctx *gin.Context) {
	res, metaRes, err := c.packageService.GetAll(ctx.Request.Context(), meta.New(ctx))
	if err != nil {
		response.NewFailed("failed to get packages", err).Send(ctx)
		return
	}

	response.NewSuccess("success get package", res, metaRes).Send(ctx)
}

func (c *packageController) GetAllByUser(ctx *gin.Context) {
	userId, err := utils.GetUserIdFromCtx(ctx)
	if err != nil {
		response.NewFailed("failed get data from body", myerror.New(err.Error(), http.StatusBadRequest)).Send(ctx)
		return
	}

	res, err := c.packageService.GetAllByUser(ctx.Request.Context(), userId)
	if err != nil {
		response.NewFailed("failed to get packages", err).Send(ctx)
		return
	}

	response.NewSuccess("success get package", res).Send(ctx)
}

func (c *packageController) UpdatePackage(ctx *gin.Context) {
	var req dto.UpdatePackageRequest
	if err := ctx.ShouldBind(&req); err != nil {
		err = myerror.GetErrBodyRequest(err, dto.UpdatePackageRequest{})
		response.NewFailed("failed get data from body", err).Send(ctx)
		return
	}

	res, err := c.packageService.UpdatePackage(ctx.Request.Context(), req)
	if err != nil {
		response.NewFailed("failed to update package", err).Send(ctx)
		return
	}

	response.NewSuccess("success update package", res).Send(ctx)
}

func (c *packageController) DeletePackage(ctx *gin.Context) {
	id := ctx.Param("id")

	if err := c.packageService.DeletePackage(ctx.Request.Context(), id); err != nil {
		response.NewFailed("failed to delete package", err).Send(ctx)
		return
	}

	response.NewSuccess("success delete package", nil).Send(ctx)
}

func (c *packageController) GetByID(ctx *gin.Context) {
	id := ctx.Param("id")

	res, err := c.packageService.GetByID(ctx.Request.Context(), id)
	if err != nil {
		response.NewFailed("failed to get package", err).Send(ctx)
		return
	}

	response.NewSuccess("success get package", res).Send(ctx)
}
