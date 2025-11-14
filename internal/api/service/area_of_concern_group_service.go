package service

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/CRS-Project/crs-backend/internal/api/repository"
	"github.com/CRS-Project/crs-backend/internal/dto"
	"github.com/CRS-Project/crs-backend/internal/entity"
	myerror "github.com/CRS-Project/crs-backend/internal/pkg/error"
	"github.com/CRS-Project/crs-backend/internal/pkg/meta"
	mypdf "github.com/CRS-Project/crs-backend/internal/pkg/pdf"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type (
	AreaOfConcernGroupService interface {
		Create(ctx context.Context, req dto.AreaOfConcernGroupRequest) (dto.AreaOfConcernGroupResponse, error)
		GetById(ctx context.Context, areaOfConcernGroupId string) (dto.AreaOfConcernGroupResponse, error)
		GetAll(ctx context.Context, userId string, metaReq meta.Meta) ([]dto.AreaOfConcernGroupResponse, meta.Meta, error)
		Update(ctx context.Context, req dto.AreaOfConcernGroupRequest) error
		Delete(ctx context.Context, userId, areaOfConcernGroupId string) error
		GeneratePDF(ctx context.Context, userId, areaOfConcernGroupId string) (*bytes.Buffer, string, error)
		GetStatistic(ctx context.Context, packageId string) (dto.AreaOfConcernGroupStatistic, error)
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
	pkg, user, err := s.getPackagePermission(ctx, req.UserId)
	if err != nil {
		return dto.AreaOfConcernGroupResponse{}, err
	}

	var contractor entity.User
	if pkg == nil {
		contractor, err = s.userRepository.GetContractorByPackage(ctx, nil, req.PackageID, "Package")
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return dto.AreaOfConcernGroupResponse{}, myerror.New("this package not have contractor", http.StatusBadRequest)
			}
			return dto.AreaOfConcernGroupResponse{}, err
		}

		pkg = contractor.Package
	} else {
		contractor = user
	}

	areaOfConcernGroupResult, err := s.areaOfConcernGroupRepository.Create(ctx, nil, entity.AreaOfConcernGroup{
		ReviewFocus:    req.ReviewFocus,
		UserDiscipline: req.UserDiscipline,
		PackageID:      uuid.MustParse(req.PackageID),
	})
	if err != nil {
		return dto.AreaOfConcernGroupResponse{}, err
	}

	return dto.AreaOfConcernGroupResponse{
		ID:             areaOfConcernGroupResult.ID.String(),
		ReviewFocus:    areaOfConcernGroupResult.ReviewFocus,
		Package:        pkg.Name,
		UserDiscipline: areaOfConcernGroupResult.UserDiscipline,
	}, nil
}

func (s *areaOfConcernGroupService) GetById(ctx context.Context, id string) (dto.AreaOfConcernGroupResponse, error) {
	areaOfConcernGroup, err := s.areaOfConcernGroupRepository.GetByID(ctx, nil, id, "Package")
	if err != nil {
		return dto.AreaOfConcernGroupResponse{}, err
	}

	return dto.AreaOfConcernGroupResponse{
		ID:             areaOfConcernGroup.ID.String(),
		ReviewFocus:    areaOfConcernGroup.ReviewFocus,
		Package:        areaOfConcernGroup.Package.Name,
		UserDiscipline: areaOfConcernGroup.UserDiscipline,
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

	areaOfConcernGroups, metaRes, err := s.areaOfConcernGroupRepository.GetAll(ctx, nil, packageId, metaReq, "Package")
	if err != nil {
		return nil, meta.Meta{}, err
	}

	var areaOfConcernGroupResponse []dto.AreaOfConcernGroupResponse
	for _, areaOfConcernGroup := range areaOfConcernGroups {
		areaOfConcernGroupResponse = append(areaOfConcernGroupResponse, dto.AreaOfConcernGroupResponse{
			ID:             areaOfConcernGroup.ID.String(),
			ReviewFocus:    areaOfConcernGroup.ReviewFocus,
			Package:        areaOfConcernGroup.Package.Name,
			UserDiscipline: areaOfConcernGroup.UserDiscipline,
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

	areaOfConcernGroup.UserDiscipline = req.UserDiscipline
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

func (s *areaOfConcernGroupService) GeneratePDF(ctx context.Context, userId, areaOfConcernGroupId string) (*bytes.Buffer, string, error) {
	data, err := s.areaOfConcernGroupRepository.GetByID(ctx, nil, areaOfConcernGroupId, "AreaOfConcerns.Consolidators.User", "AreaOfConcerns.Comments.CommentReplies", "AreaOfConcerns.Comments.User", "AreaOfConcerns.Comments.Document", "Package")
	if err != nil {
		return nil, "", err
	}

	contractor, err := s.userRepository.GetContractorByPackage(ctx, nil, data.PackageID.String(), "Package")
	if err != nil {
		return nil, "", err
	}

	var requestData []mypdf.GenerateRequestData
	for _, aoc := range data.AreaOfConcerns {
		consolidator := ""
		for i, c := range aoc.Consolidators {
			userName := "deleted user"
			if c.User != nil {
				userName = c.User.Name
			}
			if i > 0 {
				consolidator += fmt.Sprintf("\n%d. %s", i+1, userName)
			} else {
				consolidator += fmt.Sprintf("%d. %s", i+1, userName)
			}
		}

		var comments []mypdf.CommentRow
		for i, c := range aoc.Comments {
			status := "N/A"
			if c.Status != nil {
				status = string(*c.Status)
			}

			closeOutComments := "N/A"
			for _, cr := range c.CommentReplies {
				if cr.IsCloseOutComment {
					closeOutComments = cr.Comment
					break
				}
			}

			comments = append(comments, mypdf.CommentRow{
				No:              fmt.Sprintf("%d", i+1),
				Page:            c.Section,
				SMEInitial:      c.User.Name,
				SMEComment:      c.Comment,
				RefDocNo:        c.Document.CompanyDocumentNumber,
				RefDocTitle:     c.Document.DocumentTitle,
				DocStatus:       string(c.Document.Status),
				Status:          status,
				SMECloseComment: closeOutComments,
			})
		}

		requestData = append(requestData, mypdf.GenerateRequestData{
			PackageInfoData: mypdf.PackageInfoData{
				Package:           contractor.Package.Name,
				ContractorInitial: contractor.Name,
			},
			DisciplineSectionData: mypdf.DisciplineSectionData{
				Discipline:               data.UserDiscipline,
				AreaOfConcernID:          aoc.AreaOfConcernId,
				AreaOfConcernDescription: aoc.Description,
				Consolidator:             consolidator,
			},
			CommentRow: comments,
		})
	}

	pdfBuffer, filename, err := mypdf.Generate(requestData)
	if err != nil {
		return nil, "", err
	}

	return pdfBuffer, filename, nil
}

func (s *areaOfConcernGroupService) GetStatistic(ctx context.Context, packageId string) (dto.AreaOfConcernGroupStatistic, error) {
	return s.areaOfConcernGroupRepository.Statistic(ctx, nil, packageId)
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
