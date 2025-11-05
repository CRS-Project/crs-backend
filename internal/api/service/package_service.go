package service

import (
	"context"
	"net/http"

	"github.com/CRS-Project/crs-backend/internal/api/repository"
	"github.com/CRS-Project/crs-backend/internal/dto"
	"github.com/CRS-Project/crs-backend/internal/entity"
	myerror "github.com/CRS-Project/crs-backend/internal/pkg/error"
	"github.com/CRS-Project/crs-backend/internal/pkg/meta"
	"gorm.io/gorm"
)

type (
	PackageService interface {
		CreatePackage(ctx context.Context, req dto.CreatePackageRequest) (dto.PackageInfo, error)
		GetAll(ctx context.Context, metaReq meta.Meta) ([]dto.PackageInfo, meta.Meta, error)
		UpdatePackage(ctx context.Context, req dto.UpdatePackageRequest) error
		DeletePackage(ctx context.Context, id string) error
		GetByID(ctx context.Context, id string) (dto.PackageInfo, error)
	}

	packageService struct {
		packageRepository repository.PackageRepository
		db                *gorm.DB
	}
)

func NewPackage(packageRepository repository.PackageRepository, db *gorm.DB) PackageService {
	return &packageService{
		packageRepository: packageRepository,
		db:                db,
	}
}

func (s *packageService) CreatePackage(ctx context.Context, req dto.CreatePackageRequest) (dto.PackageInfo, error) {
	_, err := s.packageRepository.GetByName(ctx, nil, req.Name)
	if err == nil {
		return dto.PackageInfo{}, myerror.New("package with this name already exists", http.StatusConflict)
	}

	pkgCreation := entity.Package{
		Name: req.Name,
	}

	pkgResult, err := s.packageRepository.Create(ctx, nil, pkgCreation)
	if err != nil {
		return dto.PackageInfo{}, err
	}

	return dto.PackageInfo{
		ID:   pkgResult.ID.String(),
		Name: pkgResult.Name,
	}, nil
}

func (s *packageService) GetAll(ctx context.Context, metaReq meta.Meta) ([]dto.PackageInfo, meta.Meta, error) {
	pkgs, metaRes, err := s.packageRepository.GetAll(ctx, nil, metaReq)
	if err != nil {
		return nil, meta.Meta{}, err
	}

	var pkgInfos []dto.PackageInfo
	for _, pkg := range pkgs {
		pkgInfos = append(pkgInfos, dto.PackageInfo{
			ID:   pkg.ID.String(),
			Name: pkg.Name,
		})
	}

	return pkgInfos, metaRes, nil
}

func (s *packageService) UpdatePackage(ctx context.Context, req dto.UpdatePackageRequest) error {
	pkg, err := s.packageRepository.GetByID(ctx, nil, req.ID)
	if err != nil {
		return err
	}
	pkg.Name = req.Name

	if err = s.packageRepository.Update(ctx, nil, pkg); err != nil {
		return err
	}

	return nil
}

func (s *packageService) DeletePackage(ctx context.Context, id string) error {
	pkg, err := s.packageRepository.GetByID(ctx, nil, id)
	if err != nil {
		return err
	}

	if err = s.packageRepository.Delete(ctx, nil, pkg); err != nil {
		return err
	}

	return nil
}

func (s *packageService) GetByID(ctx context.Context, id string) (dto.PackageInfo, error) {
	pkg, err := s.packageRepository.GetByID(ctx, nil, id)
	if err != nil {
		return dto.PackageInfo{}, err
	}

	return dto.PackageInfo{
		ID:   pkg.ID.String(),
		Name: pkg.Name,
	}, nil
}
