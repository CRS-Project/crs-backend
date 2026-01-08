package service

import (
	"context"
	"errors"
	"net/http"

	"github.com/CRS-Project/crs-backend/internal/api/repository"
	"github.com/CRS-Project/crs-backend/internal/dto"
	"github.com/CRS-Project/crs-backend/internal/entity"
	myerror "github.com/CRS-Project/crs-backend/internal/pkg/error"
	mylog "github.com/CRS-Project/crs-backend/internal/pkg/logger"
	"github.com/CRS-Project/crs-backend/internal/pkg/meta"
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
			Status:                   string(disciplineListDocument.Document.Status),
		},
		Consolidators: consolidatorResponse,
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
				Status:                   string(disciplineListDocument.Document.Status),
			},
			Consolidators: consolidatorResponse,
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
