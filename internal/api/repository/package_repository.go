package repository

import (
	"context"

	"github.com/CRS-Project/crs-backend/internal/entity"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type (
	PackageRepository interface {
		GetByID(ctx context.Context, tx *gorm.DB, pkgID uuid.UUID, preload ...string) (entity.Package, error)
		GetByName(ctx context.Context, tx *gorm.DB, pkgName string, preload ...string) (entity.Package, error)
		Create(ctx context.Context, tx *gorm.DB, pkg entity.Package, preload ...string) (entity.Package, error)
		GetAll(ctx context.Context, tx *gorm.DB, preload ...string) ([]entity.Package, error)
		Update(ctx context.Context, tx *gorm.DB, pkg entity.Package, preload ...string) error
		Delete(ctx context.Context, tx *gorm.DB, pkg entity.Package, preload ...string) error
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

func (r *packageRepository) GetByID(ctx context.Context, tx *gorm.DB, pkgID uuid.UUID, preload ...string) (entity.Package, error) {
	if tx == nil {
		tx = r.db
	}

	var pkg entity.Package
	if err := tx.WithContext(ctx).First(&pkg, pkgID).Error; err != nil {
		return entity.Package{}, err
	}

	return pkg, nil
}

func (r *packageRepository) GetByName(ctx context.Context, tx *gorm.DB, pkgName string, preload ...string) (entity.Package, error) {
	if tx == nil {
		tx = r.db
	}

	var pkg entity.Package
	if err := tx.WithContext(ctx).Take(&pkg, "name = ?", pkgName).Error; err != nil {
		return entity.Package{}, err
	}

	return pkg, nil
}

func (r *packageRepository) Create(ctx context.Context, tx *gorm.DB, pkg entity.Package, preload ...string) (entity.Package, error) {
	if tx == nil {
		tx = r.db
	}

	if err := tx.WithContext(ctx).Create(&pkg).Error; err != nil {
		return pkg, err
	}

	return pkg, nil
}

func (r *packageRepository) GetAll(ctx context.Context, tx *gorm.DB, preload ...string) ([]entity.Package, error) {
	if tx == nil {
		tx = r.db
	}

	var pkgs []entity.Package
	if err := tx.WithContext(ctx).Find(&pkgs).Error; err != nil {
		return []entity.Package{}, err
	}

	return pkgs, nil
}

func (r *packageRepository) Update(ctx context.Context, tx *gorm.DB, pkg entity.Package, preload ...string) error {
	if tx == nil {
		tx = r.db
	}

	if err := tx.WithContext(ctx).Model(&pkg).Update("name", pkg.Name).Error; err != nil {
		return err
	}

	return nil
}

func (r *packageRepository) Delete(ctx context.Context, tx *gorm.DB, pkg entity.Package, preload ...string) error {
	if tx == nil {
		tx = r.db
	}

	if err := tx.WithContext(ctx).Delete(&pkg).Error; err != nil {
		return err
	}

	return nil
}
