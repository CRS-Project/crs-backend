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
