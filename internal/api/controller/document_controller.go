package controller

import (
	"github.com/CRS-Project/crs-backend/internal/api/service"
	"github.com/CRS-Project/crs-backend/internal/dto"
	"github.com/CRS-Project/crs-backend/internal/pkg/meta"
	"github.com/CRS-Project/crs-backend/internal/pkg/response"
	"github.com/gin-gonic/gin"
)

type (
	DocumentController interface {
		CreateDocument(ctx *gin.Context)
		GetAll(ctx *gin.Context)
		DeleteDocument(ctx *gin.Context)
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

func (c *documentController) CreateDocument(ctx *gin.Context) {
	var req dto.CreateDocumentRequest
	if err := ctx.ShouldBind(&req); err != nil {
		response.NewFailed("failed get data from body", err).Send(ctx)
		return
	}

	res, err := c.documentService.Create(ctx, req)
	if err != nil {
		response.NewFailed("failed to create document", err).Send(ctx)
		return
	}

	response.NewSuccess("success create document", res).Send(ctx)
}

func (c *documentController) GetAll(ctx *gin.Context) {
	res, metaRes, err := c.documentService.GetAll(ctx, meta.New(ctx))
	if err != nil {
		response.NewFailed("failed to get documents", err).Send(ctx)
		return
	}

	response.NewSuccess("success get documents", res, metaRes).Send(ctx)
}

func (c *documentController) DeleteDocument(ctx *gin.Context) {
	id := ctx.Param("id")

	err := c.documentService.Delete(ctx, id)
	if err != nil {
		response.NewFailed("failed to delete document", err).Send(ctx)
		return
	}

	response.NewSuccess("success delete document", nil).Send(ctx)
}
