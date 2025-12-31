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
	"github.com/xuri/excelize/v2"

	"gorm.io/gorm"
)

type (
	DocumentService interface {
		Create(ctx context.Context, req dto.CreateDocumentRequest) (dto.DocumentDetailResponse, error)
		CreateBulk(ctx context.Context, req dto.CreateBulkDocumentRequest) ([]dto.GetAllDocumentResponse, error)
		GetAll(ctx context.Context, userId string, metaReq meta.Meta) ([]dto.GetAllDocumentResponse, meta.Meta, error)
		GetByID(ctx context.Context, documentId string) (dto.DocumentDetailResponse, error)
		Update(ctx context.Context, req dto.UpdateDocumentRequest) (dto.DocumentDetailResponse, error)
		Delete(ctx context.Context, userId, documentId string) error
	}

	documentService struct {
		documentRepository               repository.DocumentRepository
		disciplineListDocumentRepository repository.DisciplineListDocumentRepository
		packageRepository                repository.PackageRepository
		userRepository                   repository.UserRepository
		db                               *gorm.DB ``
	}
)

func NewDocument(documentRepository repository.DocumentRepository,
	disciplineListDocumentRepository repository.DisciplineListDocumentRepository,
	packageRepository repository.PackageRepository,
	userRepository repository.UserRepository,
	db *gorm.DB) DocumentService {
	return &documentService{
		documentRepository:               documentRepository,
		disciplineListDocumentRepository: disciplineListDocumentRepository,
		packageRepository:                packageRepository,
		userRepository:                   userRepository,
		db:                               db,
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
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return dto.DocumentDetailResponse{}, myerror.New("this package not have contractor, please set it first", http.StatusNotFound)
			}
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

func (s *documentService) CreateBulk(ctx context.Context, req dto.CreateBulkDocumentRequest) ([]dto.GetAllDocumentResponse, error) {
	pkg, user, err := s.getPackagePermission(ctx, req.UserID)
	if err != nil {
		return nil, err
	}

	if pkg == nil {
		pkgVal, err := s.packageRepository.GetByID(ctx, nil, req.PackageID)
		if err != nil {
			return nil, err
		}
		pkg = &pkgVal
	}

	file, err := req.FileSheet.Open()
	if err != nil {
		return nil, err
	}

	xlsx, err := excelize.OpenReader(file)
	if err != nil {
		return nil, err
	}

	if len(xlsx.GetSheetList()) == 0 {
		return nil, myerror.New("received sheet does not have a worksheet", http.StatusNotFound)
	}

	sheetName := xlsx.GetSheetList()[0]

	rows, err := xlsx.Rows(sheetName)
	if err != nil {
		return nil, err
	}

	rowCounter := 0
	var documents []entity.Document

	for rows.Next() {
		rowCounter++
		if rowCounter < 4 {
			continue
		}

		row, err := rows.Columns()
		if err != nil {
			return nil, err
		}

		documents = append(documents, entity.Document{})
		documents[len(documents)-1].Package = pkg

		for i, val := range row {
			switch i {
			case 2:
				contractor, err := s.userRepository.GetContractorByPackage(ctx, nil, pkg.ID.String())
				if err != nil {
					break
				}

				if user.Role == entity.RoleContractor && user.ID != contractor.ID {
					break
				}

				documents[len(documents)-1].Contractor = &contractor
			case 3:
				documents[len(documents)-1].DocumentSerialNumber = val
			case 4:
				documents[len(documents)-1].CTRNumber = val
			case 5:
				documents[len(documents)-1].WBS = val
			case 6:
				documents[len(documents)-1].CompanyDocumentNumber = val
			case 7:
				documents[len(documents)-1].ContractorDocumentNumber = val
			case 8:
				documents[len(documents)-1].DocumentTitle = val
			case 9:
				documents[len(documents)-1].Discipline = val
			// case 10: // Discipline Ori? what the hell is that?
			// 	documents[len(documents)-1].Discipline = val
			case 11:
				if val == "" {
					break
				}
				documents[len(documents)-1].SubDiscipline = &val
			case 12:
				documents[len(documents)-1].DocumentType = val
			case 13:
				documents[len(documents)-1].DocumentCategory = val
			case 14:
				documents[len(documents)-1].Status = entity.StatusDocument(val)
			}
		}

		validRow := documents[len(documents)-1].Contractor != nil
		validRow = validRow && documents[len(documents)-1].CompanyDocumentNumber != ""
		validRow = validRow && documents[len(documents)-1].ContractorDocumentNumber != ""
		validRow = validRow && documents[len(documents)-1].DocumentTitle != ""

		switch documents[len(documents)-1].Status {
		case entity.StatusDocumentIFR, entity.StatusDocumentIFU:
		default:
			validRow = false
		}

		if !validRow {
			documents = documents[:len(documents)-1]
		}
	}

	if err = rows.Close(); err != nil {
		return nil, err
	}

	var documentsRes []dto.GetAllDocumentResponse
	for _, document := range documents {
		document, err = s.documentRepository.Create(ctx, nil, document)
		if err != nil {
			return nil, err
		}

		documentsRes = append(documentsRes, dto.GetAllDocumentResponse{
			ID:                       document.ID.String(),
			CompanyDocumentNumber:    document.CompanyDocumentNumber,
			ContractorDocumentNumber: document.ContractorDocumentNumber,
			DocumentTitle:            document.DocumentTitle,
			DocumentType:             document.DocumentType,
			DocumentCategory:         document.DocumentCategory,
			Package:                  pkg.Name,
			Status:                   string(document.Status),
		})
	}

	if len(documentsRes) == 0 {
		return nil, myerror.New("no valid data in sheets", http.StatusBadRequest)
	}

	return documentsRes, nil
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
	document.UpdatedBy = uuid.MustParse(req.UserID)

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

	document, err := s.documentRepository.GetByID(ctx, nil, documentId, "DisciplineListDocuments")
	if err != nil {
		return err
	}

	if pkg != nil && document.PackageID != pkg.ID {
		return myerror.New("you don't have permission for this package", http.StatusUnauthorized)
	}

	if len(document.DisciplineListDocuments) > 0 {
		return myerror.New("Unable to delete the document because it is linked to a discipline list", http.StatusBadRequest)
	}

	document.DeletedBy = uuid.MustParse(userId)
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
