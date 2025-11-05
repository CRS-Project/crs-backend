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
	deadlineDate, err := time.Parse(time.RFC822, req.Deadline)
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

	return dto.GetDocument{
		DocumentInfo: dto.DocumentInfo{
			ID:                              documentRes.ID.String(),
			DocumentSerialDisciplineNumber:  documentRes.DocumentSerialDisciplineNumber,
			CTRDisciplineNumber:             documentRes.CTRDisciplineNumber,
			WBS:                             documentRes.WBS,
			CompanyDocumentDisciplineNumber: documentRes.CompanyDocumentDisciplineNumber,
			ContractorDocumentNumber:        documentRes.ContractorDocumentNumber,
			DocumentTitle:                   documentRes.DocumentTitle,
			Discipline:                      documentRes.Discipline,
			SubDiscipline:                   documentRes.SubDiscipline,
			DocumentType:                    documentRes.DocumentType,
			DocumentCategory:                documentRes.DocumentCategory,
			Deadline:                        documentRes.Deadline.Format(time.RFC822),
		},
		PackageInfo: dto.PackageInfo{
			ID:   documentRes.Package.ID.String(),
			Name: documentRes.Package.Name,
		},
		ContractorInfo: dto.PersonalInfo{
			ID:           documentRes.Contractor.ID.String(),
			Name:         documentRes.Contractor.Name,
			Email:        documentRes.Contractor.Email,
			Initial:      documentRes.Contractor.Initial,
			Institution:  documentRes.Contractor.Institution,
			PhotoProfile: documentRes.Contractor.PhotoProfile,
			Role:         string(documentRes.Contractor.Role),
		},
	}, nil
}

func (s *documentService) GetAll(ctx context.Context, metaReq meta.Meta) ([]dto.GetDocument, meta.Meta, error) {
	documents, metaRes, err := s.documentRepository.GetAll(ctx, nil, metaReq, "Contractor", "Package")
	if err != nil {
		return nil, meta.Meta{}, err
	}

	var getDocuments []dto.GetDocument
	for _, document := range documents {
		getDocuments = append(getDocuments, dto.GetDocument{
			DocumentInfo: dto.DocumentInfo{
				ID:                              document.ID.String(),
				DocumentUrl:                     document.DocumentUrl,
				DocumentSerialDisciplineNumber:  document.DocumentSerialDisciplineNumber,
				CTRDisciplineNumber:             document.CTRDisciplineNumber,
				WBS:                             document.WBS,
				CompanyDocumentDisciplineNumber: document.CompanyDocumentDisciplineNumber,
				DocumentTitle:                   document.DocumentTitle,
				Discipline:                      document.Discipline,
				SubDiscipline:                   document.SubDiscipline,
				DocumentType:                    document.DocumentType,
				DocumentCategory:                document.DocumentCategory,
				Deadline:                        document.Deadline.Format(time.RFC822),
			},
			PackageInfo: dto.PackageInfo{
				ID:   document.Package.ID.String(),
				Name: document.Package.Name,
			},
			ContractorInfo: dto.PersonalInfo{
				ID:           document.Contractor.ID.String(),
				Name:         document.Contractor.Name,
				Email:        document.Contractor.Email,
				Initial:      document.Contractor.Initial,
				Institution:  document.Contractor.Institution,
				PhotoProfile: document.Contractor.PhotoProfile,
				Role:         string(document.Contractor.Role),
			},
		})
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
