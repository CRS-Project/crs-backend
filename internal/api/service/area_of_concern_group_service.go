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
	AreaOfConcernGroupService interface {
		Create(ctx context.Context, req dto.AreaOfConcernGroupRequest) (dto.AreaOfConcernGroupResponse, error)
		GetById(ctx context.Context, areaOfConcernGroupId string) (dto.AreaOfConcernGroupResponse, error)
		GetAll(ctx context.Context, userId string, metaReq meta.Meta) ([]dto.AreaOfConcernGroupResponse, meta.Meta, error)
		Update(ctx context.Context, req dto.AreaOfConcernGroupRequest) error
		Delete(ctx context.Context, userId, areaOfConcernGroupId string) error
	}

	areaOfConcernGroupService struct {
		areaOfConcernGroupRepository repository.AreaOfConcernGroupRepository
		packageRepository            repository.PackageRepository
		userRepository               repository.UserRepository
		userDisciplineRepository     repository.UserDisciplineRepository
		db                           *gorm.DB
	}
)

func NewAreaOfConcernGroup(areaOfConcernGroupRepository repository.AreaOfConcernGroupRepository,
	packageRepository repository.PackageRepository,
	userRepository repository.UserRepository,
	userDisciplineRepository repository.UserDisciplineRepository,
	db *gorm.DB) AreaOfConcernGroupService {
	return &areaOfConcernGroupService{
		areaOfConcernGroupRepository: areaOfConcernGroupRepository,
		packageRepository:            packageRepository,
		userRepository:               userRepository,
		userDisciplineRepository:     userDisciplineRepository,
		db:                           db,
	}
}

func (s *areaOfConcernGroupService) Create(ctx context.Context, req dto.AreaOfConcernGroupRequest) (dto.AreaOfConcernGroupResponse, error) {
	pkg, _, err := s.getPackagePermission(ctx, req.UserId)
	if err != nil {
		return dto.AreaOfConcernGroupResponse{}, err
	}

	areaOfConcernGroupResult, err := s.areaOfConcernGroupRepository.Create(ctx, nil, entity.AreaOfConcernGroup{})
	if err != nil {
		return dto.AreaOfConcernGroupResponse{}, err
	}

	userDiscipline, err := s.userDisciplineRepository.GetByID(ctx, nil, req.UserDisciplineID)
	if err != nil {
		return dto.AreaOfConcernGroupResponse{}, err
	}

	return dto.AreaOfConcernGroupResponse{
		ID:             areaOfConcernGroupResult.ID.String(),
		Package:        pkg.Name,
		UserDiscipline: userDiscipline.Name,
	}, nil
}

func (s *areaOfConcernGroupService) GetById(ctx context.Context, id string) (dto.AreaOfConcernGroupResponse, error) {
	areaOfConcernGroup, err := s.areaOfConcernGroupRepository.GetByID(ctx, nil, id, "Package", "UserDiscipline")
	if err != nil {
		return dto.AreaOfConcernGroupResponse{}, err
	}

	return dto.AreaOfConcernGroupResponse{
		ID:             areaOfConcernGroup.ID.String(),
		ReviewFocus:    areaOfConcernGroup.ReviewFocus,
		Package:        areaOfConcernGroup.Package.Name,
		UserDiscipline: areaOfConcernGroup.UserDiscipline.Name,
	}, nil
}

func (s *areaOfConcernGroupService) GetAll(ctx context.Context, userId string, metaReq meta.Meta) ([]dto.AreaOfConcernGroupResponse, meta.Meta, error) {
	pkg, _, err := s.getPackagePermission(ctx, userId)
	if err != nil {
		return nil, meta.Meta{}, err
	}

	packageId := ""
	if pkg != nil {
		packageId = pkg.ID.String()
	}

	areaOfConcernGroups, metaRes, err := s.areaOfConcernGroupRepository.GetAll(ctx, nil, packageId, metaReq, "Package", "UserDiscipline")
	if err != nil {
		return nil, meta.Meta{}, err
	}

	var areaOfConcernGroupResponse []dto.AreaOfConcernGroupResponse
	for _, areaOfConcernGroup := range areaOfConcernGroups {
		areaOfConcernGroupResponse = append(areaOfConcernGroupResponse, dto.AreaOfConcernGroupResponse{
			ID:             areaOfConcernGroup.ID.String(),
			ReviewFocus:    areaOfConcernGroup.ReviewFocus,
			Package:        areaOfConcernGroup.Package.Name,
			UserDiscipline: areaOfConcernGroup.UserDiscipline.Name,
		})
	}

	return areaOfConcernGroupResponse, metaRes, nil
}

func (s *areaOfConcernGroupService) Update(ctx context.Context, req dto.AreaOfConcernGroupRequest) error {
	pkg, _, err := s.getPackagePermission(ctx, req.UserId)
	if err != nil {
		return err
	}

	areaOfConcernGroup, err := s.areaOfConcernGroupRepository.GetByID(ctx, nil, req.ID)
	if err != nil {
		return err
	}

	if pkg != nil && pkg.ID != areaOfConcernGroup.PackageID {
		return myerror.New("you not allowed to this package", http.StatusUnauthorized)
	}

	userDiscipline, err := s.userDisciplineRepository.GetByID(ctx, nil, req.UserDisciplineID)
	if err != nil {
		return err
	}

	areaOfConcernGroup.UserDisciplineID = userDiscipline.ID
	areaOfConcernGroup.ReviewFocus = req.ReviewFocus

	if err = s.areaOfConcernGroupRepository.Update(ctx, nil, areaOfConcernGroup); err != nil {
		return err
	}

	return nil
}

func (s *areaOfConcernGroupService) Delete(ctx context.Context, userId, areaOfConcernGroupId string) error {
	pkg, _, err := s.getPackagePermission(ctx, userId)
	if err != nil {
		return err
	}

	areaOfConcernGroup, err := s.areaOfConcernGroupRepository.GetByID(ctx, nil, areaOfConcernGroupId)
	if err != nil {
		return err
	}

	if pkg != nil && pkg.ID != areaOfConcernGroup.PackageID {
		return myerror.New("you not allowed to this package", http.StatusUnauthorized)
	}

	if err = s.areaOfConcernGroupRepository.Delete(ctx, nil, areaOfConcernGroup); err != nil {
		return err
	}

	return nil
}

func (s *areaOfConcernGroupService) getPackagePermission(ctx context.Context, userId string) (*entity.Package, entity.User, error) {
	user, err := s.userRepository.GetById(ctx, nil, userId)
	if err != nil {
		return nil, entity.User{}, err
	}

	if user.PackageID == nil {
		return nil, user, nil
	}

	pkg, err := s.packageRepository.GetByID(ctx, nil, user.PackageID.String())
	if err != nil {
		return nil, entity.User{}, err
	}

	return &pkg, user, nil
}
