package repository

import (
	"context"

	"github.com/CRS-Project/crs-backend/internal/entity"
	"github.com/CRS-Project/crs-backend/internal/pkg/meta"
	"gorm.io/gorm"
)

type (
	PackageRepository interface {
		GetByID(ctx context.Context, tx *gorm.DB, pkgID string, preloads ...string) (entity.Package, error)
		GetByName(ctx context.Context, tx *gorm.DB, pkgName string, preloads ...string) (entity.Package, error)
		Create(ctx context.Context, tx *gorm.DB, pkg entity.Package, preloads ...string) (entity.Package, error)
		GetAll(ctx context.Context, tx *gorm.DB, metaReq meta.Meta, preloads ...string) ([]entity.Package, meta.Meta, error)
		Update(ctx context.Context, tx *gorm.DB, pkg entity.Package, preloads ...string) (entity.Package, error)
		Delete(ctx context.Context, tx *gorm.DB, pkg entity.Package, preloads ...string) error
	}

	packageRepository struct {
		db *gorm.DB
	}
)

func NewPackage(db *gorm.DB) PackageRepository {
	return &packageRepository{
		db: db,
	}
}

func (r *packageRepository) GetByID(ctx context.Context, tx *gorm.DB, pkgID string, preloads ...string) (entity.Package, error) {
	if tx == nil {
		tx = r.db
	}

	for _, preload := range preloads {
		tx = tx.Preload(preload)
	}

	var pkg entity.Package
	if err := tx.WithContext(ctx).First(&pkg, "id = ?", pkgID).Error; err != nil {
		return entity.Package{}, err
	}

	return pkg, nil
}

func (r *packageRepository) GetByName(ctx context.Context, tx *gorm.DB, pkgName string, preloads ...string) (entity.Package, error) {
	if tx == nil {
		tx = r.db
	}

	for _, preload := range preloads {
		tx = tx.Preload(preload)
	}

	var pkg entity.Package
	if err := tx.WithContext(ctx).Take(&pkg, "name = ?", pkgName).Error; err != nil {
		return entity.Package{}, err
	}

	return pkg, nil
}

func (r *packageRepository) Create(ctx context.Context, tx *gorm.DB, pkg entity.Package, preloads ...string) (entity.Package, error) {
	if tx == nil {
		tx = r.db
	}

	for _, preload := range preloads {
		tx = tx.Preload(preload)
	}

	if err := tx.WithContext(ctx).Create(&pkg).Error; err != nil {
		return entity.Package{}, err
	}

	return pkg, nil
}

func (r *packageRepository) GetAll(ctx context.Context, tx *gorm.DB, metaReq meta.Meta, preloads ...string) ([]entity.Package, meta.Meta, error) {
	if tx == nil {
		tx = r.db
	}

	for _, preload := range preloads {
		tx = tx.Preload(preload)
	}

	var pkgs []entity.Package

	tx = tx.WithContext(ctx).Model(&entity.Package{})

	if err := WithFilters(tx, &metaReq, AddModels(entity.Package{})).Find(&pkgs).Error; err != nil {
		return nil, meta.Meta{}, err
	}

	return pkgs, metaReq, nil
}

func (r *packageRepository) Update(ctx context.Context, tx *gorm.DB, pkg entity.Package, preloads ...string) (entity.Package, error) {
	if tx == nil {
		tx = r.db
	}

	for _, preload := range preloads {
		tx = tx.Preload(preload)
	}

	if err := tx.WithContext(ctx).Save(&pkg).Error; err != nil {
		return entity.Package{}, err
	}

	return pkg, nil
}

func (r *packageRepository) Delete(ctx context.Context, tx *gorm.DB, pkg entity.Package, preloads ...string) error {
	if tx == nil {
		tx = r.db
	}

	for _, preload := range preloads {
		tx = tx.Preload(preload)
	}

	if err := tx.WithContext(ctx).Delete(&pkg).Error; err != nil {
		return err
	}

	return nil
}
