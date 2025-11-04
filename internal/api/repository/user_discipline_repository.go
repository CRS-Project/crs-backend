package repository

import (
	"context"

	"github.com/CRS-Project/crs-backend/internal/entity"
	"github.com/CRS-Project/crs-backend/internal/pkg/meta"
	"gorm.io/gorm"
)

type (
	UserDisciplineRepository interface {
		GetAll(ctx context.Context, tx *gorm.DB, metaReq meta.Meta, preloads ...string) ([]entity.UserDiscipline, error)
		GetByID(ctx context.Context, tx *gorm.DB, userDisciplineId string, preloads ...string) (entity.UserDiscipline, error)
		GetContractorDiscipline(ctx context.Context, tx *gorm.DB) (entity.UserDiscipline, error)
	}

	userDisciplineRepository struct {
		db *gorm.DB
	}
)

func NewUserDiscipline(db *gorm.DB) UserDisciplineRepository {
	return &userDisciplineRepository{db}
}

func (r *userDisciplineRepository) GetAll(ctx context.Context, tx *gorm.DB, metaReq meta.Meta, preloads ...string) ([]entity.UserDiscipline, error) {
	if tx == nil {
		tx = r.db
	}

	for _, preload := range preloads {
		tx = tx.Preload(preload)
	}

	var userDiscipline []entity.UserDiscipline

	tx = tx.WithContext(ctx).Model(entity.UserDiscipline{})
	if err := WithFilters(tx.Debug(), &metaReq, AddModels(entity.UserDiscipline{})).Find(&userDiscipline).Error; err != nil {
		return nil, err
	}

	return userDiscipline, nil
}

func (r *userDisciplineRepository) GetByID(ctx context.Context, tx *gorm.DB, userDisciplineId string, preloads ...string) (entity.UserDiscipline, error) {
	if tx == nil {
		tx = r.db
	}

	var userDiscipline entity.UserDiscipline

	if err := tx.WithContext(ctx).Where("id = ?", userDisciplineId).First(&userDiscipline).Error; err != nil {
		return entity.UserDiscipline{}, err
	}

	return userDiscipline, nil
}

func (r *userDisciplineRepository) GetContractorDiscipline(ctx context.Context, tx *gorm.DB) (entity.UserDiscipline, error) {
	if tx == nil {
		tx = r.db
	}

	var userDiscipline entity.UserDiscipline

	if err := tx.WithContext(ctx).Where("initial = ?", "CONTRACTOR").First(&userDiscipline).Error; err != nil {
		return entity.UserDiscipline{}, err
	}

	return userDiscipline, nil
}
