package repository

import (
	"context"

	"github.com/CRS-Project/crs-backend/internal/entity"
	"github.com/CRS-Project/crs-backend/internal/pkg/meta"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type (
	DisciplineListDocumentRepository interface {
		Create(ctx context.Context, tx *gorm.DB, disciplineListDocument entity.DisciplineListDocument, preloads ...string) (entity.DisciplineListDocument, error)
		GetAll(ctx context.Context, tx *gorm.DB, metaReq meta.Meta, preloads ...string) ([]entity.DisciplineListDocument, meta.Meta, error)
		GetAllByDisciplineGroupID(ctx context.Context, tx *gorm.DB, disciplineGroupId string, metaReq meta.Meta, preloads ...string) ([]entity.DisciplineListDocument, meta.Meta, error)
		GetByID(ctx context.Context, tx *gorm.DB, disciplineListDocumentID string, preloads ...string) (entity.DisciplineListDocument, error)
		Update(ctx context.Context, tx *gorm.DB, disciplineListDocument entity.DisciplineListDocument, preloads ...string) error
		Delete(ctx context.Context, tx *gorm.DB, disciplineListDocument entity.DisciplineListDocument, preloads ...string) error
		DeleteByDisciplineGroupID(ctx context.Context, tx *gorm.DB, disciplineGroupID string, preloads ...string) error
	}

	disciplineListDocumentRepository struct {
		db *gorm.DB
	}
)

func NewDisciplineListDocument(db *gorm.DB) DisciplineListDocumentRepository {
	return &disciplineListDocumentRepository{
		db: db,
	}
}

func (r *disciplineListDocumentRepository) Create(ctx context.Context, tx *gorm.DB, disciplineListDocument entity.DisciplineListDocument, preloads ...string) (entity.DisciplineListDocument, error) {
	if tx == nil {
		tx = r.db
	}

	for _, preload := range preloads {
		tx = tx.Preload(preload)
	}

	if err := tx.WithContext(ctx).Create(&disciplineListDocument).Error; err != nil {
		return entity.DisciplineListDocument{}, err
	}

	return disciplineListDocument, nil
}

func (r *disciplineListDocumentRepository) GetAll(ctx context.Context, tx *gorm.DB, metaReq meta.Meta, preloads ...string) ([]entity.DisciplineListDocument, meta.Meta, error) {
	if tx == nil {
		tx = r.db
	}

	for _, preload := range preloads {
		tx = tx.Preload(preload)
	}

	var disciplineListDocuments []entity.DisciplineListDocument

	tx = tx.WithContext(ctx).Model(&entity.DisciplineListDocument{})

	if err := WithFilters(tx, &metaReq, AddModels(entity.DisciplineListDocument{})).Find(&disciplineListDocuments).Error; err != nil {
		return nil, meta.Meta{}, err
	}

	return disciplineListDocuments, metaReq, nil
}

func (r *disciplineListDocumentRepository) GetAllByDisciplineGroupID(ctx context.Context, tx *gorm.DB, disciplineGroupId string, metaReq meta.Meta, preloads ...string) ([]entity.DisciplineListDocument, meta.Meta, error) {
	if tx == nil {
		tx = r.db
	}

	for _, preload := range preloads {
		tx = tx.Preload(preload)
	}

	var disciplineListDocuments []entity.DisciplineListDocument

	tx = tx.WithContext(ctx).Model(&entity.DisciplineListDocument{}).Where("discipline_group_id = ?", disciplineGroupId)

	filterMap := metaReq.SeparateFilter()

	// Check if we need to join documents
	// Join if search is present OR if sorting by due_date
	_, hasSearch := filterMap["search"]
	if hasSearch || metaReq.SortBy == "due_date" {
		tx = tx.Joins("LEFT JOIN documents d ON d.id = discipline_list_documents.document_id")
	}

	if find, ok := filterMap["search"]; ok {
		tx = tx.Where("d.company_document_number ILIKE ? OR d.document_serial_number ILIKE ?",
			"%"+find+"%",
			"%"+find+"%")
	}
	if err := WithFilters(tx, &metaReq, AddModels(entity.DisciplineListDocument{}),
		AddCustomField("search", ""),
		AddCustomField("due_date", "", "d.due_date")).Find(&disciplineListDocuments).Error; err != nil {
		return nil, meta.Meta{}, err
	}

	return disciplineListDocuments, metaReq, nil
}

func (r *disciplineListDocumentRepository) GetByID(ctx context.Context, tx *gorm.DB, disciplineListDocumentID string, preloads ...string) (entity.DisciplineListDocument, error) {
	if tx == nil {
		tx = r.db
	}

	for _, preload := range preloads {
		tx = tx.Preload(preload)
	}

	var disciplineListDocument entity.DisciplineListDocument
	if err := tx.WithContext(ctx).First(&disciplineListDocument, "id = ?", disciplineListDocumentID).Error; err != nil {
		return entity.DisciplineListDocument{}, err
	}

	return disciplineListDocument, nil
}

func (r *disciplineListDocumentRepository) Update(ctx context.Context, tx *gorm.DB, disciplineListDocument entity.DisciplineListDocument, preloads ...string) error {
	if tx == nil {
		tx = r.db
	}

	for _, preload := range preloads {
		tx = tx.Preload(preload)
	}

	if err := tx.WithContext(ctx).Save(&disciplineListDocument).Error; err != nil {
		return err
	}

	return nil
}

func (r *disciplineListDocumentRepository) Delete(ctx context.Context, tx *gorm.DB, disciplineListDocument entity.DisciplineListDocument, preloads ...string) error {
	if tx == nil {
		tx = r.db
	}

	for _, preload := range preloads {
		tx = tx.Preload(preload)
	}

	// persist deleted_by if provided
	if disciplineListDocument.DeletedBy != uuid.Nil {
		if err := tx.WithContext(ctx).Model(&entity.DisciplineListDocument{}).
			Where("id = ?", disciplineListDocument.ID).
			Updates(map[string]interface{}{"deleted_by": disciplineListDocument.DeletedBy}).Error; err != nil {
			return err
		}
	}

	if err := tx.WithContext(ctx).Delete(&disciplineListDocument).Error; err != nil {
		return err
	}

	return nil
}

func (r *disciplineListDocumentRepository) DeleteByDisciplineGroupID(ctx context.Context, tx *gorm.DB, disciplineGroupID string, preloads ...string) error {
	if tx == nil {
		tx = r.db
	}

	for _, preload := range preloads {
		tx = tx.Preload(preload)
	}

	if err := tx.WithContext(ctx).Where("discipline_group_id = ?", disciplineGroupID).Delete(&entity.DisciplineListDocument{}).Error; err != nil {
		return err
	}

	return nil
}
