package repository

import (
	"context"

	"github.com/CRS-Project/crs-backend/internal/entity"
	"github.com/CRS-Project/crs-backend/internal/pkg/meta"
	"gorm.io/gorm"
)

type (
	CommentRepository interface {
		Create(ctx context.Context, tx *gorm.DB, comment entity.Comment, preloads ...string) (entity.Comment, error)
		GetByID(ctx context.Context, tx *gorm.DB, commentID string, preloads ...string) (entity.Comment, error)
		GetAll(ctx context.Context, tx *gorm.DB, metaReq meta.Meta, preloads ...string) ([]entity.Comment, meta.Meta, error)
		GetAllByDocumentID(ctx context.Context, tx *gorm.DB, documentId string, metaReq meta.Meta, preloads ...string) ([]entity.Comment, meta.Meta, error)
		GetAllByAreaOfConcernID(ctx context.Context, tx *gorm.DB, areaOfConcernId string, metaReq meta.Meta, preloads ...string) ([]entity.Comment, meta.Meta, error)
		GetAllByReplyID(ctx context.Context, tx *gorm.DB, replyId string, metaReq meta.Meta, preloads ...string) ([]entity.Comment, meta.Meta, error)
		Update(ctx context.Context, tx *gorm.DB, comment entity.Comment, preloads ...string) error
		Delete(ctx context.Context, tx *gorm.DB, comment entity.Comment, preloads ...string) error
	}

	commentRepository struct {
		db *gorm.DB
	}
)

func NewComment(db *gorm.DB) CommentRepository {
	return &commentRepository{
		db: db,
	}
}

func (r *commentRepository) Create(ctx context.Context, tx *gorm.DB, comment entity.Comment, preloads ...string) (entity.Comment, error) {
	if tx == nil {
		tx = r.db
	}

	for _, preload := range preloads {
		tx = tx.Preload(preload)
	}

	if err := tx.WithContext(ctx).Create(&comment).Error; err != nil {
		return entity.Comment{}, err
	}

	return comment, nil
}

func (r *commentRepository) GetByID(ctx context.Context, tx *gorm.DB, commentID string, preloads ...string) (entity.Comment, error) {
	if tx == nil {
		tx = r.db
	}

	for _, preload := range preloads {
		tx = tx.Preload(preload)
	}

	var comment entity.Comment
	if err := tx.WithContext(ctx).Take(&comment, "id = ?", commentID).Error; err != nil {
		return entity.Comment{}, err
	}

	return comment, nil
}

func (r *commentRepository) GetAll(ctx context.Context, tx *gorm.DB, metaReq meta.Meta, preloads ...string) ([]entity.Comment, meta.Meta, error) {
	if tx == nil {
		tx = r.db
	}

	for _, preload := range preloads {
		tx = tx.Preload(preload)
	}

	var comments []entity.Comment

	tx = tx.WithContext(ctx).Model(&entity.Comment{})
	if err := WithFilters(tx, &metaReq, AddModels(entity.Comment{})).Find(&comments).Error; err != nil {
		return nil, meta.Meta{}, err
	}

	return comments, metaReq, nil
}

func (r *commentRepository) GetAllByDocumentID(ctx context.Context, tx *gorm.DB, documentId string, metaReq meta.Meta, preloads ...string) ([]entity.Comment, meta.Meta, error) {
	if tx == nil {
		tx = r.db
	}

	for _, preload := range preloads {
		tx = tx.Preload(preload)
	}

	var comments []entity.Comment

	tx = tx.WithContext(ctx).Model(&entity.Comment{}).Where("document_id = ? AND comment_reply_id IS NULL", documentId)
	if err := WithFilters(tx, &metaReq, AddModels(entity.Comment{})).Find(&comments).Error; err != nil {
		return nil, meta.Meta{}, err
	}

	return comments, metaReq, nil
}

func (r *commentRepository) GetAllByAreaOfConcernID(ctx context.Context, tx *gorm.DB, areaOfConcernId string, metaReq meta.Meta, preloads ...string) ([]entity.Comment, meta.Meta, error) {
	if tx == nil {
		tx = r.db
	}

	for _, preload := range preloads {
		tx = tx.Preload(preload)
	}

	var comments []entity.Comment

	tx = tx.WithContext(ctx).Model(&entity.Comment{}).Where("area_of_concern_id = ? AND comment_reply_id IS NULL", areaOfConcernId)
	if err := WithFilters(tx, &metaReq, AddModels(entity.Comment{})).Find(&comments).Error; err != nil {
		return nil, meta.Meta{}, err
	}

	return comments, metaReq, nil
}

func (r *commentRepository) GetAllByReplyID(ctx context.Context, tx *gorm.DB, replyId string, metaReq meta.Meta, preloads ...string) ([]entity.Comment, meta.Meta, error) {
	if tx == nil {
		tx = r.db
	}

	for _, preload := range preloads {
		tx = tx.Preload(preload)
	}

	var comments []entity.Comment

	tx = tx.WithContext(ctx).Model(&entity.Comment{}).Where("comment_reply_id = ?", replyId)
	if err := WithFilters(tx, &metaReq, AddModels(entity.Comment{})).Find(&comments).Error; err != nil {
		return nil, meta.Meta{}, err
	}

	return comments, metaReq, nil
}

func (r *commentRepository) Update(ctx context.Context, tx *gorm.DB, comment entity.Comment, preloads ...string) error {
	if tx == nil {
		tx = r.db
	}

	for _, preload := range preloads {
		tx = tx.Preload(preload)
	}

	if err := tx.WithContext(ctx).Save(&comment).Error; err != nil {
		return err
	}

	return nil
}

func (r *commentRepository) Delete(ctx context.Context, tx *gorm.DB, comment entity.Comment, preloads ...string) error {
	if tx == nil {
		tx = r.db
	}

	for _, preload := range preloads {
		tx = tx.Preload(preload)
	}

	if err := tx.WithContext(ctx).Delete(&comment).Error; err != nil {
		return err
	}

	return nil
}
