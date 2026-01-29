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
	DisciplineListDocumentController interface {
		Create(ctx *gin.Context)
		GetAll(ctx *gin.Context)
		GetById(ctx *gin.Context)
		Update(ctx *gin.Context)
		Delete(ctx *gin.Context)
		GenerateExcel(ctx *gin.Context)
	}

	disciplineListDocumentController struct {
		disciplineListDocumentService service.DisciplineListDocumentService
	}
)

func NewDisciplineListDocument(disciplineListDocumentService service.DisciplineListDocumentService) DisciplineListDocumentController {
	return &disciplineListDocumentController{
		disciplineListDocumentService: disciplineListDocumentService,
	}
}

func (c *disciplineListDocumentController) Create(ctx *gin.Context) {
	disciplineGroupId := ctx.Param("discipline_group_id")

	userId, err := utils.GetUserIdFromCtx(ctx)
	if err != nil {
		response.NewFailed("failed get data from body", myerror.New(err.Error(), http.StatusBadRequest)).Send(ctx)
		return
	}

	var req dto.DisciplineListDocumentRequest
	if err := ctx.ShouldBind(&req); err != nil {
		err = myerror.GetErrBodyRequest(err, dto.DisciplineListDocumentRequest{})
		response.NewFailed("failed get data from body", err).Send(ctx)
		return
	}

	req.UserId = userId
	req.DisciplineGroupID = disciplineGroupId
	disciplineListDocument, err := c.disciplineListDocumentService.Create(ctx, req)
	if err != nil {
		response.NewFailed("failed create discipline list document", err).Send(ctx)
		return
	}

	response.NewSuccess("success create discipline list document", disciplineListDocument).Send(ctx)
}

func (c *disciplineListDocumentController) GetAll(ctx *gin.Context) {
	disciplineGroupId := ctx.Param("discipline_group_id")
	userId, err := utils.GetUserIdFromCtx(ctx)
	if err != nil {
		response.NewFailed("failed get data from body", myerror.New(err.Error(), http.StatusBadRequest)).Send(ctx)
		return
	}

	disciplineListDocuments, metaRes, err := c.disciplineListDocumentService.GetAll(ctx.Request.Context(), disciplineGroupId, userId, meta.NewWithDefault(ctx, 0, 0, "desc", "due_date"))
	if err != nil {
		response.NewFailed("failed get all area of Concerns", err).Send(ctx)
		return
	}

	response.NewSuccess("success get all discipline list documents", disciplineListDocuments, metaRes).Send(ctx)
}

func (c *disciplineListDocumentController) GetById(ctx *gin.Context) {
	disciplineListDocumentId := ctx.Param("discipline_list_document_id")
	result, err := c.disciplineListDocumentService.GetById(ctx.Request.Context(), disciplineListDocumentId)
	if err != nil {
		response.NewFailed("failed get detail discipline list document", err).Send(ctx)
		return
	}

	response.NewSuccess("success get detail discipline list document", result).Send(ctx)
}

func (c *disciplineListDocumentController) Update(ctx *gin.Context) {
	disciplineListDocumentId := ctx.Param("discipline_list_document_id")
	userId, err := utils.GetUserIdFromCtx(ctx)
	if err != nil {
		response.NewFailed("failed get data from body", myerror.New(err.Error(), http.StatusBadRequest)).Send(ctx)
		return
	}

	var req dto.UpdateDisciplineListDocumentRequest
	if err := ctx.ShouldBind(&req); err != nil {
		err = myerror.GetErrBodyRequest(err, dto.DisciplineListDocumentRequest{})
		response.NewFailed("failed get data from body", myerror.New(err.Error(), http.StatusBadRequest)).Send(ctx)
		return
	}

	req.ID = disciplineListDocumentId
	req.UserId = userId
	err = c.disciplineListDocumentService.Update(ctx.Request.Context(), req)
	if err != nil {
		response.NewFailed("failed update discipline list document", err).Send(ctx)
		return
	}

	response.NewSuccess("success update discipline list document", nil).Send(ctx)
}

func (c *disciplineListDocumentController) Delete(ctx *gin.Context) {
	disciplineListDocumentId := ctx.Param("discipline_list_document_id")
	userId, err := utils.GetUserIdFromCtx(ctx)
	if err != nil {
		response.NewFailed("failed get data from body", myerror.New(err.Error(), http.StatusBadRequest)).Send(ctx)
		return
	}

	err = c.disciplineListDocumentService.Delete(ctx.Request.Context(), userId, disciplineListDocumentId)
	if err != nil {
		response.NewFailed("failed delete discipline list document", err).Send(ctx)
		return
	}

	response.NewSuccess("success delete discipline list document", nil).Send(ctx)
}

func (c *disciplineListDocumentController) GenerateExcel(ctx *gin.Context) {
	disciplineListDocumentId := ctx.Param("discipline_list_document_id")
	userId, err := utils.GetUserIdFromCtx(ctx)
	if err != nil {
		response.NewFailed("failed get data from body", myerror.New(err.Error(), http.StatusBadRequest)).Send(ctx)
		return
	}

	excelBuffer, filename, err := c.disciplineListDocumentService.GenerateExcel(ctx.Request.Context(), userId, disciplineListDocumentId)
	if err != nil {
		response.NewFailed("failed generate excel", err).Send(ctx)
		return
	}

	ctx.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	ctx.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))
	ctx.Data(http.StatusOK, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", excelBuffer.Bytes())
}
