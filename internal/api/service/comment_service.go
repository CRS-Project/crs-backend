package service

import (
	"context"
	"net/http"

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
		GetAllByAreaOfConcernId(ctx context.Context, userId, areaOfConcernId string, metaReq meta.Meta) ([]dto.CommentResponse, meta.Meta, error)
		GetAllByReplyId(ctx context.Context, userId, areaOfConcernId, replyId string, metaReq meta.Meta) ([]dto.CommentResponse, meta.Meta, error)
		Update(ctx context.Context, req dto.UpdateCommentRequest) error
		Delete(ctx context.Context, userId, areaOfConcernId, commentId string) error
	}

	commentService struct {
		commentRepository       repository.CommentRepository
		documentRepository      repository.DocumentRepository
		areaOfConcernRepository repository.AreaOfConcernRepository
		userRepository          repository.UserRepository
		db                      *gorm.DB
	}
)

func NewComment(commentRepository repository.CommentRepository,
	documentRepository repository.DocumentRepository,
	areaOfConcernRepository repository.AreaOfConcernRepository,
	userRepository repository.UserRepository,
	db *gorm.DB) CommentService {
	return &commentService{
		commentRepository:       commentRepository,
		documentRepository:      documentRepository,
		areaOfConcernRepository: areaOfConcernRepository,
		userRepository:          userRepository,
		db:                      db,
	}
}

func (s *commentService) Create(ctx context.Context, req dto.CommentRequest) (dto.CommentResponse, error) {
	areaOfConcern, _, _, err := s.checkPackagePermission(ctx, req.AreaOfConcernId, req.DocumentId, req.UserId)
	if err != nil {
		return dto.CommentResponse{}, err
	}

	if req.IsCloseOutComment {
		return dto.CommentResponse{}, myerror.New("you can't set this comment as close out comment", http.StatusBadRequest)
	}

	commentResult, err := s.commentRepository.Create(ctx, nil, entity.Comment{
		Section:           req.Section,
		Comment:           req.Comment,
		Baseline:          req.Baseline,
		AreaOfConcernID:   areaOfConcern.ID,
		IsCloseOutComment: req.IsCloseOutComment,
		DocumentID:        uuid.MustParse(req.DocumentId),
		UserID:            uuid.MustParse(req.UserId),
	})
	if err != nil {
		return dto.CommentResponse{}, err
	}

	return dto.CommentResponse{
		ID:        commentResult.ID.String(),
		Section:   commentResult.Section,
		Comment:   commentResult.Comment,
		Baseline:  commentResult.Baseline,
		Status:    (*string)(commentResult.Status),
		CommentAt: commentResult.CreatedAt.Format("15.04 • 02 Jan 2006"),
	}, nil
}

func (s *commentService) Reply(ctx context.Context, req dto.CommentRequest) (dto.CommentResponse, error) {
	areaOfConcern, document, user, err := s.checkPackagePermission(ctx, req.AreaOfConcernId, req.DocumentId, req.UserId)
	if err != nil {
		return dto.CommentResponse{}, err
	}

	commentReplied, err := s.commentRepository.GetByID(ctx, nil, req.ReplyId, "Document")
	if err != nil {
		return dto.CommentResponse{}, err
	}

	if commentReplied.Status != nil {
		return dto.CommentResponse{}, myerror.New("this comment already has a status", http.StatusUnauthorized)
	}

	if req.IsCloseOutComment {
		cs := entity.CommentStatusReject
		commentReplied.Status = &cs
		if err = s.commentRepository.Update(ctx, nil, commentReplied); err != nil {
			return dto.CommentResponse{}, err
		}
	}

	if commentReplied.CommentReplyID != nil {
		req.ReplyId = commentReplied.CommentReplyID.String()
	}

	replyId := uuid.MustParse(req.ReplyId)
	commentResult, err := s.commentRepository.Create(ctx, nil, entity.Comment{
		Section:         req.Section,
		Comment:         req.Comment,
		Baseline:        req.Baseline,
		DocumentID:      document.ID,
		UserID:          user.ID,
		AreaOfConcernID: areaOfConcern.ID,
		CommentReplyID:  &replyId,
	}, "Document")
	if err != nil {
		return dto.CommentResponse{}, err
	}

	return dto.CommentResponse{
		ID:                    commentResult.ID.String(),
		Section:               commentResult.Section,
		Comment:               commentResult.Comment,
		Baseline:              commentResult.Baseline,
		Status:                (*string)(commentResult.Status),
		CommentAt:             commentResult.CreatedAt.Format("15.04 • 02 Jan 2006"),
		CompanyDocumentNumber: commentResult.Document.CompanyDocumentNumber,
	}, nil
}

func (s *commentService) GetById(ctx context.Context, id string) (dto.CommentResponse, error) {
	comment, err := s.commentRepository.GetByID(ctx, nil, id, "User", "Document")
	if err != nil {
		return dto.CommentResponse{}, err
	}

	var replies []dto.CommentResponse
	if len(comment.CommentReplies) > 0 {
		for _, reply := range comment.CommentReplies {
			replies = append(replies, dto.CommentResponse{
				ID:                    reply.ID.String(),
				Section:               reply.Section,
				Comment:               reply.Comment,
				Baseline:              reply.Baseline,
				Status:                (*string)(reply.Status),
				CommentAt:             reply.CreatedAt.Format("15.04 • 02 Jan 2006"),
				DocumentID:            reply.DocumentID.String(),
				CompanyDocumentNumber: comment.Document.CompanyDocumentNumber,
				UserComment: &dto.UserComment{
					ID:           comment.User.ID.String(),
					Name:         comment.User.Name,
					PhotoProfile: comment.User.PhotoProfile,
					Role:         string(comment.User.Role),
				},
			})
		}
	}

	return dto.CommentResponse{
		ID:                    comment.ID.String(),
		Section:               comment.Section,
		Comment:               comment.Comment,
		Baseline:              comment.Baseline,
		Status:                (*string)(comment.Status),
		DocumentID:            comment.Document.ID.String(),
		CommentAt:             comment.CreatedAt.Format("15.04 • 02 Jan 2006"),
		CompanyDocumentNumber: comment.Document.CompanyDocumentNumber,
		UserComment: &dto.UserComment{
			Name:         comment.User.Name,
			PhotoProfile: comment.User.PhotoProfile,
			Role:         string(comment.User.Role),
		},
		CommentReplies: replies,
	}, nil
}

func (s *commentService) GetAllByAreaOfConcernId(ctx context.Context, userId, areaOfConcernId string, metaReq meta.Meta) ([]dto.CommentResponse, meta.Meta, error) {
	_, _, _, err := s.checkPackagePermission(ctx, areaOfConcernId, "", userId)
	if err != nil {
		return nil, meta.Meta{}, err
	}

	comments, metaRes, err := s.commentRepository.GetAllByAreaOfConcernID(ctx, nil, areaOfConcernId, metaReq, "User", "CommentReplies.User", "Document")
	if err != nil {
		return nil, meta.Meta{}, err
	}

	var commentResponse []dto.CommentResponse
	for _, comment := range comments {
		var replies []dto.CommentResponse
		if len(comment.CommentReplies) > 0 {
			for _, reply := range comment.CommentReplies {
				replies = append(replies, dto.CommentResponse{
					ID:                    reply.ID.String(),
					Section:               reply.Section,
					Comment:               reply.Comment,
					Baseline:              reply.Baseline,
					Status:                (*string)(reply.Status),
					CommentAt:             reply.CreatedAt.Format("15.04 • 02 Jan 2006"),
					DocumentID:            reply.DocumentID.String(),
					CompanyDocumentNumber: comment.Document.CompanyDocumentNumber,
					UserComment: &dto.UserComment{
						ID:           comment.User.ID.String(),
						Name:         comment.User.Name,
						PhotoProfile: comment.User.PhotoProfile,
						Role:         string(comment.User.Role),
					},
				})
			}
		}

		commentResponse = append(commentResponse, dto.CommentResponse{
			ID:                    comment.ID.String(),
			Section:               comment.Section,
			Comment:               comment.Comment,
			Baseline:              comment.Baseline,
			Status:                (*string)(comment.Status),
			CommentAt:             comment.CreatedAt.Format("15.04 • 02 Jan 2006"),
			DocumentID:            comment.DocumentID.String(),
			CompanyDocumentNumber: comment.Document.CompanyDocumentNumber,
			UserComment: &dto.UserComment{
				ID:           comment.User.ID.String(),
				Name:         comment.User.Name,
				PhotoProfile: comment.User.PhotoProfile,
				Role:         string(comment.User.Role),
			},
			CommentReplies: replies,
		})
	}

	return commentResponse, metaRes, nil
}

func (s *commentService) GetAllByReplyId(ctx context.Context, userId, areaOfConcernId, replyId string, metaReq meta.Meta) ([]dto.CommentResponse, meta.Meta, error) {
	_, _, _, err := s.checkPackagePermission(ctx, areaOfConcernId, "", userId)
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
			Status:    (*string)(comment.Status),
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
	_, _, user, err := s.checkPackagePermission(ctx, req.AreaOfConcernId, req.DocumentId, req.UserId)
	if err != nil {
		return err
	}

	comment, err := s.commentRepository.GetByID(ctx, nil, req.ID)
	if err != nil {
		return err
	}

	if user.PackageID != nil && comment.UserID != user.ID {
		return myerror.New("you dont have permission in this comment", http.StatusUnauthorized)
	}

	if req.Status != nil && comment.CommentReplyID != nil {
		return myerror.New("this is not parent comment", http.StatusUnauthorized)
	}

	if comment.Status != nil {
		return myerror.New("this comment already has a status", http.StatusUnauthorized)
	}

	comment.Comment = req.Comment
	comment.Baseline = req.Baseline
	comment.Section = req.Section
	comment.Status = (*entity.CommentStatus)(req.Status)

	if err = s.commentRepository.Update(ctx, nil, comment); err != nil {
		return err
	}

	return nil
}

func (s *commentService) Delete(ctx context.Context, userId, areaOfConcernId, commentId string) error {
	_, _, user, err := s.checkPackagePermission(ctx, areaOfConcernId, "", userId)
	if err != nil {
		return err
	}

	comment, err := s.commentRepository.GetByID(ctx, nil, commentId)
	if err != nil {
		return err
	}

	if user.PackageID != nil && comment.UserID != user.ID {
		return myerror.New("you don't have permission for this comment", http.StatusUnauthorized)
	}

	if err = s.commentRepository.Delete(ctx, nil, comment); err != nil {
		return err
	}

	return nil
}

func (s *commentService) checkPackagePermission(ctx context.Context, areaOfConcernId, documentId, userId string) (entity.AreaOfConcern, *entity.Document, entity.User, error) {
	areaOfConcern, err := s.areaOfConcernRepository.GetByID(ctx, nil, areaOfConcernId)
	if err != nil {
		return entity.AreaOfConcern{}, nil, entity.User{}, err
	}

	var document *entity.Document
	if documentId != "" {
		documentR, err := s.documentRepository.GetByID(ctx, nil, documentId)
		if err != nil {
			return entity.AreaOfConcern{}, nil, entity.User{}, err
		}

		document = &documentR
	}

	user, err := s.userRepository.GetById(ctx, nil, userId)
	if err != nil {
		return entity.AreaOfConcern{}, nil, entity.User{}, err
	}

	if user.PackageID != nil && areaOfConcern.PackageID != *user.PackageID {
		return entity.AreaOfConcern{}, nil, entity.User{}, myerror.New("you dont have permission in this package", http.StatusUnauthorized)
	}

	return areaOfConcern, document, user, nil
}
