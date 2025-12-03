package service

import (
	"context"
	"errors"
	"net/http"

	"github.com/CRS-Project/crs-backend/internal/api/repository"
	"github.com/CRS-Project/crs-backend/internal/dto"
	"github.com/CRS-Project/crs-backend/internal/entity"
	myerror "github.com/CRS-Project/crs-backend/internal/pkg/error"
	"github.com/CRS-Project/crs-backend/internal/pkg/meta"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type (
	AreaOfConcernService interface {
		Create(ctx context.Context, req dto.AreaOfConcernRequest) (dto.AreaOfConcernResponse, error)
		GetById(ctx context.Context, areaOfConcernId string) (dto.AreaOfConcernResponse, error)
		GetAll(ctx context.Context, areaOfConcernGroupId, userId string, metaReq meta.Meta) ([]dto.AreaOfConcernResponse, meta.Meta, error)
		Update(ctx context.Context, req dto.UpdateAreaOfConcernRequest) error
		Delete(ctx context.Context, userId, areaOfConcernId string) error
	}

	areaOfConcernService struct {
		areaOfConcernRepository             repository.AreaOfConcernRepository
		areaOfConcernGroupRepository        repository.AreaOfConcernGroupRepository
		areaOfConcernConsolidatorRepository repository.AreaOfConcernConsolidatorRepository
		commentRepository                   repository.CommentRepository
		packageRepository                   repository.PackageRepository
		userRepository                      repository.UserRepository
		userDisciplineRepository            repository.UserDisciplineRepository
		db                                  *gorm.DB
	}
)

func NewAreaOfConcern(areaOfConcernRepository repository.AreaOfConcernRepository,
	areaOfConcernGroupRepository repository.AreaOfConcernGroupRepository,
	areaOfConcernConsolidatorRepository repository.AreaOfConcernConsolidatorRepository,
	commentRepository repository.CommentRepository,
	packageRepository repository.PackageRepository,
	userRepository repository.UserRepository,
	userDisciplineRepository repository.UserDisciplineRepository,
	db *gorm.DB) AreaOfConcernService {
	return &areaOfConcernService{
		areaOfConcernRepository:             areaOfConcernRepository,
		areaOfConcernGroupRepository:        areaOfConcernGroupRepository,
		areaOfConcernConsolidatorRepository: areaOfConcernConsolidatorRepository,
		commentRepository:                   commentRepository,
		packageRepository:                   packageRepository,
		userRepository:                      userRepository,
		userDisciplineRepository:            userDisciplineRepository,
		db:                                  db,
	}
}

func (s *areaOfConcernService) Create(ctx context.Context, req dto.AreaOfConcernRequest) (dto.AreaOfConcernResponse, error) {
	pkg, user, err := s.getPackagePermission(ctx, req.UserId)
	if err != nil {
		return dto.AreaOfConcernResponse{}, err
	}

	var contractor entity.User
	if pkg == nil {
		contractor, err = s.userRepository.GetContractorByPackage(ctx, nil, req.PackageID, "Package")
		if err != nil {
			if errors.Is(err,gorm.ErrRecordNotFound){
				return dto.AreaOfConcernResponse{}, myerror.New("this package not have contractor", http.StatusBadRequest)
			}
			return dto.AreaOfConcernResponse{}, err
		}

		pkg = contractor.Package
	} else {
		contractor = user
	}

	var consolidatorsInput []entity.AreaOfConcernConsolidator
	for _, consolidator := range req.Consolidators {
		consolidatorsInput = append(consolidatorsInput, entity.AreaOfConcernConsolidator{
			UserID: uuid.MustParse(consolidator.UserID),
		})
	}

	areaofconcerngroup, err := s.areaOfConcernGroupRepository.GetByID(ctx, nil, req.AreaOfConcernGroupID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return dto.AreaOfConcernResponse{}, myerror.New("area of concern group not found", http.StatusNotFound)
		}
		return dto.AreaOfConcernResponse{}, err
	}

	areaOfConcernResult, err := s.areaOfConcernRepository.Create(ctx, nil, entity.AreaOfConcern{
		Description:          req.Description,
		AreaOfConcernId:      req.AreaOfConcernId,
		AreaOfConcernGroupID: areaofconcerngroup.ID,
		PackageID:            pkg.ID,
		Consolidators:        consolidatorsInput,
	})
	if err != nil {
		return dto.AreaOfConcernResponse{}, err
	}

	return dto.AreaOfConcernResponse{
		ID:              areaOfConcernResult.ID.String(),
		AreaOfConcernId: req.AreaOfConcernId,
		Description:     req.Description,
		Package:         pkg.Name,
	}, nil
}

func (s *areaOfConcernService) GetById(ctx context.Context, id string) (dto.AreaOfConcernResponse, error) {
	areaOfConcern, err := s.areaOfConcernRepository.GetByID(ctx, nil, id, "Package", "Consolidators.User")
	if err != nil {
		return dto.AreaOfConcernResponse{}, err
	}

	var consolidatorResponse []dto.AreaOfConcernConsolidatorResponse
	for _, c := range areaOfConcern.Consolidators {
		consolidatorResponse = append(consolidatorResponse, dto.AreaOfConcernConsolidatorResponse{
			ID:   c.User.ID.String(),
			Name: c.User.Name,
		})
	}

	return dto.AreaOfConcernResponse{
		ID:              areaOfConcern.ID.String(),
		AreaOfConcernId: areaOfConcern.AreaOfConcernId,
		Description:     areaOfConcern.Description,
		Package:         areaOfConcern.Package.Name,
		Consolidators:   consolidatorResponse,
	}, nil
}

func (s *areaOfConcernService) GetAll(ctx context.Context, areaOfConcernGroupId, userId string, metaReq meta.Meta) ([]dto.AreaOfConcernResponse, meta.Meta, error) {
	_, _, err := s.getPackagePermission(ctx, userId)
	if err != nil {
		return nil, meta.Meta{}, err
	}

	areaOfConcerns, metaRes, err := s.areaOfConcernRepository.GetAllByAreaOfConcernGroupID(ctx, nil, areaOfConcernGroupId, metaReq, "Package")
	if err != nil {
		return nil, meta.Meta{}, err
	}

	var areaOfConcernResponse []dto.AreaOfConcernResponse
	for _, areaOfConcern := range areaOfConcerns {
		areaOfConcernResponse = append(areaOfConcernResponse, dto.AreaOfConcernResponse{
			ID:              areaOfConcern.ID.String(),
			AreaOfConcernId: areaOfConcern.AreaOfConcernId,
			Description:     areaOfConcern.Description,
			Package:         areaOfConcern.Package.Name,
		})
	}

	return areaOfConcernResponse, metaRes, nil
}

func (s *areaOfConcernService) Update(ctx context.Context, req dto.UpdateAreaOfConcernRequest) error {
	pkg, _, err := s.getPackagePermission(ctx, req.UserId)
	if err != nil {
		return err
	}

	areaOfConcern, err := s.areaOfConcernRepository.GetByID(ctx, nil, req.ID, "Consolidators.User")
	if err != nil {
		return err
	}

	if pkg != nil && pkg.ID != areaOfConcern.PackageID {
		return myerror.New("you not allowed to this package", http.StatusUnauthorized)
	}

	areaOfConcern.AreaOfConcernId = req.AreaOfConcernId
	areaOfConcern.Description = req.Description

	consolidatorMap := map[string]bool{}
	for _, c := range areaOfConcern.Consolidators {
		consolidatorMap[c.UserID.String()] = true
	}

	reqConsolidatorMap := map[string]bool{}
	for _, c := range req.Consolidators {
		reqConsolidatorMap[c.UserID] = true
	}

	var deletedConsolidators []string
	for _, c := range areaOfConcern.Consolidators {
		if _, ok := reqConsolidatorMap[c.UserID.String()]; !ok {
			deletedConsolidators = append(deletedConsolidators, c.ID.String())
		}
	}

	if err := s.areaOfConcernConsolidatorRepository.DeleteBulk(ctx, nil, deletedConsolidators); err != nil {
		return err
	}

	var newConsolidator []entity.AreaOfConcernConsolidator
	for _, c := range req.Consolidators {
		if _, ok := consolidatorMap[c.UserID]; !ok {
			newConsolidator = append(newConsolidator, entity.AreaOfConcernConsolidator{
				AreaOfConcernID: areaOfConcern.ID,
				UserID:          uuid.MustParse(c.UserID),
			})
		}
	}

	if err := s.areaOfConcernConsolidatorRepository.CreateBulk(ctx, nil, newConsolidator); err != nil {
		return err
	}

	if err = s.areaOfConcernRepository.Update(ctx, nil, areaOfConcern); err != nil {
		return err
	}

	return nil
}

func (s *areaOfConcernService) Delete(ctx context.Context, userId, areaOfConcernId string) error {
	pkg, _, err := s.getPackagePermission(ctx, userId)
	if err != nil {
		return err
	}

	areaOfConcern, err := s.areaOfConcernRepository.GetByID(ctx, nil, areaOfConcernId)
	if err != nil {
		return err
	}

	if pkg != nil && pkg.ID != areaOfConcern.PackageID {
		return myerror.New("you not allowed to this package", http.StatusUnauthorized)
	}

	if err := s.commentRepository.DeleteByAreaOfConcernID(ctx, nil, []string{areaOfConcern.ID.String()}); err != nil {
		return err
	}

	if err = s.areaOfConcernRepository.Delete(ctx, nil, areaOfConcern); err != nil {
		return err
	}

	return nil
}

func (s *areaOfConcernService) getPackagePermission(ctx context.Context, userId string) (*entity.Package, entity.User, error) {
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
