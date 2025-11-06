package service

import (
	"context"
	"net/http"
	"time"

	"github.com/CRS-Project/crs-backend/internal/api/repository"
	"github.com/CRS-Project/crs-backend/internal/dto"
	"github.com/CRS-Project/crs-backend/internal/entity"
	myerror "github.com/CRS-Project/crs-backend/internal/pkg/error"
	"github.com/CRS-Project/crs-backend/internal/pkg/meta"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type (
	CommentService interface {
		Create(ctx context.Context, req dto.CommentRequest) (dto.CommentResponse, error)
		Reply(ctx context.Context, req dto.CommentRequest) (dto.CommentResponse, error)
		GetById(ctx context.Context, id string) (dto.CommentResponse, error)
		GetAllByDocumentId(ctx context.Context, userId, documentId string, metaReq meta.Meta) ([]dto.CommentResponse, meta.Meta, error)
		GetAllByReplyId(ctx context.Context, userId, documentId, replyId string, metaReq meta.Meta) ([]dto.CommentResponse, meta.Meta, error)
		Update(ctx context.Context, req dto.UpdateCommentRequest) error
		Delete(ctx context.Context, userId, documentId, commentId string) error
	}

	commentService struct {
		commentRepository  repository.CommentRepository
		documentRepository repository.DocumentRepository
		userRepository     repository.UserRepository
		db                 *gorm.DB
	}
)

func NewComment(commentRepository repository.CommentRepository,
	documentRepository repository.DocumentRepository,
	userRepository repository.UserRepository,
	db *gorm.DB) CommentService {
	return &commentService{
		commentRepository:  commentRepository,
		documentRepository: documentRepository,
		userRepository:     userRepository,
		db:                 db,
	}
}

func (s *commentService) Create(ctx context.Context, req dto.CommentRequest) (dto.CommentResponse, error) {
	document, _, err := s.checkPackagePermission(ctx, req.DocumentId, req.UserId)
	if err != nil {
		return dto.CommentResponse{}, err
	}

	if time.Now().After(document.Deadline) {
		return dto.CommentResponse{}, myerror.New("document deadline has passed", http.StatusBadRequest)
	}

	commentResult, err := s.commentRepository.Create(ctx, nil, entity.Comment{
		Section:    req.Section,
		Comment:    req.Comment,
		Baseline:   req.Baseline,
		DocumentID: uuid.MustParse(req.DocumentId),
		UserID:     uuid.MustParse(req.UserId),
	})
	if err != nil {
		return dto.CommentResponse{}, err
	}

	return dto.CommentResponse{
		ID:        commentResult.ID.String(),
		Section:   commentResult.Section,
		Comment:   commentResult.Comment,
		Baseline:  commentResult.Baseline,
		Status:    string(commentResult.Status),
		CommentAt: commentResult.CreatedAt.Format("15.04 • 02 Jan 2006"),
	}, nil
}

func (s *commentService) Reply(ctx context.Context, req dto.CommentRequest) (dto.CommentResponse, error) {
	document, user, err := s.checkPackagePermission(ctx, req.DocumentId, req.UserId)
	if err != nil {
		return dto.CommentResponse{}, err
	}

	if time.Now().After(document.Deadline) {
		return dto.CommentResponse{}, myerror.New("document deadline has passed", http.StatusBadRequest)
	}

	commentReplied, err := s.commentRepository.GetByID(ctx, nil, req.ReplyId)
	if err != nil {
		return dto.CommentResponse{}, err
	}

	if commentReplied.CommentReplyID != nil {
		req.ReplyId = commentReplied.CommentReplyID.String()
	}

	replyId := uuid.MustParse(req.ReplyId)
	commentResult, err := s.commentRepository.Create(ctx, nil, entity.Comment{
		Section:        req.Section,
		Comment:        req.Comment,
		Baseline:       req.Baseline,
		DocumentID:     document.ID,
		UserID:         user.ID,
		CommentReplyID: &replyId,
	})
	if err != nil {
		return dto.CommentResponse{}, err
	}

	return dto.CommentResponse{
		ID:        commentResult.ID.String(),
		Section:   commentResult.Section,
		Comment:   commentResult.Comment,
		Baseline:  commentResult.Baseline,
		Status:    string(commentResult.Status),
		CommentAt: commentResult.CreatedAt.Format("15.04 • 02 Jan 2006"),
	}, nil
}

func (s *commentService) GetById(ctx context.Context, id string) (dto.CommentResponse, error) {
	comment, err := s.commentRepository.GetByID(ctx, nil, id, "User")
	if err != nil {
		return dto.CommentResponse{}, err
	}

	return dto.CommentResponse{
		ID:        comment.ID.String(),
		Section:   comment.Section,
		Comment:   comment.Comment,
		Baseline:  comment.Baseline,
		Status:    string(comment.Status),
		CommentAt: comment.CreatedAt.Format("15.04 • 02 Jan 2006"),
		UserComment: &dto.UserComment{
			Name: comment.User.Name,
			Role: string(comment.User.Role),
		},
	}, nil
}

func (s *commentService) GetAllByDocumentId(ctx context.Context, userId, documentId string, metaReq meta.Meta) ([]dto.CommentResponse, meta.Meta, error) {
	_, _, err := s.checkPackagePermission(ctx, documentId, userId)
	if err != nil {
		return nil, meta.Meta{}, err
	}

	comments, metaRes, err := s.commentRepository.GetAllByDocumentID(ctx, nil, documentId, metaReq, "User")
	if err != nil {
		return nil, meta.Meta{}, err
	}

	var commentResponse []dto.CommentResponse
	for _, comment := range comments {
		commentResponse = append(commentResponse, dto.CommentResponse{
			ID:        comment.ID.String(),
			Section:   comment.Section,
			Comment:   comment.Comment,
			Baseline:  comment.Baseline,
			Status:    string(comment.Status),
			CommentAt: comment.CreatedAt.Format("15.04 • 02 Jan 2006"),
			UserComment: &dto.UserComment{
				Name: comment.User.Name,
				Role: string(comment.User.Role),
			},
		})
	}

	return commentResponse, metaRes, nil
}

func (s *commentService) GetAllByReplyId(ctx context.Context, userId, documentId, replyId string, metaReq meta.Meta) ([]dto.CommentResponse, meta.Meta, error) {
	_, _, err := s.checkPackagePermission(ctx, documentId, userId)
	if err != nil {
		return nil, meta.Meta{}, err
	}

	comments, metaRes, err := s.commentRepository.GetAllByReplyID(ctx, nil, replyId, metaReq, "User")
	if err != nil {
		return nil, meta.Meta{}, err
	}

	var commentResponse []dto.CommentResponse
	for _, comment := range comments {
		commentResponse = append(commentResponse, dto.CommentResponse{
			ID:        comment.ID.String(),
			Section:   comment.Section,
			Comment:   comment.Comment,
			Baseline:  comment.Baseline,
			Status:    string(comment.Status),
			CommentAt: comment.CreatedAt.Format("15.04 • 02 Jan 2006"),
			UserComment: &dto.UserComment{
				Name: comment.User.Name,
				Role: string(comment.User.Role),
			},
		})
	}

	return commentResponse, metaRes, nil
}

func (s *commentService) Update(ctx context.Context, req dto.UpdateCommentRequest) error {
	document, user, err := s.checkPackagePermission(ctx, req.DocumentId, req.UserId)
	if err != nil {
		return err
	}

	if time.Now().After(document.Deadline) {
		return myerror.New("document deadline has passed", http.StatusBadRequest)
	}

	comment, err := s.commentRepository.GetByID(ctx, nil, req.ID)
	if err != nil {
		return err
	}

	if user.PackageID != nil && comment.UserID != user.ID {
		return myerror.New("you dont have permission in this comment", http.StatusUnauthorized)
	}

	if req.Status == string(entity.CommentStatusClose) && comment.CommentReplyID != nil {
		return myerror.New("this is not parent comment", http.StatusUnauthorized)
	}

	comment.Comment = req.Comment
	comment.Baseline = req.Baseline
	comment.Section = req.Section
	comment.Status = entity.CommentStatus(req.Status)

	if err = s.commentRepository.Update(ctx, nil, comment); err != nil {
		return err
	}

	return nil
}

func (s *commentService) Delete(ctx context.Context, userId, documentId, commentId string) error {
	_, user, err := s.checkPackagePermission(ctx, documentId, userId)
	if err != nil {
		return err
	}

	comment, err := s.commentRepository.GetByID(ctx, nil, commentId)
	if err != nil {
		return err
	}

	if user.PackageID != nil && comment.UserID != user.ID {
		return myerror.New("you dont have permission in this comment", http.StatusUnauthorized)
	}

	if err = s.commentRepository.Delete(ctx, nil, comment); err != nil {
		return err
	}

	return nil
}

func (s *commentService) checkPackagePermission(ctx context.Context, documentId, userId string) (entity.Document, entity.User, error) {
	document, err := s.documentRepository.GetByID(ctx, nil, documentId)
	if err != nil {
		return entity.Document{}, entity.User{}, err
	}

	user, err := s.userRepository.GetById(ctx, nil, userId)
	if err != nil {
		return entity.Document{}, entity.User{}, err
	}

	if user.PackageID != nil && document.PackageID != *user.PackageID {
		return entity.Document{}, entity.User{}, myerror.New("you dont have permission in this package", http.StatusUnauthorized)
	}

	return document, user, nil
}
