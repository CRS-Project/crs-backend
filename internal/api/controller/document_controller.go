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
	DocumentController interface {
		Create(ctx *gin.Context)
		CreateBulk(ctx *gin.Context)
		GetByID(ctx *gin.Context)
		GetAll(ctx *gin.Context)
		Update(ctx *gin.Context)
		Delete(ctx *gin.Context)
	}

	documentController struct {
		documentService service.DocumentService
	}
)

func NewDocument(documentService service.DocumentService) DocumentController {
	return &documentController{
		documentService: documentService,
	}
}

func (c *documentController) Create(ctx *gin.Context) {
	userId, err := utils.GetUserIdFromCtx(ctx)
	if err != nil {
		response.NewFailed("failed get data from body", myerror.New(err.Error(), http.StatusBadRequest)).Send(ctx)
		return
	}

	var req dto.CreateDocumentRequest
	if err := ctx.ShouldBind(&req); err != nil {
		response.NewFailed("failed get data from body", err).Send(ctx)
		return
	}

	req.UserID = userId

	res, err := c.documentService.Create(ctx.Request.Context(), req)
	if err != nil {
		response.NewFailed("failed to create document", err).Send(ctx)
		return
	}

	response.NewSuccess("success create document", res).Send(ctx)
}

func (c *documentController) CreateBulk(ctx *gin.Context) {
	userId, err := utils.GetUserIdFromCtx(ctx)
	if err != nil {
		response.NewFailed("failed get data from body", myerror.New(err.Error(), http.StatusBadRequest)).Send(ctx)
		return
	}

	var req dto.CreateBulkDocumentRequest
	if err := ctx.ShouldBind(&req); err != nil {
		response.NewFailed("failed get data from body", err).Send(ctx)
		return
	}

	packageId := ctx.Param("package_id")

	req.UserID = userId
	req.PackageID = packageId

	res, err := c.documentService.CreateBulk(ctx.Request.Context(), req)
	if err != nil {
		response.NewFailed("failed to create bulk document", err).Send(ctx)
		return
	}

	response.NewSuccess("success to create bulk document", res).Send(ctx)
}

func (c *documentController) GetByID(ctx *gin.Context) {
	documentId := ctx.Param("document_id")

	res, err := c.documentService.GetByID(ctx.Request.Context(), documentId)
	if err != nil {
		response.NewFailed("failed to get document", err).Send(ctx)
		return
	}

	response.NewSuccess("success get document", res).Send(ctx)
}

func (c *documentController) GetAll(ctx *gin.Context) {
	userId, err := utils.GetUserIdFromCtx(ctx)
	if err != nil {
		response.NewFailed("failed get data from body", myerror.New(err.Error(), http.StatusBadRequest)).Send(ctx)
		return
	}

	res, metaRes, err := c.documentService.GetAll(ctx.Request.Context(), userId, meta.New(ctx))
	if err != nil {
		response.NewFailed("failed to get documents", err).Send(ctx)
		return
	}

	response.NewSuccess("success get documents", res, metaRes).Send(ctx)
}

func (c *documentController) Update(ctx *gin.Context) {
	documentId := ctx.Param("document_id")
	userId, err := utils.GetUserIdFromCtx(ctx)
	if err != nil {
		response.NewFailed("failed get data from body", myerror.New(err.Error(), http.StatusBadRequest)).Send(ctx)
		return
	}

	var req dto.UpdateDocumentRequest
	if err := ctx.ShouldBind(&req); err != nil {
		response.NewFailed("failed get data from body", err).Send(ctx)
		return
	}

	req.UserID = userId
	req.ID = documentId

	res, err := c.documentService.Update(ctx.Request.Context(), req)
	if err != nil {
		response.NewFailed("failed to update document", err).Send(ctx)
		return
	}

	response.NewSuccess("success update document", res).Send(ctx)
}

func (c *documentController) Delete(ctx *gin.Context) {
	documentId := ctx.Param("document_id")
	userId, err := utils.GetUserIdFromCtx(ctx)
	if err != nil {
		response.NewFailed("failed get data from body", myerror.New(err.Error(), http.StatusBadRequest)).Send(ctx)
		return
	}

	err = c.documentService.Delete(ctx.Request.Context(), userId, documentId)
	if err != nil {
		response.NewFailed("failed to delete document", err).Send(ctx)
		return
	}

	response.NewSuccess("success delete document", nil).Send(ctx)
}
