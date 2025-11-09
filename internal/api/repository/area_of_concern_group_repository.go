package repository

import (
	"context"

	"github.com/CRS-Project/crs-backend/internal/entity"
	"github.com/CRS-Project/crs-backend/internal/pkg/meta"
	"gorm.io/gorm"
)

type (
	AreaOfConcernGroupRepository interface {
		Create(ctx context.Context, tx *gorm.DB, areaOfConcernGroup entity.AreaOfConcernGroup, preloads ...string) (entity.AreaOfConcernGroup, error)
		GetAll(ctx context.Context, tx *gorm.DB, packageId string, metaReq meta.Meta, preloads ...string) ([]entity.AreaOfConcernGroup, meta.Meta, error)
		GetByID(ctx context.Context, tx *gorm.DB, areaOfConcernGroupID string, preloads ...string) (entity.AreaOfConcernGroup, error)
		Update(ctx context.Context, tx *gorm.DB, areaOfConcernGroup entity.AreaOfConcernGroup, preloads ...string) error
		Delete(ctx context.Context, tx *gorm.DB, areaOfConcernGroup entity.AreaOfConcernGroup, preloads ...string) error
	}

	areaOfConcernGroupRepository struct {
		db *gorm.DB
	}
)

func NewAreaOfConcernGroup(db *gorm.DB) AreaOfConcernGroupRepository {
	return &areaOfConcernGroupRepository{
		db: db,
	}
}

func (r *areaOfConcernGroupRepository) Create(ctx context.Context, tx *gorm.DB, areaOfConcernGroup entity.AreaOfConcernGroup, preloads ...string) (entity.AreaOfConcernGroup, error) {
	if tx == nil {
		tx = r.db
	}

	for _, preload := range preloads {
		tx = tx.Preload(preload)
	}

	if err := tx.WithContext(ctx).Create(&areaOfConcernGroup).Error; err != nil {
		return entity.AreaOfConcernGroup{}, err
	}

	return areaOfConcernGroup, nil
}

func (r *areaOfConcernGroupRepository) GetAll(ctx context.Context, tx *gorm.DB, packageId string, metaReq meta.Meta, preloads ...string) ([]entity.AreaOfConcernGroup, meta.Meta, error) {
	if tx == nil {
		tx = r.db
	}

	for _, preload := range preloads {
		tx = tx.Preload(preload)
	}

	var areaOfConcernGroups []entity.AreaOfConcernGroup

	tx = tx.WithContext(ctx).Model(&entity.AreaOfConcernGroup{})

	if packageId != "" {
		tx = tx.Where("package_id = ?", packageId)
	}

	if err := WithFilters(tx, &metaReq, AddModels(entity.AreaOfConcernGroup{})).Find(&areaOfConcernGroups).Error; err != nil {
		return nil, meta.Meta{}, err
	}

	return areaOfConcernGroups, metaReq, nil
}

func (r *areaOfConcernGroupRepository) GetByID(ctx context.Context, tx *gorm.DB, areaOfConcernGroupID string, preloads ...string) (entity.AreaOfConcernGroup, error) {
	if tx == nil {
		tx = r.db
	}

	for _, preload := range preloads {
		tx = tx.Preload(preload)
	}

	var areaOfConcernGroup entity.AreaOfConcernGroup
	if err := tx.WithContext(ctx).First(&areaOfConcernGroup, "id = ?", areaOfConcernGroupID).Error; err != nil {
		return entity.AreaOfConcernGroup{}, err
	}

	return areaOfConcernGroup, nil
}

func (r *areaOfConcernGroupRepository) Update(ctx context.Context, tx *gorm.DB, areaOfConcernGroup entity.AreaOfConcernGroup, preloads ...string) error {
	if tx == nil {
		tx = r.db
	}

	for _, preload := range preloads {
		tx = tx.Preload(preload)
	}

	if err := tx.WithContext(ctx).Save(&areaOfConcernGroup).Error; err != nil {
		return err
	}

	return nil
}

func (r *areaOfConcernGroupRepository) Delete(ctx context.Context, tx *gorm.DB, areaOfConcernGroup entity.AreaOfConcernGroup, preloads ...string) error {
	if tx == nil {
		tx = r.db
	}

	for _, preload := range preloads {
		tx = tx.Preload(preload)
	}

	if err := tx.WithContext(ctx).Delete(&areaOfConcernGroup).Error; err != nil {
		return err
	}

	return nil
}
