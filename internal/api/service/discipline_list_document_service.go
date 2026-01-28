package service

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/CRS-Project/crs-backend/internal/api/repository"
	"github.com/CRS-Project/crs-backend/internal/dto"
	"github.com/CRS-Project/crs-backend/internal/entity"
	myerror "github.com/CRS-Project/crs-backend/internal/pkg/error"
	mylog "github.com/CRS-Project/crs-backend/internal/pkg/logger"
	"github.com/CRS-Project/crs-backend/internal/pkg/meta"
	mypdf "github.com/CRS-Project/crs-backend/internal/pkg/pdf"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type (
	DisciplineListDocumentService interface {
		Create(ctx context.Context, req dto.DisciplineListDocumentRequest) (dto.DisciplineListDocumentResponse, error)
		GetById(ctx context.Context, disciplineListDocumentId string) (dto.DisciplineListDocumentResponse, error)
		GetAll(ctx context.Context, disciplineGroupId, userId string, metaReq meta.Meta) ([]dto.DisciplineListDocumentResponse, meta.Meta, error)
		Update(ctx context.Context, req dto.UpdateDisciplineListDocumentRequest) error
		Delete(ctx context.Context, userId, disciplineListDocumentId string) error
		GenerateExcel(ctx context.Context, userId, disciplineListDocumentId string) (*bytes.Buffer, string, error)
	}

	disciplineListDocumentService struct {
		disciplineListDocumentRepository             repository.DisciplineListDocumentRepository
		disciplineGroupRepository                    repository.DisciplineGroupRepository
		disciplineListDocumentConsolidatorRepository repository.DisciplineListDocumentConsolidatorRepository
		commentRepository                            repository.CommentRepository
		packageRepository                            repository.PackageRepository
		documentRepository                           repository.DocumentRepository
		userRepository                               repository.UserRepository
		userDisciplineRepository                     repository.UserDisciplineRepository
		db                                           *gorm.DB
	}
)

func NewDisciplineListDocument(disciplineListDocumentRepository repository.DisciplineListDocumentRepository,
	disciplineGroupRepository repository.DisciplineGroupRepository,
	disciplineListDocumentConsolidatorRepository repository.DisciplineListDocumentConsolidatorRepository,
	commentRepository repository.CommentRepository,
	packageRepository repository.PackageRepository,
	documentRepository repository.DocumentRepository,
	userRepository repository.UserRepository,
	userDisciplineRepository repository.UserDisciplineRepository,
	db *gorm.DB) DisciplineListDocumentService {
	return &disciplineListDocumentService{
		disciplineListDocumentRepository:             disciplineListDocumentRepository,
		disciplineGroupRepository:                    disciplineGroupRepository,
		disciplineListDocumentConsolidatorRepository: disciplineListDocumentConsolidatorRepository,
		commentRepository:                            commentRepository,
		packageRepository:                            packageRepository,
		documentRepository:                           documentRepository,
		userRepository:                               userRepository,
		userDisciplineRepository:                     userDisciplineRepository,
		db:                                           db,
	}
}

func (s *disciplineListDocumentService) Create(ctx context.Context, req dto.DisciplineListDocumentRequest) (dto.DisciplineListDocumentResponse, error) {
	pkg, user, err := s.getPackagePermission(ctx, req.UserId)
	if err != nil {
		return dto.DisciplineListDocumentResponse{}, err
	}

	var contractor entity.User
	if pkg == nil {
		contractor, err = s.userRepository.GetContractorByPackage(ctx, nil, req.PackageID, "Package")
		if err != nil {
			return dto.DisciplineListDocumentResponse{}, err
		}

		pkg = contractor.Package
	} else {
		contractor = user
	}

	var consolidatorsInput []entity.DisciplineListDocumentConsolidator
	for _, consolidator := range req.Consolidators {
		consolidatorsInput = append(consolidatorsInput, entity.DisciplineListDocumentConsolidator{
			DisciplineGroupConsolidatorID: uuid.MustParse(consolidator.DisciplineGroupConsolidatorID),
		})
	}

	disciplinegroup, err := s.disciplineGroupRepository.GetByID(ctx, nil, req.DisciplineGroupID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return dto.DisciplineListDocumentResponse{}, myerror.New("discipline group not found", http.StatusNotFound)
		}
		return dto.DisciplineListDocumentResponse{}, err
	}

	document, err := s.documentRepository.GetByID(ctx, nil, req.DocumentID)
	if err != nil {
		return dto.DisciplineListDocumentResponse{}, err
	}

	disciplineListDocumentResult, err := s.disciplineListDocumentRepository.Create(ctx, nil, entity.DisciplineListDocument{
		DisciplineGroupID: disciplinegroup.ID,
		DocumentID:        document.ID,
		PackageID:         pkg.ID,
		Consolidators:     consolidatorsInput,
	})
	if err != nil {
		return dto.DisciplineListDocumentResponse{}, err
	}

	return dto.DisciplineListDocumentResponse{
		ID:      disciplineListDocumentResult.ID.String(),
		Package: pkg.Name,
		IsDueDate: document.DueDate != nil &&
			time.Now().After(*document.DueDate),
	}, nil
}

func (s *disciplineListDocumentService) GetById(ctx context.Context, id string) (dto.DisciplineListDocumentResponse, error) {
	disciplineListDocument, err := s.disciplineListDocumentRepository.GetByID(ctx, nil, id, "Package", "Document", "Consolidators.DisciplineGroupConsolidator.User")
	if err != nil {
		return dto.DisciplineListDocumentResponse{}, err
	}

	var consolidatorResponse []dto.DisciplineListDocumentConsolidatorResponse
	for _, c := range disciplineListDocument.Consolidators {
		mylog.Infoln(c)
		consolidatorResponse = append(consolidatorResponse, dto.DisciplineListDocumentConsolidatorResponse{
			UserID:                   c.DisciplineGroupConsolidator.User.ID.String(),
			DisciplineListDocumentID: disciplineListDocument.ID.String(),
			Name:                     c.DisciplineGroupConsolidator.User.Name,
		})
	}

	return dto.DisciplineListDocumentResponse{
		ID:      disciplineListDocument.ID.String(),
		Package: disciplineListDocument.Package.Name,
		Document: &dto.DocumentDetailResponse{
			ID:                       disciplineListDocument.DocumentID.String(),
			DocumentUrl:              disciplineListDocument.Document.DocumentUrl,
			DocumentSerialNumber:     disciplineListDocument.Document.DocumentSerialNumber,
			CTRNumber:                disciplineListDocument.Document.CTRNumber,
			WBS:                      disciplineListDocument.Document.WBS,
			CompanyDocumentNumber:    disciplineListDocument.Document.CompanyDocumentNumber,
			ContractorDocumentNumber: disciplineListDocument.Document.ContractorDocumentNumber,
			DocumentTitle:            disciplineListDocument.Document.DocumentTitle,
			Discipline:               disciplineListDocument.Document.Discipline,
			SubDiscipline:            disciplineListDocument.Document.SubDiscipline,
			DocumentType:             disciplineListDocument.Document.DocumentType,
			DocumentCategory:         disciplineListDocument.Document.DocumentCategory,
			Package:                  disciplineListDocument.Package.Name,
			DueDate:                  disciplineListDocument.Document.DueDate,
			Status:                   string(disciplineListDocument.Document.Status),
		},
		Consolidators: consolidatorResponse,
		IsDueDate: disciplineListDocument.Document.DueDate != nil &&
			time.Now().After(*disciplineListDocument.Document.DueDate),
	}, nil
}

func (s *disciplineListDocumentService) GetAll(ctx context.Context, disciplineGroupId, userId string, metaReq meta.Meta) ([]dto.DisciplineListDocumentResponse, meta.Meta, error) {
	_, _, err := s.getPackagePermission(ctx, userId)
	if err != nil {
		return nil, meta.Meta{}, err
	}

	disciplineListDocuments, metaRes, err := s.disciplineListDocumentRepository.GetAllByDisciplineGroupID(ctx, nil, disciplineGroupId, metaReq, "Package", "Document", "Consolidators.DisciplineGroupConsolidator.User")
	if err != nil {
		return nil, meta.Meta{}, err
	}

	var disciplineListDocumentResponse []dto.DisciplineListDocumentResponse
	for _, disciplineListDocument := range disciplineListDocuments {
		var consolidatorResponse []dto.DisciplineListDocumentConsolidatorResponse
		for _, c := range disciplineListDocument.Consolidators {
			consolidatorResponse = append(consolidatorResponse, dto.DisciplineListDocumentConsolidatorResponse{
				UserID:                   c.DisciplineGroupConsolidator.User.ID.String(),
				DisciplineListDocumentID: disciplineListDocument.ID.String(),
				Name:                     c.DisciplineGroupConsolidator.User.Name,
			})
		}

		disciplineListDocumentResponse = append(disciplineListDocumentResponse, dto.DisciplineListDocumentResponse{
			ID:      disciplineListDocument.ID.String(),
			Package: disciplineListDocument.Package.Name,
			Document: &dto.DocumentDetailResponse{
				ID:                       disciplineListDocument.DocumentID.String(),
				DocumentUrl:              disciplineListDocument.Document.DocumentUrl,
				DocumentSerialNumber:     disciplineListDocument.Document.DocumentSerialNumber,
				CTRNumber:                disciplineListDocument.Document.CTRNumber,
				WBS:                      disciplineListDocument.Document.WBS,
				CompanyDocumentNumber:    disciplineListDocument.Document.CompanyDocumentNumber,
				ContractorDocumentNumber: disciplineListDocument.Document.ContractorDocumentNumber,
				DocumentTitle:            disciplineListDocument.Document.DocumentTitle,
				Discipline:               disciplineListDocument.Document.Discipline,
				SubDiscipline:            disciplineListDocument.Document.SubDiscipline,
				DocumentType:             disciplineListDocument.Document.DocumentType,
				DocumentCategory:         disciplineListDocument.Document.DocumentCategory,
				Package:                  disciplineListDocument.Package.Name,
				DueDate:                  disciplineListDocument.Document.DueDate,
				Status:                   string(disciplineListDocument.Document.Status),
			},
			Consolidators: consolidatorResponse,
			IsDueDate: disciplineListDocument.Document.DueDate != nil &&
				time.Now().After(*disciplineListDocument.Document.DueDate),
		})
	}

	return disciplineListDocumentResponse, metaRes, nil
}

func (s *disciplineListDocumentService) Update(ctx context.Context, req dto.UpdateDisciplineListDocumentRequest) error {
	pkg, _, err := s.getPackagePermission(ctx, req.UserId)
	if err != nil {
		return err
	}

	disciplineListDocument, err := s.disciplineListDocumentRepository.GetByID(ctx, nil, req.ID, "Consolidators.DisciplineGroupConsolidator.User")
	if err != nil {
		return err
	}

	if pkg != nil && pkg.ID != disciplineListDocument.PackageID {
		return myerror.New("you not allowed to this package", http.StatusUnauthorized)
	}

	document, err := s.documentRepository.GetByID(ctx, nil, req.DocumentID)
	if err != nil {
		return err
	}
	disciplineListDocument.DocumentID = document.ID

	consolidatorMap := map[string]bool{}
	for _, c := range disciplineListDocument.Consolidators {
		consolidatorMap[c.DisciplineGroupConsolidatorID.String()] = true
	}

	reqConsolidatorMap := map[string]bool{}
	for _, c := range req.Consolidators {
		reqConsolidatorMap[c.DisciplineGroupConsolidatorID] = true
	}

	var deletedConsolidators []string
	for _, c := range disciplineListDocument.Consolidators {
		if _, ok := reqConsolidatorMap[c.DisciplineGroupConsolidatorID.String()]; !ok {
			deletedConsolidators = append(deletedConsolidators, c.ID.String())
		}
	}

	if err := s.disciplineListDocumentConsolidatorRepository.DeleteBulk(ctx, nil, deletedConsolidators); err != nil {
		return err
	}

	var newConsolidator []entity.DisciplineListDocumentConsolidator
	for _, c := range req.Consolidators {
		if _, ok := consolidatorMap[c.DisciplineGroupConsolidatorID]; !ok {
			newConsolidator = append(newConsolidator, entity.DisciplineListDocumentConsolidator{
				DisciplineListDocumentID:      disciplineListDocument.ID,
				DisciplineGroupConsolidatorID: uuid.MustParse(c.DisciplineGroupConsolidatorID),
			})
		}
	}

	if err := s.disciplineListDocumentConsolidatorRepository.CreateBulk(ctx, nil, newConsolidator); err != nil {
		return err
	}

	// set updated_by and updated_at
	disciplineListDocument.UpdatedBy = uuid.MustParse(req.UserId)

	if err := s.disciplineListDocumentRepository.Update(ctx, nil, disciplineListDocument); err != nil {
		return err
	}

	return nil
}

func (s *disciplineListDocumentService) Delete(ctx context.Context, userId, disciplineListDocumentId string) error {
	pkg, _, err := s.getPackagePermission(ctx, userId)
	if err != nil {
		return err
	}

	disciplineListDocument, err := s.disciplineListDocumentRepository.GetByID(ctx, nil, disciplineListDocumentId)
	if err != nil {
		return err
	}

	if pkg != nil && pkg.ID != disciplineListDocument.PackageID {
		return myerror.New("you not allowed to this package", http.StatusUnauthorized)
	}

	if err := s.commentRepository.DeleteByDisciplineListDocumentID(ctx, nil, []string{disciplineListDocument.ID.String()}); err != nil {
		return err
	}

	// mark who deleted
	disciplineListDocument.DeletedBy = uuid.MustParse(userId)
	if err = s.disciplineListDocumentRepository.Delete(ctx, nil, disciplineListDocument); err != nil {
		return err
	}

	return nil
}

func (s *disciplineListDocumentService) GenerateExcel(ctx context.Context, userId, disciplineListDocumentId string) (*bytes.Buffer, string, error) {
	dld, err := s.disciplineListDocumentRepository.GetByID(ctx, nil, disciplineListDocumentId,
		"Package",
		"DisciplineGroup",
		"Consolidators.DisciplineGroupConsolidator.User",
		"Comments.CommentReplies",
		"Comments.User",
		"Document")
	if err != nil {
		return nil, "", err
	}

	contractor, err := s.userRepository.GetContractorByPackage(ctx, nil, dld.PackageID.String(), "Package")
	if err != nil {
		return nil, "", err
	}

	// Prepare data for Excel
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

	consolidatorStr := ""
	for i, c := range dld.Consolidators {
		userName := "deleted user"
		if c.DisciplineGroupConsolidator.User != nil {
			userName = c.DisciplineGroupConsolidator.User.Name
		}
		if i > 0 {
			consolidatorStr += fmt.Sprintf("\n%d. %s", i+1, userName)
		} else {
			consolidatorStr += fmt.Sprintf("%d. %s", i+1, userName)
		}
	}

	reqData := []mypdf.GenerateRequestData{
		{
			PackageInfoData: mypdf.PackageInfoData{
				Package:           contractor.Package.Name,
				ContractorInitial: contractor.Name,
			},
			DisciplineSectionData: mypdf.DisciplineSectionData{
				Discipline:   dld.DisciplineGroup.UserDiscipline,
				Consolidator: consolidatorStr,
			},
			CommentRow: comments,
		},
	}

	excelBuffer, filename, err := mypdf.GenerateExcel(reqData)
	if err != nil {
		return nil, "", err
	}

	return excelBuffer, filename, nil
}

func (s *disciplineListDocumentService) getPackagePermission(ctx context.Context, userId string) (*entity.Package, entity.User, error) {
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
