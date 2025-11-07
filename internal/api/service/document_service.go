package service

import (
	"context"
	"time"

	"github.com/CRS-Project/crs-backend/internal/api/repository"
	"github.com/CRS-Project/crs-backend/internal/dto"
	"github.com/CRS-Project/crs-backend/internal/entity"
	"github.com/CRS-Project/crs-backend/internal/pkg/meta"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type (
	DocumentService interface {
		Create(ctx context.Context, req dto.CreateDocumentRequest) (dto.GetDocument, error)
		GetAll(ctx context.Context, metaReq meta.Meta) ([]dto.GetDocument, meta.Meta, error)
		Delete(ctx context.Context, id string) error
		GetByID(ctx context.Context, id string) (dto.GetDocument, error)
		Update(ctx context.Context, req dto.UpdateDocumentRequest) (dto.GetDocument, error)
	}

	documentService struct {
		documentRepository repository.DocumentRepository
		db                 *gorm.DB
	}
)

func NewDocument(documentRepository repository.DocumentRepository, db *gorm.DB) DocumentService {
	return &documentService{
		documentRepository: documentRepository,
		db:                 db,
	}
}

func (s *documentService) Create(ctx context.Context, req dto.CreateDocumentRequest) (dto.GetDocument, error) {
	deadlineDate, err := time.Parse("15.04 • 02 Jan 2006", req.Deadline)
	if err != nil {
		return dto.GetDocument{}, err
	}

	contractorUUID, err := uuid.Parse(req.ContractorID)
	if err != nil {
		return dto.GetDocument{}, err
	}

	packageUUID, err := uuid.Parse(req.PackageID)
	if err != nil {
		return dto.GetDocument{}, err
	}

	documentCreation := entity.Document{
		ContractorID:                    contractorUUID,
		PackageID:                       packageUUID,
		DocumentSerialDisciplineNumber:  req.DocumentSerialDisciplineNumber,
		CTRDisciplineNumber:             req.CTRDisciplineNumber,
		WBS:                             req.WBS,
		CompanyDocumentDisciplineNumber: req.CompanyDocumentDisciplineNumber,
		ContractorDocumentNumber:        req.ContractorDocumentNumber,
		DocumentTitle:                   req.DocumentTitle,
		Discipline:                      req.Discipline,
		SubDiscipline:                   req.SubDiscipline,
		DocumentType:                    req.DocumentType,
		DocumentCategory:                req.DocumentCategory,
		Deadline:                        deadlineDate,
	}

	documentRes, err := s.documentRepository.Create(ctx, nil, documentCreation, "Contractor", "Package")
	if err != nil {
		return dto.GetDocument{}, err
	}

	return documentRes.GetDetail(), nil
}

func (s *documentService) GetAll(ctx context.Context, metaReq meta.Meta) ([]dto.GetDocument, meta.Meta, error) {
	documents, metaRes, err := s.documentRepository.GetAll(ctx, nil, metaReq, "Contractor", "Package")
	if err != nil {
		return nil, meta.Meta{}, err
	}

	var getDocuments []dto.GetDocument
	for _, document := range documents {
		getDocuments = append(getDocuments, document.GetDetail())
	}

	return getDocuments, metaRes, nil
}

func (s *documentService) Delete(ctx context.Context, id string) error {
	document, err := s.documentRepository.GetByID(ctx, nil, id)
	if err != nil {
		return err
	}

	if err = s.documentRepository.Delete(ctx, nil, document); err != nil {
		return err
	}

	return nil
}

func (s *documentService) GetByID(ctx context.Context, id string) (dto.GetDocument, error) {
	document, err := s.documentRepository.GetByID(ctx, nil, id, "Contractor", "Package")
	if err != nil {
		return dto.GetDocument{}, err
	}

	return document.GetDetail(), nil
}

func (s *documentService) Update(ctx context.Context, req dto.UpdateDocumentRequest) (dto.GetDocument, error) {
	document, err := s.documentRepository.GetByID(ctx, nil, req.ID, "Contractor", "Package")
	if err != nil {
		return dto.GetDocument{}, err
	}

	if req.DocumentUrl != nil {
		document.DocumentUrl = req.DocumentUrl
	}
	if req.DocumentSerialDisciplineNumber != nil {
		document.DocumentSerialDisciplineNumber = *req.DocumentSerialDisciplineNumber
	}
	if req.CTRDisciplineNumber != nil {
		document.CTRDisciplineNumber = *req.CTRDisciplineNumber
	}
	if req.WBS != nil {
		document.WBS = *req.WBS
	}
	if req.CompanyDocumentDisciplineNumber != nil {
		document.CompanyDocumentDisciplineNumber = *req.CompanyDocumentDisciplineNumber
	}
	if req.ContractorDocumentNumber != nil {
		document.ContractorDocumentNumber = *req.ContractorDocumentNumber
	}
	if req.DocumentTitle != nil {
		document.DocumentTitle = *req.DocumentTitle
	}
	if req.Discipline != nil {
		document.Discipline = *req.Discipline
	}
	if req.SubDiscipline != nil {
		document.SubDiscipline = req.SubDiscipline
	}
	if req.DocumentType != nil {
		document.DocumentType = *req.DocumentType
	}
	if req.DocumentCategory != nil {
		document.DocumentCategory = *req.DocumentCategory
	}
	if req.Deadline != nil {
		deadline, err := time.Parse("15.04 • 02 Jan 2006", *req.Deadline)
		if err != nil {
			return dto.GetDocument{}, err
		}

		document.Deadline = deadline
	}

	document, err = s.documentRepository.Update(ctx, nil, document)
	if err != nil {
		return dto.GetDocument{}, err
	}

	return document.GetDetail(), nil
}
