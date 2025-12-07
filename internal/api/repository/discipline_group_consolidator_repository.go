package repository

import (
	"context"

	"github.com/CRS-Project/crs-backend/internal/entity"
	"github.com/CRS-Project/crs-backend/internal/pkg/meta"
	"gorm.io/gorm"
)

type (
	DisciplineGroupConsolidatorRepository interface {
		Create(ctx context.Context, tx *gorm.DB, disciplineGroupConsolidator entity.DisciplineGroupConsolidator, preloads ...string) (entity.DisciplineGroupConsolidator, error)
		CreateBulk(ctx context.Context, tx *gorm.DB, disciplineGroupConsolidator []entity.DisciplineGroupConsolidator, preloads ...string) error
		GetAll(ctx context.Context, tx *gorm.DB, metaReq meta.Meta, preloads ...string) ([]entity.DisciplineGroupConsolidator, meta.Meta, error)
		GetAllConsolidator(ctx context.Context, tx *gorm.DB, search, disciplineGroupId string, preloads ...string) ([]entity.DisciplineGroupConsolidator, error)
		GetByID(ctx context.Context, tx *gorm.DB, disciplineGroupConsolidatorID string, preloads ...string) (entity.DisciplineGroupConsolidator, error)
		GetByUserID(ctx context.Context, tx *gorm.DB, userID string, preloads ...string) ([]entity.DisciplineGroupConsolidator, error)
		Update(ctx context.Context, tx *gorm.DB, disciplineGroupConsolidator entity.DisciplineGroupConsolidator, preloads ...string) error
		Delete(ctx context.Context, tx *gorm.DB, disciplineGroupConsolidator entity.DisciplineGroupConsolidator, preloads ...string) error
		DeleteByID(ctx context.Context, tx *gorm.DB, id string) error
		DeleteByUserID(ctx context.Context, tx *gorm.DB, userID string) error
		DeleteByDisciplineGroupID(ctx context.Context, tx *gorm.DB, disciplineGroupID string) error
		DeleteBulk(ctx context.Context, tx *gorm.DB, disciplineGroupConcolidatorIDs []string) error
	}

	disciplineGroupConsolidatorRepository struct {
		db *gorm.DB
	}
)

func NewDisciplineGroupConsolidator(db *gorm.DB) DisciplineGroupConsolidatorRepository {
	return &disciplineGroupConsolidatorRepository{
		db: db,
	}
}

func (r *disciplineGroupConsolidatorRepository) Create(ctx context.Context, tx *gorm.DB, disciplineGroupConsolidator entity.DisciplineGroupConsolidator, preloads ...string) (entity.DisciplineGroupConsolidator, error) {
	if tx == nil {
		tx = r.db
	}

	for _, preload := range preloads {
		tx = tx.Preload(preload)
	}

	if err := tx.WithContext(ctx).Create(&disciplineGroupConsolidator).Error; err != nil {
		return entity.DisciplineGroupConsolidator{}, err
	}

	return disciplineGroupConsolidator, nil
}

func (r *disciplineGroupConsolidatorRepository) CreateBulk(ctx context.Context, tx *gorm.DB, disciplineGroupConsolidators []entity.DisciplineGroupConsolidator, preloads ...string) error {
	if tx == nil {
		tx = r.db
	}

	for _, preload := range preloads {
		tx = tx.Preload(preload)
	}

	if err := tx.WithContext(ctx).Create(&disciplineGroupConsolidators).Error; err != nil {
		return err
	}

	return nil
}

func (r *disciplineGroupConsolidatorRepository) GetAll(ctx context.Context, tx *gorm.DB, metaReq meta.Meta, preloads ...string) ([]entity.DisciplineGroupConsolidator, meta.Meta, error) {
	if tx == nil {
		tx = r.db
	}

	for _, preload := range preloads {
		tx = tx.Preload(preload)
	}

	var disciplineGroupConsolidators []entity.DisciplineGroupConsolidator

	tx = tx.WithContext(ctx).Model(&entity.DisciplineGroupConsolidator{})
	if err := WithFilters(tx, &metaReq, AddModels(entity.DisciplineGroupConsolidator{})).Find(&disciplineGroupConsolidators).Error; err != nil {
		return nil, meta.Meta{}, err
	}

	return disciplineGroupConsolidators, metaReq, nil
}

func (r *disciplineGroupConsolidatorRepository) GetAllConsolidator(ctx context.Context, tx *gorm.DB, search, disciplineGroupId string, preloads ...string) ([]entity.DisciplineGroupConsolidator, error) {
	if tx == nil {
		tx = r.db
	}

	for _, preload := range preloads {
		tx = tx.Preload(preload)
	}

	var disciplineGroupConsolidators []entity.DisciplineGroupConsolidator

	tx = tx.WithContext(ctx).Model(&entity.DisciplineGroupConsolidator{}).
		Joins("LEFT JOIN users u ON u.id = discipline_group_consolidators.user_id").
		Where("discipline_group_consolidators.discipline_group_id = ?", disciplineGroupId)

	if search != "" {
		tx.Where("u.name ILIKE ?", search)
	}

	if err := tx.Find(&disciplineGroupConsolidators).Error; err != nil {
		return nil, err
	}

	return disciplineGroupConsolidators, nil
}

func (r *disciplineGroupConsolidatorRepository) GetByID(ctx context.Context, tx *gorm.DB, disciplineGroupConsolidatorID string, preloads ...string) (entity.DisciplineGroupConsolidator, error) {
	if tx == nil {
		tx = r.db
	}

	for _, preload := range preloads {
		tx = tx.Preload(preload)
	}

	var disciplineGroupConsolidator entity.DisciplineGroupConsolidator
	if err := tx.WithContext(ctx).First(&disciplineGroupConsolidator, "id = ?", disciplineGroupConsolidatorID).Error; err != nil {
		return entity.DisciplineGroupConsolidator{}, err
	}

	return disciplineGroupConsolidator, nil
}

func (r *disciplineGroupConsolidatorRepository) GetByUserID(ctx context.Context, tx *gorm.DB, userID string, preloads ...string) ([]entity.DisciplineGroupConsolidator, error) {
	if tx == nil {
		tx = r.db
	}

	for _, preload := range preloads {
		tx = tx.Preload(preload)
	}

	var disciplineGroupConsolidator []entity.DisciplineGroupConsolidator
	if err := tx.WithContext(ctx).Find(&disciplineGroupConsolidator, "user_id = ?", userID).Error; err != nil {
		return nil, err
	}

	return disciplineGroupConsolidator, nil
}

func (r *disciplineGroupConsolidatorRepository) Update(ctx context.Context, tx *gorm.DB, disciplineGroupConsolidator entity.DisciplineGroupConsolidator, preloads ...string) error {
	if tx == nil {
		tx = r.db
	}

	for _, preload := range preloads {
		tx = tx.Preload(preload)
	}

	if err := tx.WithContext(ctx).Save(&disciplineGroupConsolidator).Error; err != nil {
		return err
	}

	return nil
}

func (r *disciplineGroupConsolidatorRepository) Delete(ctx context.Context, tx *gorm.DB, disciplineGroupConsolidator entity.DisciplineGroupConsolidator, preloads ...string) error {
	if tx == nil {
		tx = r.db
	}

	for _, preload := range preloads {
		tx = tx.Preload(preload)
	}

	if err := tx.WithContext(ctx).Delete(&disciplineGroupConsolidator).Error; err != nil {
		return err
	}

	return nil
}

func (r *disciplineGroupConsolidatorRepository) DeleteByID(ctx context.Context, tx *gorm.DB, id string) error {
	if tx == nil {
		tx = r.db
	}

	if err := tx.WithContext(ctx).
		Where("id = ?", id).
		Delete(&entity.DisciplineGroupConsolidator{}).Error; err != nil {
		return err
	}

	return nil
}

func (r *disciplineGroupConsolidatorRepository) DeleteByUserID(ctx context.Context, tx *gorm.DB, userID string) error {
	if tx == nil {
		tx = r.db
	}

	if err := tx.WithContext(ctx).
		Where("user_id = ?", userID).
		Delete(&entity.DisciplineGroupConsolidator{}).Error; err != nil {
		return err
	}

	return nil
}

func (r *disciplineGroupConsolidatorRepository) DeleteByDisciplineGroupID(ctx context.Context, tx *gorm.DB, disciplineGroupID string) error {
	if tx == nil {
		tx = r.db
	}

	if err := tx.WithContext(ctx).
		Where("discipline_group_id = ?", disciplineGroupID).
		Delete(&entity.DisciplineGroupConsolidator{}).Error; err != nil {
		return err
	}

	return nil
}

func (r *disciplineGroupConsolidatorRepository) DeleteBulk(ctx context.Context, tx *gorm.DB, disciplineGroupConcolidatorIDs []string) error {
	if tx == nil {
		tx = r.db
	}

	if err := tx.WithContext(ctx).
		Where("id IN ?", disciplineGroupConcolidatorIDs).
		Delete(&entity.DisciplineGroupConsolidator{}).Error; err != nil {
		return err
	}

	return nil
}
