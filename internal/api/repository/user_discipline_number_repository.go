package repository

import (
	"context"

	"github.com/CRS-Project/crs-backend/internal/entity"
	"gorm.io/gorm"
)

type (
	UserDisciplineNumberRepository interface {
		CountByUserDisciplineID(ctx context.Context, tx *gorm.DB, userDisciplineId string) (int, error)
		Update(ctx context.Context, tx *gorm.DB, userDisciplineNumber entity.UserDisciplineNumber) (entity.UserDisciplineNumber, error)
	}

	userDisciplineNumberRepository struct {
		db *gorm.DB
	}
)

func NewUserDisciplineNumber(db *gorm.DB) UserDisciplineNumberRepository {
	return &userDisciplineNumberRepository{db}
}

func (r *userDisciplineNumberRepository) CountByUserDisciplineID(ctx context.Context, tx *gorm.DB, userDisciplineId string) (int, error) {
	if tx == nil {
		tx = r.db
	}

	var count int64

	tx = tx.WithContext(ctx).Model(entity.UserDisciplineNumber{})
	if err := tx.Where("user_discipline_id = ?", userDisciplineId).Count(&count).Error; err != nil {
		return 0, err
	}

	return int(count), nil
}

func (r *userDisciplineNumberRepository) Update(ctx context.Context, tx *gorm.DB, userDisciplineNumber entity.UserDisciplineNumber) (entity.UserDisciplineNumber, error) {
	if tx == nil {
		tx = r.db
	}
	tx = tx.WithContext(ctx)

	if err := tx.Save(&userDisciplineNumber).Error; err != nil {
		return entity.UserDisciplineNumber{}, err
	}

	return userDisciplineNumber, nil
}
