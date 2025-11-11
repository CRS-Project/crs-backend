package repository

import (
	"context"

	"github.com/CRS-Project/crs-backend/internal/dto"
	"gorm.io/gorm"
)

type (
	StatisticRepository interface {
		GetAOCAndCommentChart(ctx context.Context, tx *gorm.DB, packageId string) ([]dto.StatisticAOCAndCommentChart, error)
		GetCommentCard(ctx context.Context, tx *gorm.DB, packageId string) (dto.StatisticAOCAndCommentCard, error)
		GetCommentUserChart(ctx context.Context, tx *gorm.DB, packageId string) ([]dto.StatisticCommentUsersChart, error)
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

func (r *statisticRepository) GetAOCAndCommentChart(ctx context.Context, tx *gorm.DB, packageId string) ([]dto.StatisticAOCAndCommentChart, error) {
	if tx == nil {
		tx = r.db
	}

	query := `
	WITH date_series AS (
		SELECT generate_series(
			LEAST(
				(SELECT MIN(created_at)::date FROM area_of_concerns WHERE package_id = ?),
				(SELECT MIN(created_at)::date FROM documents WHERE package_id = ?)
			),
			NOW()::date,
			'2 days'
		) AS start_date
	),
	aoc_count AS (
		SELECT
			DATE_TRUNC('day', created_at)::date AS created_date,
			COUNT(*) AS total_area_of_concern
		FROM area_of_concerns
		WHERE deleted_at IS NULL
		AND package_id = ?
		GROUP BY 1
	),
	doc_count AS (
		SELECT
			DATE_TRUNC('day', created_at)::date AS created_date,
			COUNT(*) AS total_documents
		FROM documents
		WHERE deleted_at IS NULL
		AND package_id = ?
		GROUP BY 1
	),
	comment_count AS (
		SELECT
			DATE_TRUNC('day', c.created_at)::date AS created_date,
			COUNT(*) FILTER (WHERE c.comment_reply_id IS NULL) AS total_comments,
			COUNT(*) FILTER (WHERE c.status = 'REJECT') AS total_comment_rejected
		FROM comments c
		JOIN area_of_concerns a ON a.id = c.area_of_concern_id
		WHERE c.deleted_at IS NULL
		AND a.deleted_at IS NULL
		AND a.package_id = ?
		GROUP BY 1
	),
	aoc_by_interval AS (
		SELECT
			ds.start_date,
			COALESCE(SUM(a.total_area_of_concern), 0) AS total_area_of_concern
		FROM date_series ds
		LEFT JOIN aoc_count a ON a.created_date >= ds.start_date 
			AND a.created_date < ds.start_date + INTERVAL '2 days'
		GROUP BY ds.start_date
	),
	doc_by_interval AS (
		SELECT
			ds.start_date,
			COALESCE(SUM(d.total_documents), 0) AS total_documents
		FROM date_series ds
		LEFT JOIN doc_count d ON d.created_date >= ds.start_date 
			AND d.created_date < ds.start_date + INTERVAL '2 days'
		GROUP BY ds.start_date
	),
	comment_by_interval AS (
		SELECT
			ds.start_date,
			COALESCE(SUM(c.total_comments), 0) AS total_comments,
			COALESCE(SUM(c.total_comment_rejected), 0) AS total_comment_rejected
		FROM date_series ds
		LEFT JOIN comment_count c ON c.created_date >= ds.start_date 
			AND c.created_date < ds.start_date + INTERVAL '2 days'
		GROUP BY ds.start_date
	)
	SELECT
		TO_CHAR(ds.start_date, 'DD-Mon') AS name,
		a.total_area_of_concern,
		d.total_documents,
		c.total_comments,
		c.total_comment_rejected
	FROM date_series ds
	LEFT JOIN aoc_by_interval a USING (start_date)
	LEFT JOIN doc_by_interval d USING (start_date)
	LEFT JOIN comment_by_interval c USING (start_date)
	ORDER BY ds.start_date;
	`

	var stats []dto.StatisticAOCAndCommentChart
	err := tx.Raw(query, packageId, packageId, packageId, packageId, packageId).Scan(&stats).Error

	if err != nil {
		return nil, err
	}

	return stats, nil
}

func (r *statisticRepository) GetCommentCard(ctx context.Context, tx *gorm.DB, packageId string) (dto.StatisticAOCAndCommentCard, error) {
	if tx == nil {
		tx = r.db
	}

	var stats dto.StatisticAOCAndCommentCard
	err := tx.Raw(`
		SELECT
		(SELECT COUNT(*) FROM area_of_concerns a WHERE a.package_id = ? AND deleted_at is null) AS total_area_of_concern,
		(SELECT COUNT(*) FROM documents d WHERE d.package_id = ? AND deleted_at is null) AS total_documents,
		(SELECT COUNT(*) FROM comments c
			JOIN area_of_concerns a ON a.id = c.area_of_concern_id
			WHERE a.package_id = ? AND c.comment_reply_id IS NULL AND c.deleted_at is null) AS total_comments,
		(SELECT COUNT(*) FROM comments c
			JOIN area_of_concerns a ON a.id = c.area_of_concern_id
			WHERE a.package_id = ? AND c.status = 'REJECT' AND c.deleted_at is null) AS total_comment_rejected;
	`, packageId, packageId, packageId, packageId).Scan(&stats).Error

	if err != nil {
		return dto.StatisticAOCAndCommentCard{}, err
	}

	return stats, nil
}

func (r *statisticRepository) GetCommentUserChart(ctx context.Context, tx *gorm.DB, packageId string) ([]dto.StatisticCommentUsersChart, error) {
	if tx == nil {
		tx = r.db
	}

	query := `
	SELECT
		u.id,
		u.initial AS name,
		COALESCE(COUNT(c.id) FILTER (WHERE c.status = 'ACCEPTED' OR c.status = 'REJECT'), 0) AS comment_closed,
		COALESCE(COUNT(c.id), 0) AS total_comment
	FROM users u
	LEFT JOIN comments c ON c.user_id = u.id
		AND c.deleted_at IS NULL
		AND c.comment_reply_id IS NULL
	WHERE u.deleted_at IS NULL
		AND u.package_id = ?
	GROUP BY u.id, u.initial
	ORDER BY u.initial;
	`

	var stats []dto.StatisticCommentUsersChart
	err := tx.Raw(query, packageId).Scan(&stats).Error

	if err != nil {
		return nil, err
	}

	return stats, nil
}
