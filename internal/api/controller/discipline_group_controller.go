package controller

import (
	"fmt"
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
	DisciplineGroupController interface {
		Create(ctx *gin.Context)
		GetAll(ctx *gin.Context)
		GetAllConsolidator(ctx *gin.Context)
		GetById(ctx *gin.Context)
		Update(ctx *gin.Context)
		Delete(ctx *gin.Context)
		GeneratePDF(ctx *gin.Context)
		GenerateExcel(ctx *gin.Context)
		Statistic(ctx *gin.Context)
	}

	disciplineGroupController struct {
		disciplineGroupService service.DisciplineGroupService
	}
)

func NewDisciplineGroup(disciplineGroupService service.DisciplineGroupService) DisciplineGroupController {
	return &disciplineGroupController{
		disciplineGroupService: disciplineGroupService,
	}
}

func (c *disciplineGroupController) Create(ctx *gin.Context) {
	userId, err := utils.GetUserIdFromCtx(ctx)
	if err != nil {
		response.NewFailed("failed get data from body", myerror.New(err.Error(), http.StatusBadRequest)).Send(ctx)
		return
	}

	var req dto.DisciplineGroupRequest
	if err := ctx.ShouldBind(&req); err != nil {
		err = myerror.GetErrBodyRequest(err, dto.DisciplineGroupRequest{})
		response.NewFailed("failed get data from body", err).Send(ctx)
		return
	}

	req.UserId = userId
	disciplineGroup, err := c.disciplineGroupService.Create(ctx, req)
	if err != nil {
		response.NewFailed("failed create discipline list document", err).Send(ctx)
		return
	}

	response.NewSuccess("success create discipline list document", disciplineGroup).Send(ctx)
}

func (c *disciplineGroupController) GetAll(ctx *gin.Context) {
	userId, err := utils.GetUserIdFromCtx(ctx)
	if err != nil {
		response.NewFailed("failed get data from body", myerror.New(err.Error(), http.StatusBadRequest)).Send(ctx)
		return
	}

	disciplineGroups, metaRes, err := c.disciplineGroupService.GetAll(ctx.Request.Context(), userId, meta.New(ctx))
	if err != nil {
		response.NewFailed("failed get all discipline groups", err).Send(ctx)
		return
	}

	response.NewSuccess("success get all discipline groups", disciplineGroups, metaRes).Send(ctx)
}

func (c *disciplineGroupController) GetAllConsolidator(ctx *gin.Context) {
	disciplineGroupId := ctx.Param("discipline_group_id")
	disciplineGroups, err := c.disciplineGroupService.GetAllConsolidator(ctx.Request.Context(), ctx.Query("search"), disciplineGroupId)
	if err != nil {
		response.NewFailed("failed get all discipline groups", err).Send(ctx)
		return
	}

	response.NewSuccess("success get all discipline groups", disciplineGroups).Send(ctx)
}

func (c *disciplineGroupController) GetById(ctx *gin.Context) {
	disciplineGroupId := ctx.Param("discipline_group_id")
	result, err := c.disciplineGroupService.GetById(ctx.Request.Context(), disciplineGroupId)
	if err != nil {
		response.NewFailed("failed get detail discipline list document", err).Send(ctx)
		return
	}

	response.NewSuccess("success get detail discipline list document", result).Send(ctx)
}

func (c *disciplineGroupController) Update(ctx *gin.Context) {
	disciplineGroupId := ctx.Param("discipline_group_id")
	userId, err := utils.GetUserIdFromCtx(ctx)
	if err != nil {
		response.NewFailed("failed get data from body", myerror.New(err.Error(), http.StatusBadRequest)).Send(ctx)
		return
	}

	var req dto.DisciplineGroupRequest
	if err := ctx.ShouldBind(&req); err != nil {
		err = myerror.GetErrBodyRequest(err, dto.DisciplineGroupRequest{})
		response.NewFailed("failed get data from body", myerror.New(err.Error(), http.StatusBadRequest)).Send(ctx)
		return
	}

	req.ID = disciplineGroupId
	req.UserId = userId
	err = c.disciplineGroupService.Update(ctx.Request.Context(), req)
	if err != nil {
		response.NewFailed("failed update discipline list document", err).Send(ctx)
		return
	}

	response.NewSuccess("success update discipline list document", nil).Send(ctx)
}

func (c *disciplineGroupController) Delete(ctx *gin.Context) {
	disciplineGroupId := ctx.Param("discipline_group_id")
	userId, err := utils.GetUserIdFromCtx(ctx)
	if err != nil {
		response.NewFailed("failed get data from body", myerror.New(err.Error(), http.StatusBadRequest)).Send(ctx)
		return
	}

	err = c.disciplineGroupService.Delete(ctx.Request.Context(), userId, disciplineGroupId)
	if err != nil {
		response.NewFailed("failed delete discipline list document", err).Send(ctx)
		return
	}

	response.NewSuccess("success delete discipline list document", nil).Send(ctx)
}

func (c *disciplineGroupController) GeneratePDF(ctx *gin.Context) {
	disciplineGroupId := ctx.Param("discipline_group_id")
	userId, err := utils.GetUserIdFromCtx(ctx)
	if err != nil {
		response.NewFailed("failed get data from body", myerror.New(err.Error(), http.StatusBadRequest)).Send(ctx)
		return
	}

	pdfBuffer, filename, err := c.disciplineGroupService.GeneratePDF(ctx.Request.Context(), userId, disciplineGroupId)
	if err != nil {
		response.NewFailed("failed generate pdf", err).Send(ctx)
		return
	}

	ctx.Header("Content-Type", "application/pdf")
	ctx.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))
	ctx.Data(http.StatusOK, "application/pdf", pdfBuffer.Bytes())
	response.NewSuccess("success generate pdf", nil).Send(ctx)
}

func (c *disciplineGroupController) GenerateExcel(ctx *gin.Context) {
	disciplineGroupId := ctx.Param("discipline_group_id")
	userId, err := utils.GetUserIdFromCtx(ctx)
	if err != nil {
		response.NewFailed("failed get data from body", myerror.New(err.Error(), http.StatusBadRequest)).Send(ctx)
		return
	}

	excelBuffer, filename, err := c.disciplineGroupService.GenerateExcel(ctx.Request.Context(), userId, disciplineGroupId)
	if err != nil {
		response.NewFailed("failed generate excel", err).Send(ctx)
		return
	}

	ctx.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	ctx.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))
	ctx.Data(http.StatusOK, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", excelBuffer.Bytes())
	response.NewSuccess("success generate excel", nil).Send(ctx)
}

func (c *disciplineGroupController) Statistic(ctx *gin.Context) {
	packageId := ctx.Param("package_id")
	res, err := c.disciplineGroupService.GetStatistic(ctx.Request.Context(), packageId)
	if err != nil {
		response.NewFailed("failed get statistic discipline group", err).Send(ctx)
		return
	}

	response.NewSuccess("success get statistic discipline group", res).Send(ctx)
}
