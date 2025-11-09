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
	DocumentService interface {
		Create(ctx context.Context, req dto.CreateDocumentRequest) (dto.DocumentDetailResponse, error)
		GetAll(ctx context.Context, userId string, metaReq meta.Meta) ([]dto.GetAllDocumentResponse, meta.Meta, error)
		GetByID(ctx context.Context, documentId string) (dto.DocumentDetailResponse, error)
		Update(ctx context.Context, req dto.UpdateDocumentRequest) (dto.DocumentDetailResponse, error)
		Delete(ctx context.Context, userId, documentId string) error
	}

	documentService struct {
		documentRepository repository.DocumentRepository
		packageRepository  repository.PackageRepository
		userRepository     repository.UserRepository
		db                 *gorm.DB ``
	}
)

func NewDocument(documentRepository repository.DocumentRepository,
	packageRepository repository.PackageRepository,
	userRepository repository.UserRepository,
	db *gorm.DB) DocumentService {
	return &documentService{
		documentRepository: documentRepository,
		packageRepository:  packageRepository,
		userRepository:     userRepository,
		db:                 db,
	}
}

func (s *documentService) Create(ctx context.Context, req dto.CreateDocumentRequest) (dto.DocumentDetailResponse, error) {
	pkg, user, err := s.getPackagePermission(ctx, req.UserID)
	if err != nil {
		return dto.DocumentDetailResponse{}, err
	}

	var contractor entity.User
	if pkg == nil {
		contractor, err = s.userRepository.GetContractorByPackage(ctx, nil, req.PackageID, "Package")
		if err != nil {
			return dto.DocumentDetailResponse{}, err
		}

		pkg = contractor.Package
	} else {
		contractor = user
	}

	documentResult, err := s.documentRepository.Create(ctx, nil, entity.Document{
		ContractorID:             contractor.ID,
		DocumentUrl:              req.DocumentUrl,
		PackageID:                pkg.ID,
		DocumentSerialNumber:     req.DocumentSerialNumber,
		CTRNumber:                req.CTRNumber,
		WBS:                      req.WBS,
		CompanyDocumentNumber:    req.CompanyDocumentNumber,
		ContractorDocumentNumber: req.ContractorDocumentNumber,
		DocumentTitle:            req.DocumentTitle,
		Discipline:               req.Discipline,
		SubDiscipline:            req.SubDiscipline,
		DocumentType:             req.DocumentType,
		DocumentCategory:         req.DocumentCategory,
		Status:                   entity.StatusDocument(req.Status),
	})
	if err != nil {
		return dto.DocumentDetailResponse{}, err
	}

	return dto.DocumentDetailResponse{
		ID:                       documentResult.ID.String(),
		DocumentUrl:              documentResult.DocumentUrl,
		DocumentSerialNumber:     documentResult.DocumentSerialNumber,
		CTRNumber:                documentResult.CTRNumber,
		WBS:                      documentResult.WBS,
		CompanyDocumentNumber:    documentResult.CompanyDocumentNumber,
		ContractorDocumentNumber: documentResult.ContractorDocumentNumber,
		DocumentTitle:            documentResult.DocumentTitle,
		Discipline:               documentResult.Discipline,
		SubDiscipline:            documentResult.SubDiscipline,
		DocumentType:             documentResult.DocumentType,
		DocumentCategory:         documentResult.DocumentCategory,
		Package:                  pkg.Name,
		Status:                   string(documentResult.Status),
	}, nil
}

func (s *documentService) GetAll(ctx context.Context, userId string, metaReq meta.Meta) ([]dto.GetAllDocumentResponse, meta.Meta, error) {
	pkg, _, err := s.getPackagePermission(ctx, userId)
	if err != nil {
		return nil, meta.Meta{}, err
	}

	pkgId := ""
	if pkg != nil {
		pkgId = pkg.ID.String()
	}

	documents, metaRes, err := s.documentRepository.GetAll(ctx, nil, pkgId, metaReq, "Contractor", "Package")
	if err != nil {
		return nil, meta.Meta{}, err
	}

	var getDocuments []dto.GetAllDocumentResponse
	for _, document := range documents {
		getDocuments = append(getDocuments, dto.GetAllDocumentResponse{
			ID:                       document.ID.String(),
			CompanyDocumentNumber:    document.CompanyDocumentNumber,
			ContractorDocumentNumber: document.ContractorDocumentNumber,
			DocumentTitle:            document.DocumentTitle,
			DocumentType:             document.DocumentType,
			DocumentCategory:         document.DocumentCategory,
			Package:                  document.Package.Name,
			Status:                   string(document.Status),
		})
	}

	return getDocuments, metaRes, nil
}

func (s *documentService) GetByID(ctx context.Context, documentId string) (dto.DocumentDetailResponse, error) {
	document, err := s.documentRepository.GetByID(ctx, nil, documentId, "Contractor", "Package")
	if err != nil {
		return dto.DocumentDetailResponse{}, err
	}

	return dto.DocumentDetailResponse{
		ID:                       document.ID.String(),
		DocumentUrl:              document.DocumentUrl,
		DocumentSerialNumber:     document.DocumentSerialNumber,
		CTRNumber:                document.CTRNumber,
		WBS:                      document.WBS,
		CompanyDocumentNumber:    document.CompanyDocumentNumber,
		ContractorDocumentNumber: document.ContractorDocumentNumber,
		DocumentTitle:            document.DocumentTitle,
		Discipline:               document.Discipline,
		SubDiscipline:            document.SubDiscipline,
		DocumentType:             document.DocumentType,
		DocumentCategory:         document.DocumentCategory,
		Package:                  document.Package.Name,
		Status:                   string(document.Status),
	}, nil
}

func (s *documentService) Update(ctx context.Context, req dto.UpdateDocumentRequest) (dto.DocumentDetailResponse, error) {
	pkg, _, err := s.getPackagePermission(ctx, req.UserID)
	if err != nil {
		return dto.DocumentDetailResponse{}, err
	}

	document, err := s.documentRepository.GetByID(ctx, nil, req.ID, "Contractor", "Package")
	if err != nil {
		return dto.DocumentDetailResponse{}, err
	}

	if pkg != nil && document.PackageID != pkg.ID {
		return dto.DocumentDetailResponse{}, myerror.New("you don't have permission for this package", http.StatusUnauthorized)
	}

	document.DocumentUrl = req.DocumentUrl
	document.DocumentSerialNumber = req.DocumentSerialNumber
	document.CTRNumber = req.CTRNumber
	document.WBS = req.WBS
	document.CompanyDocumentNumber = req.CompanyDocumentNumber
	document.ContractorDocumentNumber = req.ContractorDocumentNumber
	document.DocumentTitle = req.DocumentTitle
	document.Discipline = req.Discipline
	document.SubDiscipline = req.SubDiscipline
	document.DocumentType = req.DocumentType
	document.DocumentCategory = req.DocumentCategory
	document.Status = entity.StatusDocument(req.Status)

	document, err = s.documentRepository.Update(ctx, nil, document)
	if err != nil {
		return dto.DocumentDetailResponse{}, err
	}

	return dto.DocumentDetailResponse{
		ID:                       document.ID.String(),
		DocumentUrl:              document.DocumentUrl,
		DocumentSerialNumber:     document.DocumentSerialNumber,
		CTRNumber:                document.CTRNumber,
		WBS:                      document.WBS,
		CompanyDocumentNumber:    document.CompanyDocumentNumber,
		ContractorDocumentNumber: document.ContractorDocumentNumber,
		DocumentTitle:            document.DocumentTitle,
		Discipline:               document.Discipline,
		SubDiscipline:            document.SubDiscipline,
		DocumentType:             document.DocumentType,
		DocumentCategory:         document.DocumentCategory,
		Package:                  document.Package.Name,
		Status:                   string(document.Status),
	}, nil
}

func (s *documentService) Delete(ctx context.Context, userId, documentId string) error {
	pkg, _, err := s.getPackagePermission(ctx, userId)
	if err != nil {
		return err
	}

	document, err := s.documentRepository.GetByID(ctx, nil, documentId)
	if err != nil {
		return err
	}

	if pkg != nil && document.PackageID != pkg.ID {
		return myerror.New("you don't have permission for this package", http.StatusUnauthorized)
	}

	if err = s.documentRepository.Delete(ctx, nil, document); err != nil {
		return err
	}

	return nil
}

func (s *documentService) getPackagePermission(ctx context.Context, userId string) (*entity.Package, entity.User, error) {
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
