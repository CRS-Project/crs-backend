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
		GetAllByUser(ctx context.Context, userId string) ([]dto.PackageInfo, error)
		UpdatePackage(ctx context.Context, req dto.UpdatePackageRequest) (dto.PackageInfo, error)
		DeletePackage(ctx context.Context, id string) error
		GetByID(ctx context.Context, id string) (dto.PackageInfo, error)
	}

	packageService struct {
		packageRepository repository.PackageRepository
		userRepository    repository.UserRepository
		db                *gorm.DB
	}
)

func NewPackage(packageRepository repository.PackageRepository,
	userRepository repository.UserRepository,
	db *gorm.DB) PackageService {
	return &packageService{
		packageRepository: packageRepository,
		userRepository:    userRepository,
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

	return pkgResult.ToInfo(), nil
}

func (s *packageService) GetAll(ctx context.Context, metaReq meta.Meta) ([]dto.PackageInfo, meta.Meta, error) {
	pkgs, metaRes, err := s.packageRepository.GetAll(ctx, nil, metaReq)
	if err != nil {
		return nil, meta.Meta{}, err
	}

	var pkgInfos []dto.PackageInfo
	for _, pkg := range pkgs {
		pkgInfos = append(pkgInfos, pkg.ToInfo())
	}

	return pkgInfos, metaRes, nil
}

func (s *packageService) GetAllByUser(ctx context.Context, userId string) ([]dto.PackageInfo, error) {
	user, err := s.userRepository.GetById(ctx, nil, userId)
	if err != nil {
		return nil, err
	}

	pkgs, err := s.packageRepository.GetAllNoPag(ctx, nil)
	if err != nil {
		return nil, err
	}

	if user.PackageID != nil {
		pkg, err := s.packageRepository.GetByID(ctx, nil, user.PackageID.String())
		if err != nil {
			return nil, err
		}

		pkgs = []entity.Package{
			pkg,
		}
	}

	var pkgInfos []dto.PackageInfo
	for _, pkg := range pkgs {
		pkgInfos = append(pkgInfos, pkg.ToInfo())
	}

	return pkgInfos, nil
}

func (s *packageService) UpdatePackage(ctx context.Context, req dto.UpdatePackageRequest) (dto.PackageInfo, error) {
	pkg, err := s.packageRepository.GetByID(ctx, nil, req.ID)
	if err != nil {
		return dto.PackageInfo{}, err
	}
	pkg.Name = req.Name

	pkg, err = s.packageRepository.Update(ctx, nil, pkg)
	if err != nil {
		return dto.PackageInfo{}, err
	}

	return pkg.ToInfo(), nil
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

	return pkg.ToInfo(), nil
}
