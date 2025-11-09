package dto

type (
	CreateDocumentRequest struct {
		UserID                   string  `json:"-"`
		PackageID                string  `json:"package_id" binding:"required"`
		DocumentUrl              *string `json:"document_url" binding:""`
		DocumentSerialNumber     string  `json:"document_serial_number" binding:"required"`
		CTRNumber                string  `json:"ctr_number" binding:"required"`
		WBS                      string  `json:"wbs" binding:"required"`
		CompanyDocumentNumber    string  `json:"company_document_number" binding:"required"`
		ContractorDocumentNumber string  `json:"contractor_document_number"`
		DocumentTitle            string  `json:"document_title" binding:"required"`
		Discipline               string  `json:"discipline" binding:"required"`
		SubDiscipline            *string `json:"sub_discipline"`
		DocumentType             string  `json:"document_type" binding:"required"`
		DocumentCategory         string  `json:"document_category" binding:"required"`
		Status                   string  `json:"status" binding:"required"`
	}

	UpdateDocumentRequest struct {
		ID                       string  `json:"_"`
		UserID                   string  `json:"-"`
		DocumentUrl              *string `json:"document_url" binding:""`
		DocumentSerialNumber     string  `json:"document_serial_number" binding:"required"`
		CTRNumber                string  `json:"ctr_number" binding:"required"`
		WBS                      string  `json:"wbs" binding:"required"`
		CompanyDocumentNumber    string  `json:"company_document_number" binding:"required"`
		ContractorDocumentNumber string  `json:"contractor_document_number" binding:"required"`
		DocumentTitle            string  `json:"document_title" binding:"required"`
		Discipline               string  `json:"discipline" binding:"required"`
		SubDiscipline            *string `json:"sub_discipline"`
		DocumentType             string  `json:"document_type" binding:"required"`
		DocumentCategory         string  `json:"document_category"  binding:"required"`
		Status                   string  `json:"status" binding:"required"`
	}

	GetAllDocumentResponse struct {
		ID                       string `json:"id"`
		CompanyDocumentNumber    string `json:"company_document_number"`
		ContractorDocumentNumber string `json:"contractor_document_number"`
		DocumentTitle            string `json:"document_title"`
		DocumentType             string `json:"document_type"`
		DocumentCategory         string `json:"document_category"`
		Package                  string `json:"package"`
		Status                   string `json:"status"`
	}

	DocumentDetailResponse struct {
		ID                       string  `json:"id"`
		DocumentUrl              *string `json:"document_url"`
		DocumentSerialNumber     string  `json:"document_serial_number"`
		CTRNumber                string  `json:"ctr_number"`
		WBS                      string  `json:"wbs"`
		CompanyDocumentNumber    string  `json:"company_document_number"`
		ContractorDocumentNumber string  `json:"contractor_document_number"`
		DocumentTitle            string  `json:"document_title"`
		Discipline               string  `json:"discipline"`
		SubDiscipline            *string `json:"sub_discipline"`
		DocumentType             string  `json:"document_type"`
		DocumentCategory         string  `json:"document_category"`
		Package                  string  `json:"package"`
		Status                   string  `json:"status"`
	}
)
