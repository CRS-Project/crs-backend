package repository

import (
	"context"

	"github.com/CRS-Project/crs-backend/internal/entity"
	"gorm.io/gorm"
)

type (
	UserDisciplineNumberRepository interface {
		CountByUserDisciplineID(ctx context.Context, tx *gorm.DB, userDisciplineId string) (int, error)
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
