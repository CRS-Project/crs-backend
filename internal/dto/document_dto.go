package dto

type (
	CreateDocumentRequest struct {
		ContractorID             string  `json:"contractor_id" binding:"required"`
		PackageID                string  `json:"package_id" binding:"required"`
		DocumentSerialNumber     string  `json:"document_serial_number" binding:"required"`
		CTRNumber                string  `json:"ctr_number" binding:"required"`
		WBS                      string  `json:"wbs" binding:"required"`
		CompanyDocumentNumber    *string `json:"company_document_number"`
		ContractorDocumentNumber string  `json:"contractor_document_number"`
		DocumentTitle            string  `json:"document_title" binding:"required"`
		Discipline               string  `json:"discipline" binding:"required"`
		SubDiscipline            *string `json:"sub_discipline"`
		DocumentType             string  `json:"document_type" binding:"required"`
		DocumentCategory         string  `json:"document_category" binding:"required"`
		Deadline                 string  `json:"deadline" binding:"required"`
	}

	GetDocument struct {
		DocumentInfo   DocumentInfo `json:"document_info"`
		ContractorInfo PersonalInfo `json:"contractor_info"`
		PackageInfo    PackageInfo  `json:"package_info"`
	}

	DocumentInfo struct {
		DocumentSerialNumber     string  `json:"document_serial_number"`
		CTRNumber                string  `json:"ctr_number"`
		WBS                      string  `json:"wbs"`
		CompanyDocumentNumber    *string `json:"company_document_number"`
		ContractorDocumentNumber string  `json:"contractor_document_number"`
		DocumentTitle            string  `json:"document_title"`
		Discipline               string  `json:"discipline"`
		SubDiscipline            *string `json:"sub_discipline"`
		DocumentType             string  `json:"document_type"`
		DocumentCategory         string  `json:"document_category"`
		Deadline                 string  `json:"deadline"`
		Status                   string  `json:"status"`
	}
)
