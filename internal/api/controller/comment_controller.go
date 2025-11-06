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
	CommentController interface {
		Create(ctx *gin.Context)
		ReplyId(ctx *gin.Context)
		GetAllByDocumentId(ctx *gin.Context)
		GetAllReplyByCommentId(ctx *gin.Context)
		GetById(ctx *gin.Context)
		Update(ctx *gin.Context)
		Delete(ctx *gin.Context)
	}

	commentController struct {
		commentService service.CommentService
	}
)

func NewComment(commentService service.CommentService) CommentController {
	return &commentController{
		commentService: commentService,
	}
}

func (c *commentController) Create(ctx *gin.Context) {
	documentId := ctx.Param("document_id")

	userId, err := utils.GetUserIdFromCtx(ctx)
	if err != nil {
		response.NewFailed("failed get data from body", myerror.New(err.Error(), http.StatusBadRequest)).Send(ctx)
		return
	}

	var req dto.CommentRequest
	if err := ctx.ShouldBind(&req); err != nil {
		response.NewFailed("failed get data from body", myerror.New(err.Error(), http.StatusBadRequest)).Send(ctx)
		return
	}

	req.DocumentId = documentId
	req.UserId = userId
	comment, err := c.commentService.Create(ctx, req)
	if err != nil {
		response.NewFailed("failed create comment", err).Send(ctx)
		return
	}

	response.NewSuccess("success create comment", comment).Send(ctx)
}

func (c *commentController) ReplyId(ctx *gin.Context) {
	documentId := ctx.Param("document_id")
	commentId := ctx.Param("comment_id")

	userId, err := utils.GetUserIdFromCtx(ctx)
	if err != nil {
		response.NewFailed("failed get data from body", myerror.New(err.Error(), http.StatusBadRequest)).Send(ctx)
		return
	}

	var req dto.CommentRequest
	if err := ctx.ShouldBind(&req); err != nil {
		response.NewFailed("failed get data from body", myerror.New(err.Error(), http.StatusBadRequest)).Send(ctx)
		return
	}

	req.DocumentId = documentId
	req.UserId = userId
	req.ReplyId = commentId

	comment, err := c.commentService.Reply(ctx, req)
	if err != nil {
		response.NewFailed("failed create comment", err).Send(ctx)
		return
	}

	response.NewSuccess("success create comment", comment).Send(ctx)
}

func (c *commentController) GetAllByDocumentId(ctx *gin.Context) {
	documentId := ctx.Param("document_id")

	userId, err := utils.GetUserIdFromCtx(ctx)
	if err != nil {
		response.NewFailed("failed get data from body", myerror.New(err.Error(), http.StatusBadRequest)).Send(ctx)
		return
	}

	comments, metaRes, err := c.commentService.GetAllByDocumentId(ctx.Request.Context(), userId, documentId, meta.New(ctx))
	if err != nil {
		response.NewFailed("failed get all comments", err).Send(ctx)
		return
	}

	response.NewSuccess("success get all comments", comments, metaRes).Send(ctx)
}

func (c *commentController) GetAllReplyByCommentId(ctx *gin.Context) {
	commentId := ctx.Param("comment_id")
	documentId := ctx.Param("document_id")
	userId, err := utils.GetUserIdFromCtx(ctx)
	if err != nil {
		response.NewFailed("failed get data from body", myerror.New(err.Error(), http.StatusBadRequest)).Send(ctx)
		return
	}
	comments, metaRes, err := c.commentService.GetAllByReplyId(ctx.Request.Context(), userId, documentId, commentId, meta.New(ctx))
	if err != nil {
		response.NewFailed("failed get all comments", err).Send(ctx)
		return
	}

	response.NewSuccess("success get all comments", comments, metaRes).Send(ctx)
}

func (c *commentController) GetById(ctx *gin.Context) {
	commentId := ctx.Param("comment_id")
	result, err := c.commentService.GetById(ctx.Request.Context(), commentId)
	if err != nil {
		response.NewFailed("failed get detail comment", err).Send(ctx)
		return
	}

	response.NewSuccess("success get detail comment", result).Send(ctx)
}

func (c *commentController) Update(ctx *gin.Context) {
	commentId := ctx.Param("comment_id")
	documentId := ctx.Param("document_id")
	userId, err := utils.GetUserIdFromCtx(ctx)
	if err != nil {
		response.NewFailed("failed get data from body", myerror.New(err.Error(), http.StatusBadRequest)).Send(ctx)
		return
	}

	var req dto.UpdateCommentRequest
	if err := ctx.ShouldBind(&req); err != nil {
		response.NewFailed("failed get data from body", myerror.New(err.Error(), http.StatusBadRequest)).Send(ctx)
		return
	}

	req.ID = commentId
	req.UserId = userId
	req.DocumentId = documentId
	err = c.commentService.Update(ctx.Request.Context(), req)
	if err != nil {
		response.NewFailed("failed update comment", err).Send(ctx)
		return
	}

	response.NewSuccess("success update comment", nil).Send(ctx)
}

func (c *commentController) Delete(ctx *gin.Context) {
	commentId := ctx.Param("comment_id")
	documentId := ctx.Param("document_id")
	userId, err := utils.GetUserIdFromCtx(ctx)
	if err != nil {
		response.NewFailed("failed get data from body", myerror.New(err.Error(), http.StatusBadRequest)).Send(ctx)
		return
	}

	err = c.commentService.Delete(ctx.Request.Context(), userId, documentId, commentId)
	if err != nil {
		response.NewFailed("failed delete comment", err).Send(ctx)
		return
	}

	response.NewSuccess("success delete comment", nil).Send(ctx)
}
