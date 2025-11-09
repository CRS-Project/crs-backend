package repository

import (
	"context"

	"github.com/CRS-Project/crs-backend/internal/entity"
	"github.com/CRS-Project/crs-backend/internal/pkg/meta"
	"gorm.io/gorm"
)

type (
	AreaOfConcernRepository interface {
		Create(ctx context.Context, tx *gorm.DB, areaOfConcern entity.AreaOfConcern, preloads ...string) (entity.AreaOfConcern, error)
		GetAll(ctx context.Context, tx *gorm.DB, metaReq meta.Meta, preloads ...string) ([]entity.AreaOfConcern, meta.Meta, error)
		GetAllByAreaOfConcernGroupID(ctx context.Context, tx *gorm.DB, areaOfConcernGroupId string, metaReq meta.Meta, preloads ...string) ([]entity.AreaOfConcern, meta.Meta, error)
		GetByID(ctx context.Context, tx *gorm.DB, areaOfConcernID string, preloads ...string) (entity.AreaOfConcern, error)
		Update(ctx context.Context, tx *gorm.DB, areaOfConcern entity.AreaOfConcern, preloads ...string) error
		Delete(ctx context.Context, tx *gorm.DB, areaOfConcern entity.AreaOfConcern, preloads ...string) error
	}

	areaOfConcernRepository struct {
		db *gorm.DB
	}
)

func NewAreaOfConcern(db *gorm.DB) AreaOfConcernRepository {
	return &areaOfConcernRepository{
		db: db,
	}
}

func (r *areaOfConcernRepository) Create(ctx context.Context, tx *gorm.DB, areaOfConcern entity.AreaOfConcern, preloads ...string) (entity.AreaOfConcern, error) {
	if tx == nil {
		tx = r.db
	}

	for _, preload := range preloads {
		tx = tx.Preload(preload)
	}

	if err := tx.WithContext(ctx).Create(&areaOfConcern).Error; err != nil {
		return entity.AreaOfConcern{}, err
	}

	return areaOfConcern, nil
}

func (r *areaOfConcernRepository) GetAll(ctx context.Context, tx *gorm.DB, metaReq meta.Meta, preloads ...string) ([]entity.AreaOfConcern, meta.Meta, error) {
	if tx == nil {
		tx = r.db
	}

	for _, preload := range preloads {
		tx = tx.Preload(preload)
	}

	var areaOfConcerns []entity.AreaOfConcern

	tx = tx.WithContext(ctx).Model(&entity.AreaOfConcern{})

	if err := WithFilters(tx, &metaReq, AddModels(entity.AreaOfConcern{})).Find(&areaOfConcerns).Error; err != nil {
		return nil, meta.Meta{}, err
	}

	return areaOfConcerns, metaReq, nil
}

func (r *areaOfConcernRepository) GetAllByAreaOfConcernGroupID(ctx context.Context, tx *gorm.DB, areaOfConcernGroupId string, metaReq meta.Meta, preloads ...string) ([]entity.AreaOfConcern, meta.Meta, error) {
	if tx == nil {
		tx = r.db
	}

	for _, preload := range preloads {
		tx = tx.Preload(preload)
	}

	var areaOfConcerns []entity.AreaOfConcern

	tx = tx.WithContext(ctx).Model(&entity.AreaOfConcern{}).Where("area_of_concern_group_id = ?", areaOfConcernGroupId)
	if err := WithFilters(tx, &metaReq, AddModels(entity.AreaOfConcern{})).Find(&areaOfConcerns).Error; err != nil {
		return nil, meta.Meta{}, err
	}

	return areaOfConcerns, metaReq, nil
}

func (r *areaOfConcernRepository) GetByID(ctx context.Context, tx *gorm.DB, areaOfConcernID string, preloads ...string) (entity.AreaOfConcern, error) {
	if tx == nil {
		tx = r.db
	}

	for _, preload := range preloads {
		tx = tx.Preload(preload)
	}

	var areaOfConcern entity.AreaOfConcern
	if err := tx.WithContext(ctx).First(&areaOfConcern, "id = ?", areaOfConcernID).Error; err != nil {
		return entity.AreaOfConcern{}, err
	}

	return areaOfConcern, nil
}

func (r *areaOfConcernRepository) Update(ctx context.Context, tx *gorm.DB, areaOfConcern entity.AreaOfConcern, preloads ...string) error {
	if tx == nil {
		tx = r.db
	}

	for _, preload := range preloads {
		tx = tx.Preload(preload)
	}

	if err := tx.WithContext(ctx).Save(&areaOfConcern).Error; err != nil {
		return err
	}

	return nil
}

func (r *areaOfConcernRepository) Delete(ctx context.Context, tx *gorm.DB, areaOfConcern entity.AreaOfConcern, preloads ...string) error {
	if tx == nil {
		tx = r.db
	}

	for _, preload := range preloads {
		tx = tx.Preload(preload)
	}

	if err := tx.WithContext(ctx).Delete(&areaOfConcern).Error; err != nil {
		return err
	}

	return nil
}
