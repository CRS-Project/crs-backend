package service

// import (
// 	"context"
// 	"net/http"

// 	"github.com/CRS-Project/crs-backend/internal/api/repository"
// 	"github.com/CRS-Project/crs-backend/internal/entity"
// 	myerror "github.com/CRS-Project/crs-backend/internal/pkg/error"
// 	"github.com/CRS-Project/crs-backend/internal/pkg/meta"
// 	"gorm.io/gorm"
// )

// type (
// 	CommentService interface {
// 		Create(ctx context.Context, req dto.CreateCommentRequest) (dto.CommentInfo, error)
// 		GetAllByDocumentID(ctx context.Context, metaReq meta.Meta) ([]dto.CommentInfo, meta.Meta, error)
// 		Update(ctx context.Context, req dto.UpdateCommentRequest) error
// 		Delete(ctx context.Context, req dto.DeleteCommentRequest) error
// 	}

// 	commentService struct {
// 		commentRepository repository.CommentRepository
// 		db                *gorm.DB
// 	}
// )

// func NewComment(commentRepository repository.CommentRepository, db *gorm.DB) CommentService {
// 	return &commentService{
// 		commentRepository: commentRepository,
// 		db:                db,
// 	}
// }

// func (s *commentService) CreateComment(ctx context.Context, req dto.CreateCommentRequest) (dto.CommentInfo, error) {
// 	_, err := s.commentRepository.GetByName(ctx, nil, req.Name)
// 	if err == nil {
// 		return dto.CommentInfo{}, myerror.New("comment with this name already exists", http.StatusConflict)
// 	}

// 	pkgCreation := entity.Comment{
// 		Name: req.Name,
// 	}

// 	pkgResult, err := s.commentRepository.Create(ctx, nil, pkgCreation)
// 	if err != nil {
// 		return dto.CommentInfo{}, err
// 	}

// 	return dto.CommentInfo{
// 		ID:   pkgResult.ID.String(),
// 		Name: pkgResult.Name,
// 	}, nil
// }

// func (s *commentService) GetAll(ctx context.Context, metaReq meta.Meta) ([]dto.CommentInfo, meta.Meta, error) {
// 	pkgs, metaRes, err := s.commentRepository.GetAll(ctx, nil, metaReq)
// 	if err != nil {
// 		return nil, meta.Meta{}, err
// 	}

// 	var pkgInfos []dto.CommentInfo
// 	for _, pkg := range pkgs {
// 		pkgInfos = append(pkgInfos, dto.CommentInfo{
// 			ID:   pkg.ID.String(),
// 			Name: pkg.Name,
// 		})
// 	}

// 	return pkgInfos, metaRes, nil
// }

// func (s *commentService) UpdateComment(ctx context.Context, req dto.UpdateCommentRequest) error {
// 	pkg, err := s.commentRepository.GetByID(ctx, nil, req.ID)
// 	if err != nil {
// 		return err
// 	}
// 	pkg.Name = req.Name

// 	if err = s.commentRepository.Update(ctx, nil, pkg); err != nil {
// 		return err
// 	}

// 	return nil
// }

// func (s *commentService) DeleteComment(ctx context.Context, req dto.DeleteCommentRequest) error {
// 	pkg, err := s.commentRepository.GetByID(ctx, nil, req.ID)
// 	if err != nil {
// 		return err
// 	}

// 	if err = s.commentRepository.Delete(ctx, nil, pkg); err != nil {
// 		return err
// 	}

// 	return nil
// }
