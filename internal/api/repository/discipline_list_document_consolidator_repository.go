package repository

import (
	"context"

	"github.com/CRS-Project/crs-backend/internal/entity"
	"github.com/CRS-Project/crs-backend/internal/pkg/meta"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type (
	DisciplineListDocumentConsolidatorRepository interface {
		Create(ctx context.Context, tx *gorm.DB, disciplineListDocumentConsolidator entity.DisciplineListDocumentConsolidator, preloads ...string) (entity.DisciplineListDocumentConsolidator, error)
		CreateBulk(ctx context.Context, tx *gorm.DB, disciplineListDocumentConsolidator []entity.DisciplineListDocumentConsolidator, preloads ...string) error
		GetAll(ctx context.Context, tx *gorm.DB, metaReq meta.Meta, preloads ...string) ([]entity.DisciplineListDocumentConsolidator, meta.Meta, error)
		GetByID(ctx context.Context, tx *gorm.DB, disciplineListDocumentConsolidatorID string, preloads ...string) (entity.DisciplineListDocumentConsolidator, error)
		Update(ctx context.Context, tx *gorm.DB, disciplineListDocumentConsolidator entity.DisciplineListDocumentConsolidator, preloads ...string) error
		Delete(ctx context.Context, tx *gorm.DB, disciplineListDocumentConsolidator entity.DisciplineListDocumentConsolidator, preloads ...string) error
		DeleteByDisciplineListDocumentID(ctx context.Context, tx *gorm.DB, disciplineListDocumentID []string) error
		DeleteByDisciplineGroupConsolidatorID(ctx context.Context, tx *gorm.DB, disciplineGroupIDs []string) error
		DeleteBulk(ctx context.Context, tx *gorm.DB, disciplineListDocumentConcolidatorIDs []string) error
	}

	disciplineListDocumentConsolidatorRepository struct {
		db *gorm.DB
	}
)

func NewDisciplineListDocumentConsolidator(db *gorm.DB) DisciplineListDocumentConsolidatorRepository {
	return &disciplineListDocumentConsolidatorRepository{
		db: db,
	}
}

func (r *disciplineListDocumentConsolidatorRepository) Create(ctx context.Context, tx *gorm.DB, disciplineListDocumentConsolidator entity.DisciplineListDocumentConsolidator, preloads ...string) (entity.DisciplineListDocumentConsolidator, error) {
	if tx == nil {
		tx = r.db
	}

	for _, preload := range preloads {
		tx = tx.Preload(preload)
	}

	if err := tx.WithContext(ctx).Create(&disciplineListDocumentConsolidator).Error; err != nil {
		return entity.DisciplineListDocumentConsolidator{}, err
	}

	return disciplineListDocumentConsolidator, nil
}

func (r *disciplineListDocumentConsolidatorRepository) CreateBulk(ctx context.Context, tx *gorm.DB, disciplineListDocumentConsolidators []entity.DisciplineListDocumentConsolidator, preloads ...string) error {
	if tx == nil {
		tx = r.db
	}

	for _, preload := range preloads {
		tx = tx.Preload(preload)
	}

	if err := tx.WithContext(ctx).Create(&disciplineListDocumentConsolidators).Error; err != nil {
		return err
	}

	return nil
}

func (r *disciplineListDocumentConsolidatorRepository) GetAll(ctx context.Context, tx *gorm.DB, metaReq meta.Meta, preloads ...string) ([]entity.DisciplineListDocumentConsolidator, meta.Meta, error) {
	if tx == nil {
		tx = r.db
	}

	for _, preload := range preloads {
		tx = tx.Preload(preload)
	}

	var disciplineListDocumentConsolidators []entity.DisciplineListDocumentConsolidator

	tx = tx.WithContext(ctx).Model(&entity.DisciplineListDocumentConsolidator{})
	if err := WithFilters(tx, &metaReq, AddModels(entity.DisciplineListDocumentConsolidator{})).Find(&disciplineListDocumentConsolidators).Error; err != nil {
		return nil, meta.Meta{}, err
	}

	return disciplineListDocumentConsolidators, metaReq, nil
}

func (r *disciplineListDocumentConsolidatorRepository) GetByID(ctx context.Context, tx *gorm.DB, disciplineListDocumentConsolidatorID string, preloads ...string) (entity.DisciplineListDocumentConsolidator, error) {
	if tx == nil {
		tx = r.db
	}

	for _, preload := range preloads {
		tx = tx.Preload(preload)
	}

	var disciplineListDocumentConsolidator entity.DisciplineListDocumentConsolidator
	if err := tx.WithContext(ctx).First(&disciplineListDocumentConsolidator, "id = ?", disciplineListDocumentConsolidatorID).Error; err != nil {
		return entity.DisciplineListDocumentConsolidator{}, err
	}

	return disciplineListDocumentConsolidator, nil
}

func (r *disciplineListDocumentConsolidatorRepository) Update(ctx context.Context, tx *gorm.DB, disciplineListDocumentConsolidator entity.DisciplineListDocumentConsolidator, preloads ...string) error {
	if tx == nil {
		tx = r.db
	}

	for _, preload := range preloads {
		tx = tx.Preload(preload)
	}

	if err := tx.WithContext(ctx).Save(&disciplineListDocumentConsolidator).Error; err != nil {
		return err
	}

	return nil
}

func (r *disciplineListDocumentConsolidatorRepository) Delete(ctx context.Context, tx *gorm.DB, disciplineListDocumentConsolidator entity.DisciplineListDocumentConsolidator, preloads ...string) error {
	if tx == nil {
		tx = r.db
	}

	for _, preload := range preloads {
		tx = tx.Preload(preload)
	}

	// persist deleted_by if provided
	if disciplineListDocumentConsolidator.DeletedBy != uuid.Nil {
		if err := tx.WithContext(ctx).Model(&entity.DisciplineListDocumentConsolidator{}).
			Where("id = ?", disciplineListDocumentConsolidator.ID).
			Updates(map[string]interface{}{"deleted_by": disciplineListDocumentConsolidator.DeletedBy}).Error; err != nil {
			return err
		}
	}

	if err := tx.WithContext(ctx).Delete(&disciplineListDocumentConsolidator).Error; err != nil {
		return err
	}

	return nil
}

func (r *disciplineListDocumentConsolidatorRepository) DeleteByDisciplineListDocumentID(ctx context.Context, tx *gorm.DB, disciplineListDocumentID []string) error {
	if tx == nil {
		tx = r.db
	}

	if err := tx.WithContext(ctx).
		Where("discipline_list_document_id IN ?", disciplineListDocumentID).
		Delete(&entity.DisciplineListDocumentConsolidator{}).Error; err != nil {
		return err
	}

	return nil
}

func (r *disciplineListDocumentConsolidatorRepository) DeleteByDisciplineGroupConsolidatorID(ctx context.Context, tx *gorm.DB, disciplineGroupIDs []string) error {
	if tx == nil {
		tx = r.db
	}

	if err := tx.WithContext(ctx).
		Where("discipline_group_consolidator_id IN ?", disciplineGroupIDs).
		Delete(&entity.DisciplineListDocumentConsolidator{}).Error; err != nil {
		return err
	}

	return nil
}

func (r *disciplineListDocumentConsolidatorRepository) DeleteBulk(ctx context.Context, tx *gorm.DB, disciplineListDocumentConcolidatorIDs []string) error {
	if tx == nil {
		tx = r.db
	}

	if err := tx.WithContext(ctx).
		Where("id IN ?", disciplineListDocumentConcolidatorIDs).
		Delete(&entity.DisciplineListDocumentConsolidator{}).Error; err != nil {
		return err
	}

	return nil
}
