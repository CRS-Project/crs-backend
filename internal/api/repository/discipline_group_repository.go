package repository

import (
	"context"

	"github.com/CRS-Project/crs-backend/internal/dto"
	"github.com/CRS-Project/crs-backend/internal/entity"
	"github.com/CRS-Project/crs-backend/internal/pkg/meta"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type (
	DisciplineGroupRepository interface {
		Create(ctx context.Context, tx *gorm.DB, disciplineGroup entity.DisciplineGroup, preloads ...string) (entity.DisciplineGroup, error)
		GetAll(ctx context.Context, tx *gorm.DB, packageId string, metaReq meta.Meta, preloads ...string) ([]entity.DisciplineGroup, meta.Meta, error)
		GetByID(ctx context.Context, tx *gorm.DB, disciplineGroupID string, preloads ...string) (entity.DisciplineGroup, error)
		Update(ctx context.Context, tx *gorm.DB, disciplineGroup entity.DisciplineGroup, preloads ...string) error
		Delete(ctx context.Context, tx *gorm.DB, disciplineGroup entity.DisciplineGroup, preloads ...string) error
		Statistic(ctx context.Context, tx *gorm.DB, packageId string) (dto.DisciplineGroupStatistic, error)
	}

	disciplineGroupRepository struct {
		db *gorm.DB
	}
)

func NewDisciplineGroup(db *gorm.DB) DisciplineGroupRepository {
	return &disciplineGroupRepository{
		db: db,
	}
}

func (r *disciplineGroupRepository) Create(ctx context.Context, tx *gorm.DB, disciplineGroup entity.DisciplineGroup, preloads ...string) (entity.DisciplineGroup, error) {
	if tx == nil {
		tx = r.db
	}

	for _, preload := range preloads {
		tx = tx.Preload(preload)
	}

	if err := tx.WithContext(ctx).Create(&disciplineGroup).Error; err != nil {
		return entity.DisciplineGroup{}, err
	}

	return disciplineGroup, nil
}

func (r *disciplineGroupRepository) GetAll(ctx context.Context, tx *gorm.DB, packageId string, metaReq meta.Meta, preloads ...string) ([]entity.DisciplineGroup, meta.Meta, error) {
	if tx == nil {
		tx = r.db
	}

	for _, preload := range preloads {
		tx = tx.Preload(preload)
	}

	var disciplineGroups []entity.DisciplineGroup

	tx = tx.WithContext(ctx).Model(&entity.DisciplineGroup{})
	if packageId != "" {
		tx = tx.Where("package_id = ?", packageId)
	}

	filterMap := metaReq.SeparateFilter()
	if find, ok := filterMap["search"]; ok {
		tx = tx.Where("discipline_groups.review_focus ILIKE ? OR discipline_groups.user_discipline ILIKE ?",
			"%"+find+"%",
			"%"+find+"%")
	}

	if err := WithFilters(tx, &metaReq, AddModels(entity.DisciplineGroup{}),
		AddCustomField("search", "")).Find(&disciplineGroups).Error; err != nil {
		return nil, meta.Meta{}, err
	}

	return disciplineGroups, metaReq, nil
}

func (r *disciplineGroupRepository) GetByID(ctx context.Context, tx *gorm.DB, disciplineGroupID string, preloads ...string) (entity.DisciplineGroup, error) {
	if tx == nil {
		tx = r.db
	}

	for _, preload := range preloads {
		tx = tx.Preload(preload)
	}

	var disciplineGroup entity.DisciplineGroup
	if err := tx.WithContext(ctx).First(&disciplineGroup, "id = ?", disciplineGroupID).Error; err != nil {
		return entity.DisciplineGroup{}, err
	}

	return disciplineGroup, nil
}

func (r *disciplineGroupRepository) Update(ctx context.Context, tx *gorm.DB, disciplineGroup entity.DisciplineGroup, preloads ...string) error {
	if tx == nil {
		tx = r.db
	}

	for _, preload := range preloads {
		tx = tx.Preload(preload)
	}

	if err := tx.WithContext(ctx).Save(&disciplineGroup).Error; err != nil {
		return err
	}

	return nil
}

func (r *disciplineGroupRepository) Delete(ctx context.Context, tx *gorm.DB, disciplineGroup entity.DisciplineGroup, preloads ...string) error {
	if tx == nil {
		tx = r.db
	}

	for _, preload := range preloads {
		tx = tx.Preload(preload)
	}

	// Persist DeletedBy (if set) before performing soft-delete so we can keep track who deleted the record
	if disciplineGroup.DeletedBy != uuid.Nil {
		if err := tx.WithContext(ctx).Model(&entity.DisciplineGroup{}).
			Where("id = ?", disciplineGroup.ID).
			Updates(map[string]interface{}{"deleted_by": disciplineGroup.DeletedBy}).Error; err != nil {
			return err
		}
	}

	if err := tx.WithContext(ctx).Delete(&disciplineGroup).Error; err != nil {
		return err
	}

	return nil
}

func (r *disciplineGroupRepository) Statistic(ctx context.Context, tx *gorm.DB, packageId string) (dto.DisciplineGroupStatistic, error) {
	if tx == nil {
		tx = r.db
	}

	var stats dto.DisciplineGroupStatistic
	err := tx.Raw(`
		SELECT
		(SELECT COUNT(*) FROM discipline_groups ag WHERE ag.package_id = ? AND deleted_at is null) AS total_discipline_group,
		(SELECT COUNT(*) FROM discipline_list_documents a WHERE a.package_id = ? AND deleted_at is null) AS total_discipline_list_document,
		(SELECT COUNT(*) FROM comments c
			JOIN discipline_list_documents a ON a.id = c.discipline_list_document_id
			WHERE a.package_id = ? AND c.comment_reply_id IS NULL AND c.deleted_at is null) AS total_comment;
	`, packageId, packageId, packageId).Scan(&stats).Error

	if err != nil {
		return dto.DisciplineGroupStatistic{}, err
	}

	return stats, nil
}
