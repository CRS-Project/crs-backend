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
	AreaOfConcernGroupController interface {
		Create(ctx *gin.Context)
		GetAll(ctx *gin.Context)
		GetById(ctx *gin.Context)
		Update(ctx *gin.Context)
		Delete(ctx *gin.Context)
	}

	areaOfConcernGroupController struct {
		areaOfConcernGroupService service.AreaOfConcernGroupService
	}
)

func NewAreaOfConcernGroup(areaOfConcernGroupService service.AreaOfConcernGroupService) AreaOfConcernGroupController {
	return &areaOfConcernGroupController{
		areaOfConcernGroupService: areaOfConcernGroupService,
	}
}

func (c *areaOfConcernGroupController) Create(ctx *gin.Context) {
	userId, err := utils.GetUserIdFromCtx(ctx)
	if err != nil {
		response.NewFailed("failed get data from body", myerror.New(err.Error(), http.StatusBadRequest)).Send(ctx)
		return
	}

	var req dto.AreaOfConcernGroupRequest
	if err := ctx.ShouldBind(&req); err != nil {
		err = myerror.GetErrBodyRequest(err, dto.AreaOfConcernGroupRequest{})
		response.NewFailed("failed get data from body", err).Send(ctx)
		return
	}

	req.UserId = userId
	areaOfConcernGroup, err := c.areaOfConcernGroupService.Create(ctx, req)
	if err != nil {
		response.NewFailed("failed create area of concern", err).Send(ctx)
		return
	}

	response.NewSuccess("success create area of concern", areaOfConcernGroup).Send(ctx)
}

func (c *areaOfConcernGroupController) GetAll(ctx *gin.Context) {
	userId, err := utils.GetUserIdFromCtx(ctx)
	if err != nil {
		response.NewFailed("failed get data from body", myerror.New(err.Error(), http.StatusBadRequest)).Send(ctx)
		return
	}

	areaOfConcernGroups, metaRes, err := c.areaOfConcernGroupService.GetAll(ctx.Request.Context(), userId, meta.New(ctx))
	if err != nil {
		response.NewFailed("failed get all areaOfConcernGroups", err).Send(ctx)
		return
	}

	response.NewSuccess("success get all areaOfConcernGroups", areaOfConcernGroups, metaRes).Send(ctx)
}

func (c *areaOfConcernGroupController) GetById(ctx *gin.Context) {
	areaOfConcernGroupId := ctx.Param("area_of_concern_group_id")
	result, err := c.areaOfConcernGroupService.GetById(ctx.Request.Context(), areaOfConcernGroupId)
	if err != nil {
		response.NewFailed("failed get detail area of concern", err).Send(ctx)
		return
	}

	response.NewSuccess("success get detail area of concern", result).Send(ctx)
}

func (c *areaOfConcernGroupController) Update(ctx *gin.Context) {
	areaOfConcernGroupId := ctx.Param("area_of_concern_group_id")
	userId, err := utils.GetUserIdFromCtx(ctx)
	if err != nil {
		response.NewFailed("failed get data from body", myerror.New(err.Error(), http.StatusBadRequest)).Send(ctx)
		return
	}

	var req dto.AreaOfConcernGroupRequest
	if err := ctx.ShouldBind(&req); err != nil {
		err = myerror.GetErrBodyRequest(err, dto.AreaOfConcernGroupRequest{})
		response.NewFailed("failed get data from body", myerror.New(err.Error(), http.StatusBadRequest)).Send(ctx)
		return
	}

	req.ID = areaOfConcernGroupId
	req.UserId = userId
	err = c.areaOfConcernGroupService.Update(ctx.Request.Context(), req)
	if err != nil {
		response.NewFailed("failed update area of concern", err).Send(ctx)
		return
	}

	response.NewSuccess("success update area of concern", nil).Send(ctx)
}

func (c *areaOfConcernGroupController) Delete(ctx *gin.Context) {
	areaOfConcernGroupId := ctx.Param("area_of_concern_id")
	userId, err := utils.GetUserIdFromCtx(ctx)
	if err != nil {
		response.NewFailed("failed get data from body", myerror.New(err.Error(), http.StatusBadRequest)).Send(ctx)
		return
	}

	err = c.areaOfConcernGroupService.Delete(ctx.Request.Context(), userId, areaOfConcernGroupId)
	if err != nil {
		response.NewFailed("failed delete area of concern", err).Send(ctx)
		return
	}

	response.NewSuccess("success delete area of concern", nil).Send(ctx)
}
