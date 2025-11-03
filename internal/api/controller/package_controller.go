package controller

import (
	"github.com/CRS-Project/crs-backend/internal/api/service"
	"github.com/CRS-Project/crs-backend/internal/dto"
	"github.com/CRS-Project/crs-backend/internal/pkg/response"
	"github.com/gin-gonic/gin"
)

type (
	PackageController interface {
		CreatePackage(ctx *gin.Context)
		GetPackages(ctx *gin.Context)
		UpdatePackage(ctx *gin.Context)
		DeletePackage(ctx *gin.Context)
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
		response.NewFailed("failed get data from body", err).Send(ctx)
		return
	}

	res, err := c.packageService.CreatePackage(ctx, req)
	if err != nil {
		response.NewFailed("failed to create package", err).Send(ctx)
		return
	}

	response.NewSuccess("success create package", res).Send(ctx)
}

func (c *packageController) GetPackages(ctx *gin.Context) {
	res, err := c.packageService.GetPackages(ctx)
	if err != nil {
		response.NewFailed("failed to get packages", err).Send(ctx)
		return
	}

	response.NewSuccess("success get package", res).Send(ctx)
}

func (c *packageController) UpdatePackage(ctx *gin.Context) {
	var req dto.UpdatePackageRequest
	if err := ctx.ShouldBind(&req); err != nil {
		response.NewFailed("failed get data from body", err).Send(ctx)
		return
	}

	if err := c.packageService.UpdatePackage(ctx, req); err != nil {
		response.NewFailed("failed to update package", err).Send(ctx)
		return
	}

	response.NewSuccess("success update package", nil).Send(ctx)
}

func (c *packageController) DeletePackage(ctx *gin.Context) {
	var req dto.DeletePackageRequest
	if err := ctx.ShouldBind(&req); err != nil {
		response.NewFailed("failed get data from body", err).Send(ctx)
		return
	}

	if err := c.packageService.DeletePackage(ctx, req); err != nil {
		response.NewFailed("failed to delete package", err).Send(ctx)
		return
	}

	response.NewSuccess("success delete package", nil).Send(ctx)
}
