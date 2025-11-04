package repository

import (
	"context"

	"github.com/CRS-Project/crs-backend/internal/entity"
	"gorm.io/gorm"
)

type (
	UserPackageRepository interface {
		Create(ctx context.Context, tx *gorm.DB, userPackage entity.UserPackage) (entity.UserPackage, error)
		Save(ctx context.Context, tx *gorm.DB, userPackage entity.UserPackage) (entity.UserPackage, error)
	}

	userPackageRepository struct {
		db *gorm.DB
	}
)

func NewUserPackage(db *gorm.DB) UserPackageRepository {
	return &userPackageRepository{db}
}

func (r *userPackageRepository) Create(ctx context.Context, tx *gorm.DB, userPackage entity.UserPackage) (entity.UserPackage, error) {
	if tx == nil {
		tx = r.db
	}

	tx = tx.WithContext(ctx)
	if err := tx.Create(&userPackage).Error; err != nil {
		return entity.UserPackage{}, err
	}

	return userPackage, nil
}

func (r *userPackageRepository) Save(ctx context.Context, tx *gorm.DB, userPackage entity.UserPackage) (entity.UserPackage, error) {
	if tx == nil {
		tx = r.db
	}

	tx = tx.WithContext(ctx)
	if err := tx.Save(&userPackage).Error; err != nil {
		return entity.UserPackage{}, err
	}

	return userPackage, nil
}
