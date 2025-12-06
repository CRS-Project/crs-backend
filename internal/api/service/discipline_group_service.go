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
	DisciplineGroupService interface {
		Create(ctx context.Context, req dto.DisciplineGroupRequest) (dto.DisciplineGroupResponse, error)
		GetById(ctx context.Context, disciplineGroupId string) (dto.DisciplineGroupResponse, error)
		GetAll(ctx context.Context, userId string, metaReq meta.Meta) ([]dto.DisciplineGroupResponse, meta.Meta, error)
		GetAllConsolidator(ctx context.Context, search, userId string) ([]dto.DisciplineGroupConsolidatorResponse, error)
		Update(ctx context.Context, req dto.DisciplineGroupRequest) error
		Delete(ctx context.Context, userId, disciplineGroupId string) error
		GeneratePDF(ctx context.Context, userId, disciplineGroupId string) (*bytes.Buffer, string, error)
		GetStatistic(ctx context.Context, packageId string) (dto.DisciplineGroupStatistic, error)
		ConstructGeneratePDF(disciplineGroup entity.DisciplineGroup, contractor entity.User) []mypdf.GenerateRequestData
	}

	disciplineGroupService struct {
		disciplineGroupRepository                    repository.DisciplineGroupRepository
		disciplineGroupConsolidatorRepository        repository.DisciplineGroupConsolidatorRepository
		disciplineListDocumentRepository             repository.DisciplineListDocumentRepository
		disciplineListDocumentConsolidatorRepository repository.DisciplineListDocumentConsolidatorRepository
		packageRepository                            repository.PackageRepository
		commentRepository                            repository.CommentRepository
		userRepository                               repository.UserRepository
		userDisciplineRepository                     repository.UserDisciplineRepository
		db                                           *gorm.DB
	}
)

func NewDisciplineGroup(disciplineGroupRepository repository.DisciplineGroupRepository,
	disciplineGroupConsolidatorRepository repository.DisciplineGroupConsolidatorRepository,
	disciplineListDocumentRepository repository.DisciplineListDocumentRepository,
	disciplineListDocumentConsolidatorRepository repository.DisciplineListDocumentConsolidatorRepository,
	packageRepository repository.PackageRepository,
	commentRepository repository.CommentRepository,
	userRepository repository.UserRepository,
	userDisciplineRepository repository.UserDisciplineRepository,
	db *gorm.DB) DisciplineGroupService {
	return &disciplineGroupService{
		disciplineGroupRepository:                    disciplineGroupRepository,
		disciplineGroupConsolidatorRepository:        disciplineGroupConsolidatorRepository,
		disciplineListDocumentRepository:             disciplineListDocumentRepository,
		disciplineListDocumentConsolidatorRepository: disciplineListDocumentConsolidatorRepository,
		packageRepository:                            packageRepository,
		commentRepository:                            commentRepository,
		userRepository:                               userRepository,
		userDisciplineRepository:                     userDisciplineRepository,
		db:                                           db,
	}
}

func (s *disciplineGroupService) Create(ctx context.Context, req dto.DisciplineGroupRequest) (dto.DisciplineGroupResponse, error) {
	pkg, user, err := s.getPackagePermission(ctx, req.UserId)
	if err != nil {
		return dto.DisciplineGroupResponse{}, err
	}

	var contractor entity.User
	if pkg == nil {
		contractor, err = s.userRepository.GetContractorByPackage(ctx, nil, req.PackageID, "Package")
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return dto.DisciplineGroupResponse{}, myerror.New("this package not have contractor", http.StatusBadRequest)
			}
			return dto.DisciplineGroupResponse{}, err
		}

		pkg = contractor.Package
	} else {
		contractor = user
	}

	var consolidatorsInput []entity.DisciplineGroupConsolidator
	for _, consolidator := range req.DisciplineGroupConsolidators {
		consolidatorsInput = append(consolidatorsInput, entity.DisciplineGroupConsolidator{
			UserID: uuid.MustParse(consolidator.UserID),
		})
	}

	disciplineGroupResult, err := s.disciplineGroupRepository.Create(ctx, nil, entity.DisciplineGroup{
		ReviewFocus:                  req.ReviewFocus,
		UserDiscipline:               req.UserDiscipline,
		DisciplineInitial:            req.DisciplineInitial,
		PackageID:                    uuid.MustParse(req.PackageID),
		DisciplineGroupConsolidators: consolidatorsInput,
	})
	if err != nil {
		return dto.DisciplineGroupResponse{}, err
	}

	return dto.DisciplineGroupResponse{
		ID:                disciplineGroupResult.ID.String(),
		ReviewFocus:       disciplineGroupResult.ReviewFocus,
		DisciplineInitial: disciplineGroupResult.DisciplineInitial,
		Package:           pkg.Name,
		UserDiscipline:    disciplineGroupResult.UserDiscipline,
	}, nil
}

func (s *disciplineGroupService) GetById(ctx context.Context, id string) (dto.DisciplineGroupResponse, error) {
	disciplineGroup, err := s.disciplineGroupRepository.GetByID(ctx, nil, id, "Package", "DisciplineGroupConsolidators.User")
	if err != nil {
		return dto.DisciplineGroupResponse{}, err
	}

	var consolidatorResponse []dto.DisciplineGroupConsolidatorResponse
	for _, c := range disciplineGroup.DisciplineGroupConsolidators {
		consolidatorResponse = append(consolidatorResponse, dto.DisciplineGroupConsolidatorResponse{
			ID:   c.User.ID.String(),
			Name: c.User.Name,
		})
	}

	return dto.DisciplineGroupResponse{
		ID:                           disciplineGroup.ID.String(),
		ReviewFocus:                  disciplineGroup.ReviewFocus,
		Package:                      disciplineGroup.Package.Name,
		UserDiscipline:               disciplineGroup.UserDiscipline,
		DisciplineInitial:            disciplineGroup.DisciplineInitial,
		DisciplineGroupConsolidators: consolidatorResponse,
	}, nil
}

func (s *disciplineGroupService) GetAll(ctx context.Context, userId string, metaReq meta.Meta) ([]dto.DisciplineGroupResponse, meta.Meta, error) {
	pkg, _, err := s.getPackagePermission(ctx, userId)
	if err != nil {
		return nil, meta.Meta{}, err
	}

	packageId := ""
	if pkg != nil {
		packageId = pkg.ID.String()
	}

	disciplineGroups, metaRes, err := s.disciplineGroupRepository.GetAll(ctx, nil, packageId, metaReq, "Package")
	if err != nil {
		return nil, meta.Meta{}, err
	}

	var disciplineGroupResponse []dto.DisciplineGroupResponse
	for _, disciplineGroup := range disciplineGroups {
		disciplineGroupResponse = append(disciplineGroupResponse, dto.DisciplineGroupResponse{
			ID:                disciplineGroup.ID.String(),
			ReviewFocus:       disciplineGroup.ReviewFocus,
			Package:           disciplineGroup.Package.Name,
			UserDiscipline:    disciplineGroup.UserDiscipline,
			DisciplineInitial: disciplineGroup.DisciplineInitial,
		})
	}

	return disciplineGroupResponse, metaRes, nil
}

func (s *disciplineGroupService) GetAllConsolidator(ctx context.Context, search, disciplineGroupId string) ([]dto.DisciplineGroupConsolidatorResponse, error) {
	consolidator, err := s.disciplineGroupConsolidatorRepository.GetAllConsolidator(ctx, nil, search, disciplineGroupId, "User")
	if err != nil {
		return nil, err
	}

	var response []dto.DisciplineGroupConsolidatorResponse
	for _, c := range consolidator {
		response = append(response, dto.DisciplineGroupConsolidatorResponse{
			ID:   c.UserID.String(),
			Name: c.User.Name,
		})
	}

	return response, nil
}

func (s *disciplineGroupService) Update(ctx context.Context, req dto.DisciplineGroupRequest) error {
	// Validate user permission
	pkg, _, err := s.getPackagePermission(ctx, req.UserId)
	if err != nil {
		return err
	}

	disciplineGroup, err := s.disciplineGroupRepository.GetByID(ctx, nil, req.ID)
	if err != nil {
		return err
	}

	if pkg != nil && pkg.ID != disciplineGroup.PackageID {
		return myerror.New("you not allowed to this package", http.StatusUnauthorized)
	}

	disciplineGroup.UserDiscipline = req.UserDiscipline
	disciplineGroup.ReviewFocus = req.ReviewFocus
	disciplineGroup.DisciplineInitial = req.DisciplineInitial

	if err := s.disciplineGroupConsolidatorRepository.DeleteByDisciplineGroupID(ctx, nil, disciplineGroup.ID.String()); err != nil {
		return err
	}

	var newConsolidators []entity.DisciplineGroupConsolidator
	for _, c := range req.DisciplineGroupConsolidators {
		uid, err := uuid.Parse(c.UserID)
		if err != nil {
			return myerror.New("invalid consolidator user_id: "+c.UserID, http.StatusBadRequest)
		}

		newConsolidators = append(newConsolidators, entity.DisciplineGroupConsolidator{
			DisciplineGroupID: disciplineGroup.ID,
			UserID:            uid,
		})
	}

	if len(newConsolidators) > 0 {
		if err := s.disciplineGroupConsolidatorRepository.CreateBulk(ctx, nil, newConsolidators); err != nil {
			return err
		}
	}

	if err = s.disciplineGroupRepository.Update(ctx, nil, disciplineGroup); err != nil {
		return err
	}

	return nil
}

func (s *disciplineGroupService) Delete(ctx context.Context, userId, disciplineGroupId string) error {
	pkg, _, err := s.getPackagePermission(ctx, userId)
	if err != nil {
		return err
	}

	disciplineGroup, err := s.disciplineGroupRepository.GetByID(ctx, nil, disciplineGroupId, "DisciplineListDocuments")
	if err != nil {
		return err
	}

	if pkg != nil && pkg.ID != disciplineGroup.PackageID {
		return myerror.New("you not allowed to this package", http.StatusUnauthorized)
	}

	var disciplineListDocumentIDs []string
	for _, dld := range disciplineGroup.DisciplineListDocuments {
		disciplineListDocumentIDs = append(disciplineListDocumentIDs, dld.ID.String())
	}

	if len(disciplineListDocumentIDs) > 0 {
		if err := s.commentRepository.DeleteByDisciplineListDocumentID(ctx, nil, disciplineListDocumentIDs); err != nil {
			return err
		}

		if err := s.disciplineListDocumentConsolidatorRepository.DeleteByDisciplineListDocumentID(ctx, nil, disciplineListDocumentIDs); err != nil {
			return err
		}
	}

	if err := s.disciplineListDocumentRepository.DeleteByDisciplineGroupID(ctx, nil, disciplineGroup.ID.String()); err != nil {
		return err
	}

	if err = s.disciplineGroupRepository.Delete(ctx, nil, disciplineGroup); err != nil {
		return err
	}

	return nil
}

func (s *disciplineGroupService) GeneratePDF(ctx context.Context, userId, disciplineGroupId string) (*bytes.Buffer, string, error) {
	data, err := s.disciplineGroupRepository.GetByID(ctx, nil, disciplineGroupId, "DisciplineGroupConsolidators.User", "DisciplineListDocuments.Comments.CommentReplies", "DisciplineListDocuments.Document", "DisciplineListDocuments.Comments.User", "Package")
	if err != nil {
		return nil, "", err
	}

	contractor, err := s.userRepository.GetContractorByPackage(ctx, nil, data.PackageID.String(), "Package")
	if err != nil {
		return nil, "", err
	}

	requestData := s.ConstructGeneratePDF(data, contractor)
	pdfBuffer, filename, err := mypdf.Generate(requestData)
	if err != nil {
		return nil, "", err
	}

	return pdfBuffer, filename, nil
}

func (s *disciplineGroupService) GetStatistic(ctx context.Context, packageId string) (dto.DisciplineGroupStatistic, error) {
	return s.disciplineGroupRepository.Statistic(ctx, nil, packageId)
}

func (s *disciplineGroupService) ConstructGeneratePDF(disciplineGroup entity.DisciplineGroup, contractor entity.User) []mypdf.GenerateRequestData {
	var requestData []mypdf.GenerateRequestData
	consolidator := ""
	for i, c := range disciplineGroup.DisciplineGroupConsolidators {
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

	for _, dld := range disciplineGroup.DisciplineListDocuments {
		var comments []mypdf.CommentRow
		for i, c := range dld.Comments {
			if c.CommentReplyID != nil {
				continue
			}

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
			refDocNo, refDocTitle, docStatus := "N/A", "N/A", "N/A"
			if dld.Document != nil {
				refDocNo = dld.Document.CompanyDocumentNumber
				refDocTitle = dld.Document.DocumentTitle
				docStatus = string(dld.Document.Status)
			}
			comments = append(comments, mypdf.CommentRow{
				No:              fmt.Sprintf("%d", i+1),
				Page:            c.Section,
				SMEInitial:      c.User.Name,
				SMEComment:      c.Comment,
				RefDocNo:        refDocNo,
				RefDocTitle:     refDocTitle,
				DocStatus:       string(docStatus),
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
				Discipline:               disciplineGroup.DisciplineInitial,
				AreaOfConcernID:          disciplineGroup.Package.Name + "-" + disciplineGroup.DisciplineInitial,
				AreaOfConcernDescription: disciplineGroup.UserDiscipline,
				Consolidator:             consolidator,
			},
			CommentRow: comments,
		})
	}

	return requestData
}

func (s *disciplineGroupService) getPackagePermission(ctx context.Context, userId string) (*entity.Package, entity.User, error) {
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
