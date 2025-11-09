package repository

import (
	"context"

	"github.com/CRS-Project/crs-backend/internal/entity"
	"github.com/CRS-Project/crs-backend/internal/pkg/meta"
	"gorm.io/gorm"
)

type (
	AreaOfConcernConsolidatorRepository interface {
		Create(ctx context.Context, tx *gorm.DB, areaOfConcernConsolidator entity.AreaOfConcernConsolidator, preloads ...string) (entity.AreaOfConcernConsolidator, error)
		CreateBulk(ctx context.Context, tx *gorm.DB, areaOfConcernConsolidator []entity.AreaOfConcernConsolidator, preloads ...string) error
		GetAll(ctx context.Context, tx *gorm.DB, metaReq meta.Meta, preloads ...string) ([]entity.AreaOfConcernConsolidator, meta.Meta, error)
		GetByID(ctx context.Context, tx *gorm.DB, areaOfConcernConsolidatorID string, preloads ...string) (entity.AreaOfConcernConsolidator, error)
		Update(ctx context.Context, tx *gorm.DB, areaOfConcernConsolidator entity.AreaOfConcernConsolidator, preloads ...string) error
		Delete(ctx context.Context, tx *gorm.DB, areaOfConcernConsolidator entity.AreaOfConcernConsolidator, preloads ...string) error
		DeleteBulk(ctx context.Context, tx *gorm.DB, areaOfConcernConcolidatorIDs []string) error
	}

	areaOfConcernConsolidatorRepository struct {
		db *gorm.DB
	}
)

func NewAreaOfConcernConsolidator(db *gorm.DB) AreaOfConcernConsolidatorRepository {
	return &areaOfConcernConsolidatorRepository{
		db: db,
	}
}

func (r *areaOfConcernConsolidatorRepository) Create(ctx context.Context, tx *gorm.DB, areaOfConcernConsolidator entity.AreaOfConcernConsolidator, preloads ...string) (entity.AreaOfConcernConsolidator, error) {
	if tx == nil {
		tx = r.db
	}

	for _, preload := range preloads {
		tx = tx.Preload(preload)
	}

	if err := tx.WithContext(ctx).Create(&areaOfConcernConsolidator).Error; err != nil {
		return entity.AreaOfConcernConsolidator{}, err
	}

	return areaOfConcernConsolidator, nil
}

func (r *areaOfConcernConsolidatorRepository) CreateBulk(ctx context.Context, tx *gorm.DB, areaOfConcernConsolidators []entity.AreaOfConcernConsolidator, preloads ...string) error {
	if tx == nil {
		tx = r.db
	}

	for _, preload := range preloads {
		tx = tx.Preload(preload)
	}

	if err := tx.WithContext(ctx).Create(&areaOfConcernConsolidators).Error; err != nil {
		return err
	}

	return nil
}

func (r *areaOfConcernConsolidatorRepository) GetAll(ctx context.Context, tx *gorm.DB, metaReq meta.Meta, preloads ...string) ([]entity.AreaOfConcernConsolidator, meta.Meta, error) {
	if tx == nil {
		tx = r.db
	}

	for _, preload := range preloads {
		tx = tx.Preload(preload)
	}

	var areaOfConcernConsolidators []entity.AreaOfConcernConsolidator

	tx = tx.WithContext(ctx).Model(&entity.AreaOfConcernConsolidator{})
	if err := WithFilters(tx, &metaReq, AddModels(entity.AreaOfConcernConsolidator{})).Find(&areaOfConcernConsolidators).Error; err != nil {
		return nil, meta.Meta{}, err
	}

	return areaOfConcernConsolidators, metaReq, nil
}

func (r *areaOfConcernConsolidatorRepository) GetByID(ctx context.Context, tx *gorm.DB, areaOfConcernConsolidatorID string, preloads ...string) (entity.AreaOfConcernConsolidator, error) {
	if tx == nil {
		tx = r.db
	}

	for _, preload := range preloads {
		tx = tx.Preload(preload)
	}

	var areaOfConcernConsolidator entity.AreaOfConcernConsolidator
	if err := tx.WithContext(ctx).First(&areaOfConcernConsolidator, "id = ?", areaOfConcernConsolidatorID).Error; err != nil {
		return entity.AreaOfConcernConsolidator{}, err
	}

	return areaOfConcernConsolidator, nil
}

func (r *areaOfConcernConsolidatorRepository) Update(ctx context.Context, tx *gorm.DB, areaOfConcernConsolidator entity.AreaOfConcernConsolidator, preloads ...string) error {
	if tx == nil {
		tx = r.db
	}

	for _, preload := range preloads {
		tx = tx.Preload(preload)
	}

	if err := tx.WithContext(ctx).Save(&areaOfConcernConsolidator).Error; err != nil {
		return err
	}

	return nil
}

func (r *areaOfConcernConsolidatorRepository) Delete(ctx context.Context, tx *gorm.DB, areaOfConcernConsolidator entity.AreaOfConcernConsolidator, preloads ...string) error {
	if tx == nil {
		tx = r.db
	}

	for _, preload := range preloads {
		tx = tx.Preload(preload)
	}

	if err := tx.WithContext(ctx).Delete(&areaOfConcernConsolidator).Error; err != nil {
		return err
	}

	return nil
}

func (r *areaOfConcernConsolidatorRepository) DeleteBulk(ctx context.Context, tx *gorm.DB, areaOfConcernConcolidatorIDs []string) error {
	if tx == nil {
		tx = r.db
	}

	if err := tx.WithContext(ctx).
		Where("id IN ?", areaOfConcernConcolidatorIDs).
		Delete(&entity.AreaOfConcernConsolidator{}).Error; err != nil {
		return err
	}

	return nil
}
