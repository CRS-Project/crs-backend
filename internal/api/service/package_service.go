package service

import (
	"net/http"

	"github.com/CRS-Project/crs-backend/internal/api/repository"
	"github.com/CRS-Project/crs-backend/internal/dto"
	"github.com/CRS-Project/crs-backend/internal/entity"
	myerror "github.com/CRS-Project/crs-backend/internal/pkg/error"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type (
	PackageService interface {
		CreatePackage(ctx *gin.Context, req dto.CreatePackageRequest) (dto.PackageInfo, error)
		GetPackages(ctx *gin.Context) ([]dto.PackageInfo, error)
		UpdatePackage(ctx *gin.Context, req dto.UpdatePackageRequest) error
		DeletePackage(ctx *gin.Context, req dto.DeletePackageRequest) error
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

func (s *packageService) CreatePackage(ctx *gin.Context, req dto.CreatePackageRequest) (dto.PackageInfo, error) {
	_, err := s.packageRepository.GetByName(ctx, s.db, req.Name)
	if err == nil {
		return dto.PackageInfo{}, myerror.New("package with this name already exists", http.StatusConflict)
	}

	pkgCreation := entity.Package{
		Name: req.Name,
	}

	pkgResult, err := s.packageRepository.Create(ctx, s.db, pkgCreation)
	if err != nil {
		return dto.PackageInfo{}, err
	}

	return dto.PackageInfo{
		ID:   pkgResult.ID.String(),
		Name: pkgResult.Name,
	}, nil
}

func (s *packageService) GetPackages(ctx *gin.Context) ([]dto.PackageInfo, error) {
	pkgs, err := s.packageRepository.GetAll(ctx, s.db)
	if err != nil {
		return []dto.PackageInfo{}, err
	}

	var pkgInfos []dto.PackageInfo
	for _, pkg := range pkgs {
		pkgInfos = append(pkgInfos, dto.PackageInfo{
			ID:   pkg.ID.String(),
			Name: pkg.Name,
		})
	}

	return pkgInfos, nil
}

func (s *packageService) UpdatePackage(ctx *gin.Context, req dto.UpdatePackageRequest) error {
	uuid, err := uuid.Parse(req.ID)
	if err != nil {
		return err
	}

	_, err = s.packageRepository.GetByID(ctx, s.db, uuid)
	if err != nil {
		return err
	}

	pkg := entity.Package{
		ID:   uuid,
		Name: req.Name,
	}
	if err = s.packageRepository.Update(ctx, s.db, pkg); err != nil {
		return err
	}

	return nil
}

func (s *packageService) DeletePackage(ctx *gin.Context, req dto.DeletePackageRequest) error {
	uuid, err := uuid.Parse(req.ID)
	if err != nil {
		return err
	}

	pkg, err := s.packageRepository.GetByID(ctx, s.db, uuid)
	if err != nil {
		return err
	}

	if err = s.packageRepository.Delete(ctx, s.db, pkg); err != nil {
		return err
	}

	return nil
}
