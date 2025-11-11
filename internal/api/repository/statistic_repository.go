package repository

import (
	"context"

	"github.com/CRS-Project/crs-backend/internal/dto"
	"gorm.io/gorm"
)

type (
	StatisticRepository interface {
		GetCommentCard(ctx context.Context, tx *gorm.DB, packageId string) (dto.StatisticAOCAndCommentCard, error)
	}

	statisticRepository struct {
		db *gorm.DB
	}
)

func NewStatistic(db *gorm.DB) StatisticRepository {
	return &statisticRepository{
		db: db,
	}
}

func (r *statisticRepository) GetCommentCard(ctx context.Context, tx *gorm.DB, packageId string) (dto.StatisticAOCAndCommentCard, error) {
	if tx == nil {
		tx = r.db
	}

	var stats dto.StatisticAOCAndCommentCard
	err := tx.Raw(`
		SELECT
		(SELECT COUNT(*) FROM area_of_concerns a WHERE a.package_id = ?) AS total_area_of_concern,
		(SELECT COUNT(*) FROM documents d WHERE d.package_id = ?) AS total_documents,
		(SELECT COUNT(*) FROM comments c
			JOIN area_of_concerns a ON a.id = c.area_of_concern_id
			WHERE a.package_id = ? AND c.comment_reply_id IS NULL) AS total_comments,
		(SELECT COUNT(*) FROM comments c
			JOIN area_of_concerns a ON a.id = c.area_of_concern_id
			WHERE a.package_id = ? AND c.status = 'REJECTED') AS total_comment_rejected;
	`, packageId, packageId, packageId, packageId).Scan(&stats).Error

	if err != nil {
		return dto.StatisticAOCAndCommentCard{}, err
	}

	return stats, nil
}
