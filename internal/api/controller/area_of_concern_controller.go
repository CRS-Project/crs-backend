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
	AreaOfConcernController interface {
		Create(ctx *gin.Context)
		GetAll(ctx *gin.Context)
		GetById(ctx *gin.Context)
		Update(ctx *gin.Context)
		Delete(ctx *gin.Context)
	}

	areaOfConcernController struct {
		areaOfConcernService service.AreaOfConcernService
	}
)

func NewAreaOfConcern(areaOfConcernService service.AreaOfConcernService) AreaOfConcernController {
	return &areaOfConcernController{
		areaOfConcernService: areaOfConcernService,
	}
}

func (c *areaOfConcernController) Create(ctx *gin.Context) {
	areaOfConcernGroupId := ctx.Param("area_of_concern_group_id")

	userId, err := utils.GetUserIdFromCtx(ctx)
	if err != nil {
		response.NewFailed("failed get data from body", myerror.New(err.Error(), http.StatusBadRequest)).Send(ctx)
		return
	}

	var req dto.AreaOfConcernRequest
	if err := ctx.ShouldBind(&req); err != nil {
		err = myerror.GetErrBodyRequest(err, dto.AreaOfConcernRequest{})
		response.NewFailed("failed get data from body", err).Send(ctx)
		return
	}

	req.UserId = userId
	req.AreaOfConcernGroupID = areaOfConcernGroupId
	areaOfConcern, err := c.areaOfConcernService.Create(ctx, req)
	if err != nil {
		response.NewFailed("failed create area of concern", err).Send(ctx)
		return
	}

	response.NewSuccess("success create area of concern", areaOfConcern).Send(ctx)
}

func (c *areaOfConcernController) GetAll(ctx *gin.Context) {
	areaOfConcernGroupId := ctx.Param("area_of_concern_group_id")
	userId, err := utils.GetUserIdFromCtx(ctx)
	if err != nil {
		response.NewFailed("failed get data from body", myerror.New(err.Error(), http.StatusBadRequest)).Send(ctx)
		return
	}

	areaOfConcerns, metaRes, err := c.areaOfConcernService.GetAll(ctx.Request.Context(), areaOfConcernGroupId, userId, meta.New(ctx))
	if err != nil {
		response.NewFailed("failed get all area of Concerns", err).Send(ctx)
		return
	}

	response.NewSuccess("success get all area of concerns", areaOfConcerns, metaRes).Send(ctx)
}

func (c *areaOfConcernController) GetById(ctx *gin.Context) {
	areaOfConcernId := ctx.Param("area_of_concern_id")
	result, err := c.areaOfConcernService.GetById(ctx.Request.Context(), areaOfConcernId)
	if err != nil {
		response.NewFailed("failed get detail area of concern", err).Send(ctx)
		return
	}

	response.NewSuccess("success get detail area of concern", result).Send(ctx)
}

func (c *areaOfConcernController) Update(ctx *gin.Context) {
	areaOfConcernId := ctx.Param("area_of_concern_id")
	userId, err := utils.GetUserIdFromCtx(ctx)
	if err != nil {
		response.NewFailed("failed get data from body", myerror.New(err.Error(), http.StatusBadRequest)).Send(ctx)
		return
	}

	var req dto.UpdateAreaOfConcernRequest
	if err := ctx.ShouldBind(&req); err != nil {
		err = myerror.GetErrBodyRequest(err, dto.AreaOfConcernRequest{})
		response.NewFailed("failed get data from body", myerror.New(err.Error(), http.StatusBadRequest)).Send(ctx)
		return
	}

	req.ID = areaOfConcernId
	req.UserId = userId
	err = c.areaOfConcernService.Update(ctx.Request.Context(), req)
	if err != nil {
		response.NewFailed("failed update area of concern", err).Send(ctx)
		return
	}

	response.NewSuccess("success update area of concern", nil).Send(ctx)
}

func (c *areaOfConcernController) Delete(ctx *gin.Context) {
	areaOfConcernId := ctx.Param("area_of_concern_id")
	userId, err := utils.GetUserIdFromCtx(ctx)
	if err != nil {
		response.NewFailed("failed get data from body", myerror.New(err.Error(), http.StatusBadRequest)).Send(ctx)
		return
	}

	err = c.areaOfConcernService.Delete(ctx.Request.Context(), userId, areaOfConcernId)
	if err != nil {
		response.NewFailed("failed delete area of concern", err).Send(ctx)
		return
	}

	response.NewSuccess("success delete area of concern", nil).Send(ctx)
}
