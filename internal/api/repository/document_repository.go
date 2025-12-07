package repository

import (
	"context"

	"github.com/CRS-Project/crs-backend/internal/entity"
	"github.com/CRS-Project/crs-backend/internal/pkg/meta"
	"gorm.io/gorm"
)

type (
	DocumentRepository interface {
		GetByID(ctx context.Context, tx *gorm.DB, documentID string, preloads ...string) (entity.Document, error)
		Create(ctx context.Context, tx *gorm.DB, document entity.Document, preloads ...string) (entity.Document, error)
		GetAll(ctx context.Context, tx *gorm.DB, packageId string, metaReq meta.Meta, preloads ...string) ([]entity.Document, meta.Meta, error)
		Delete(ctx context.Context, tx *gorm.DB, document entity.Document, preloads ...string) error
		Update(ctx context.Context, tx *gorm.DB, document entity.Document, preloads ...string) (entity.Document, error)
	}

	documentRepository struct {
		db *gorm.DB
	}
)

func NewDocument(db *gorm.DB) DocumentRepository {
	return &documentRepository{
		db: db,
	}
}

func (r *documentRepository) Create(ctx context.Context, tx *gorm.DB, document entity.Document, preloads ...string) (entity.Document, error) {
	if tx == nil {
		tx = r.db
	}

	for _, preload := range preloads {
		tx = tx.Preload(preload)
	}

	if err := tx.WithContext(ctx).Create(&document).Error; err != nil {
		return entity.Document{}, err
	}

	if len(preloads) > 0 {
		if err := tx.First(&document, "id = ?", document.ID).Error; err != nil {
			return entity.Document{}, err
		}
	}

	return document, nil
}

func (r *documentRepository) GetByID(ctx context.Context, tx *gorm.DB, documentID string, preloads ...string) (entity.Document, error) {
	if tx == nil {
		tx = r.db
	}

	for _, preload := range preloads {
		tx = tx.Preload(preload)
	}

	var document entity.Document
	if err := tx.WithContext(ctx).First(&document, "id = ?", documentID).Error; err != nil {
		return entity.Document{}, err
	}

	return document, nil
}

func (r *documentRepository) GetAll(ctx context.Context, tx *gorm.DB, packageId string, metaReq meta.Meta, preloads ...string) ([]entity.Document, meta.Meta, error) {
	if tx == nil {
		tx = r.db
	}

	for _, preload := range preloads {
		tx = tx.Preload(preload)
	}

	var documents []entity.Document

	tx = tx.WithContext(ctx).Model(&entity.Document{}).
		Joins("LEFT JOIN packages ON packages.id = documents.package_id")

	if packageId != "" {
		tx = tx.Where("package_id = ?", packageId)
	}

	filterMap := metaReq.SeparateFilter()
	if find, ok := filterMap["search"]; ok {
		tx = tx.Where("document_title ILIKE ? OR company_document_number ILIKE ? OR document_type ILIKE ? OR status ILIKE ? OR packages.name ILIKE ?",
			"%"+find+"%",
			"%"+find+"%",
			"%"+find+"%",
			"%"+find+"%",
			"%"+find+"%")
	}

	if err := WithFilters(tx, &metaReq, AddModels(entity.Document{}),
		AddCustomField("search", "")).
		Find(&documents).Error; err != nil {
		return nil, meta.Meta{}, err
	}

	return documents, metaReq, nil
}

func (r *documentRepository) Update(ctx context.Context, tx *gorm.DB, document entity.Document, preloads ...string) (entity.Document, error) {
	if tx == nil {
		tx = r.db
	}

	for _, preload := range preloads {
		tx = tx.Preload(preload)
	}

	if err := tx.WithContext(ctx).Save(&document).Error; err != nil {
		return entity.Document{}, err
	}

	return document, nil
}

func (r *documentRepository) Delete(ctx context.Context, tx *gorm.DB, document entity.Document, preloads ...string) error {
	if tx == nil {
		tx = r.db
	}

	for _, preload := range preloads {
		tx = tx.Preload(preload)
	}

	return tx.WithContext(ctx).Delete(&document).Error
}
