package repository

import (
	"context"

	"github.com/CRS-Project/crs-backend/internal/dto"
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

		Statistic(ctx context.Context, tx *gorm.DB, packageId string) (dto.AreaOfConcernGroupStatistic, error)
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

	tx = tx.WithContext(ctx).Model(&entity.AreaOfConcernGroup{}).
		Joins("LEFT JOIN user_disciplines ON user_disciplines.id = area_of_concern_groups.user_discipline_id")

	if packageId != "" {
		tx = tx.Where("package_id = ?", packageId)
	}

	filterMap := metaReq.SeparateFilter()
	if find, ok := filterMap["search"]; ok {
		tx = tx.Where("area_of_concern_groups.review_focus ILIKE ? OR user_disciplines.name ILIKE ?",
			"%"+find+"%",
			"%"+find+"%")
	}

	if err := WithFilters(tx, &metaReq, AddModels(entity.AreaOfConcernGroup{}),
		AddCustomField("search", ""),
		AddCustomField("user_discipline", "", "user_disciplines.name")).Find(&areaOfConcernGroups).Error; err != nil {
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

func (r *areaOfConcernGroupRepository) Statistic(ctx context.Context, tx *gorm.DB, packageId string) (dto.AreaOfConcernGroupStatistic, error) {
	if tx == nil {
		tx = r.db
	}

	var stats dto.AreaOfConcernGroupStatistic
	err := tx.Raw(`
		SELECT
		(SELECT COUNT(*) FROM area_of_concern_groups ag WHERE ag.package_id = ? AND deleted_at is null) AS total_area_of_concern_group,
		(SELECT COUNT(*) FROM area_of_concerns a WHERE a.package_id = ? AND deleted_at is null) AS total_area_of_concern,
		(SELECT COUNT(*) FROM comments c
			JOIN area_of_concerns a ON a.id = c.area_of_concern_id
			WHERE a.package_id = ? AND c.comment_reply_id IS NULL AND c.deleted_at is null) AS total_comment;
	`, packageId, packageId, packageId).Scan(&stats).Error

	if err != nil {
		return dto.AreaOfConcernGroupStatistic{}, err
	}

	return stats, nil
}
