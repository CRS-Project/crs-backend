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
		GetAllByDisciplineListDocumentId(ctx context.Context, userId, disciplineListDocumentId string, metaReq meta.Meta) ([]dto.CommentResponse, meta.Meta, error)
		GetAllByReplyId(ctx context.Context, userId, disciplineListDocumentId, replyId string, metaReq meta.Meta) ([]dto.CommentResponse, meta.Meta, error)
		Update(ctx context.Context, req dto.UpdateCommentRequest) error
		Delete(ctx context.Context, userId, disciplineListDocumentId, commentId string) error
	}

	commentService struct {
		commentRepository                repository.CommentRepository
		documentRepository               repository.DocumentRepository
		disciplineListDocumentRepository repository.DisciplineListDocumentRepository
		userRepository                   repository.UserRepository
		db                               *gorm.DB
	}
)

func NewComment(commentRepository repository.CommentRepository,
	documentRepository repository.DocumentRepository,
	disciplineListDocumentRepository repository.DisciplineListDocumentRepository,
	userRepository repository.UserRepository,
	db *gorm.DB) CommentService {
	return &commentService{
		commentRepository:                commentRepository,
		documentRepository:               documentRepository,
		disciplineListDocumentRepository: disciplineListDocumentRepository,
		userRepository:                   userRepository,
		db:                               db,
	}
}

func (s *commentService) Create(ctx context.Context, req dto.CommentRequest) (dto.CommentResponse, error) {
	disciplineListDocument, _, err := s.checkPackagePermission(ctx, req.DisciplineListDocumentId, req.UserId)
	if err != nil {
		return dto.CommentResponse{}, err
	}

	if req.IsCloseOutComment {
		return dto.CommentResponse{}, myerror.New("you can't set this comment as close out comment", http.StatusBadRequest)
	}

	commentResult, err := s.commentRepository.Create(ctx, nil, entity.Comment{
		Section:                  req.Section,
		Comment:                  req.Comment,
		Baseline:                 req.Baseline,
		DisciplineListDocumentID: disciplineListDocument.ID,
		IsCloseOutComment:        req.IsCloseOutComment,
		AttachFileUrl:            req.AttachFileUrl,
		UserID:                   uuid.MustParse(req.UserId),
	})
	if err != nil {
		return dto.CommentResponse{}, err
	}

	return dto.CommentResponse{
		ID:                    commentResult.ID.String(),
		Section:               commentResult.Section,
		Comment:               commentResult.Comment,
		Baseline:              commentResult.Baseline,
		Status:                (*string)(commentResult.Status),
		AttachFileUrl:         commentResult.AttachFileUrl,
		CommentAt:             commentResult.CreatedAt.Format("15.04 • 02 Jan 2006"),
		CompanyDocumentNumber: disciplineListDocument.Document.CompanyDocumentNumber,
	}, nil
}

func (s *commentService) Reply(ctx context.Context, req dto.CommentRequest) (dto.CommentResponse, error) {
	disciplineListDocument, user, err := s.checkPackagePermission(ctx, req.DisciplineListDocumentId, req.UserId)
	if err != nil {
		return dto.CommentResponse{}, err
	}

	commentReplied, err := s.commentRepository.GetByID(ctx, nil, req.ReplyId)
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

	replyId := commentReplied.ID
	commentResult, err := s.commentRepository.Create(ctx, nil, entity.Comment{
		Section:                  req.Section,
		Comment:                  req.Comment,
		Baseline:                 req.Baseline,
		UserID:                   user.ID,
		IsCloseOutComment:        req.IsCloseOutComment,
		DisciplineListDocumentID: disciplineListDocument.ID,
		AttachFileUrl:            req.AttachFileUrl,
		CommentReplyID:           &replyId,
	})
	if err != nil {
		return dto.CommentResponse{}, err
	}

	return dto.CommentResponse{
		ID:                    commentResult.ID.String(),
		Section:               commentResult.Section,
		Comment:               commentResult.Comment,
		Baseline:              commentResult.Baseline,
		Status:                (*string)(commentResult.Status),
		AttachFileUrl:         commentResult.AttachFileUrl,
		CommentAt:             commentResult.CreatedAt.Format("15.04 • 02 Jan 2006"),
		CompanyDocumentNumber: disciplineListDocument.Document.CompanyDocumentNumber,
	}, nil
}

func (s *commentService) GetById(ctx context.Context, id string) (dto.CommentResponse, error) {
	comment, err := s.commentRepository.GetByID(ctx, nil, id, "User", "DisciplineListDocument.Document")
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
				CompanyDocumentNumber: comment.DisciplineListDocument.Document.CompanyDocumentNumber,
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
		DocumentID:            comment.DisciplineListDocument.Document.ID.String(),
		CommentAt:             comment.CreatedAt.Format("15.04 • 02 Jan 2006"),
		CompanyDocumentNumber: comment.DisciplineListDocument.Document.CompanyDocumentNumber,
		AttachFileUrl:         comment.AttachFileUrl,
		UserComment: &dto.UserComment{
			Name:         comment.User.Name,
			PhotoProfile: comment.User.PhotoProfile,
			Role:         string(comment.User.Role),
		},
		CommentReplies: replies,
	}, nil
}

func (s *commentService) GetAllByDisciplineListDocumentId(ctx context.Context, userId, disciplineListDocumentId string, metaReq meta.Meta) ([]dto.CommentResponse, meta.Meta, error) {
	disciplineListDocument, _, err := s.checkPackagePermission(ctx, disciplineListDocumentId, userId)
	if err != nil {
		return nil, meta.Meta{}, err
	}

	comments, metaRes, err := s.commentRepository.GetAllByDisciplineListDocumentID(ctx, nil, disciplineListDocumentId, metaReq, "User", "CommentReplies.User", "CommentReplies")
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
					DocumentID:            disciplineListDocument.Document.ID.String(),
					IsCloseOutComment:     reply.IsCloseOutComment,
					AttachFileUrl:         reply.AttachFileUrl,
					CompanyDocumentNumber: disciplineListDocument.Document.CompanyDocumentNumber,
					UserComment: &dto.UserComment{
						ID:           reply.User.ID.String(),
						Name:         reply.User.Name,
						PhotoProfile: reply.User.PhotoProfile,
						Role:         string(reply.User.Role),
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
			DocumentID:            disciplineListDocument.Document.ID.String(),
			IsCloseOutComment:     comment.IsCloseOutComment,
			AttachFileUrl:         comment.AttachFileUrl,
			CompanyDocumentNumber: disciplineListDocument.Document.CompanyDocumentNumber,
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

func (s *commentService) GetAllByReplyId(ctx context.Context, userId, disciplineListDocumentId, replyId string, metaReq meta.Meta) ([]dto.CommentResponse, meta.Meta, error) {
	_, _, err := s.checkPackagePermission(ctx, disciplineListDocumentId, userId)
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
	_, user, err := s.checkPackagePermission(ctx, req.DisciplineListDocumentId, req.UserId)
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
	comment.AttachFileUrl = req.AttachFileUrl
	if err = s.commentRepository.Update(ctx, nil, comment); err != nil {
		return err
	}

	return nil
}

func (s *commentService) Delete(ctx context.Context, userId, disciplineListDocumentId, commentId string) error {
	_, user, err := s.checkPackagePermission(ctx, disciplineListDocumentId, userId)
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

func (s *commentService) checkPackagePermission(ctx context.Context, disciplineListDocumentId, userId string) (entity.DisciplineListDocument, entity.User, error) {
	disciplineListDocument, err := s.disciplineListDocumentRepository.GetByID(ctx, nil, disciplineListDocumentId, "Document")
	if err != nil {
		return entity.DisciplineListDocument{}, entity.User{}, err
	}

	user, err := s.userRepository.GetById(ctx, nil, userId)
	if err != nil {
		return entity.DisciplineListDocument{}, entity.User{}, err
	}

	if user.PackageID != nil && disciplineListDocument.PackageID != *user.PackageID {
		return entity.DisciplineListDocument{}, entity.User{}, myerror.New("you dont have permission in this package", http.StatusUnauthorized)
	}

	return disciplineListDocument, user, nil
}
