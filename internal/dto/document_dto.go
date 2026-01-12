package dto

import (
	"mime/multipart"
	"time"
)

type (
	CreateDocumentRequest struct {
		UserID                   string     `json:"-"`
		PackageID                string     `json:"package_id" binding:""`
		DocumentUrl              *string    `json:"document_url" binding:""`
		DocumentSerialNumber     string     `json:"document_serial_number" binding:""`
		CTRNumber                string     `json:"ctr_number" binding:""`
		WBS                      string     `json:"wbs" binding:""`
		CompanyDocumentNumber    string     `json:"company_document_number" binding:""`
		ContractorDocumentNumber string     `json:"contractor_document_number"`
		DocumentTitle            string     `json:"document_title" binding:""`
		Discipline               string     `json:"discipline" binding:""`
		SubDiscipline            *string    `json:"sub_discipline"`
		DocumentType             string     `json:"document_type" binding:""`
		DocumentCategory         string     `json:"document_category" binding:""`
		DueDate                  *time.Time `json:"due_date" binding:""`
		Status                   string     `json:"status" binding:""`
	}

	UpdateDocumentRequest struct {
		ID                       string     `json:"_"`
		UserID                   string     `json:"-"`
		DocumentUrl              *string    `json:"document_url" binding:""`
		DocumentSerialNumber     string     `json:"document_serial_number" binding:""`
		CTRNumber                string     `json:"ctr_number" binding:""`
		WBS                      string     `json:"wbs" binding:""`
		CompanyDocumentNumber    string     `json:"company_document_number" binding:""`
		ContractorDocumentNumber string     `json:"contractor_document_number" binding:""`
		DocumentTitle            string     `json:"document_title" binding:""`
		Discipline               string     `json:"discipline" binding:""`
		SubDiscipline            *string    `json:"sub_discipline"`
		DocumentType             string     `json:"document_type" binding:""`
		DocumentCategory         string     `json:"document_category"  binding:""`
		DueDate                  *time.Time `json:"due_date" binding:""`
		Status                   string     `json:"status" binding:""`
	}

	CreateBulkDocumentRequest struct {
		UserID    string                `json:"_"`
		PackageID string                `json:"-"`
		FileSheet *multipart.FileHeader `multipart.FileHeader:"file_sheet" binding:"required"`
	}

	GetAllDocumentResponse struct {
		ID                       string     `json:"id"`
		CompanyDocumentNumber    string     `json:"company_document_number"`
		ContractorDocumentNumber string     `json:"contractor_document_number"`
		DocumentTitle            string     `json:"document_title"`
		DocumentType             string     `json:"document_type"`
		DocumentCategory         string     `json:"document_category"`
		Package                  string     `json:"package"`
		DueDate                  *time.Time `json:"due_date"`
		Status                   string     `json:"status"`
		TotalComments            int        `json:"total_comment"`
	}

	DocumentDetailResponse struct {
		ID                       string     `json:"id"`
		DocumentUrl              *string    `json:"document_url"`
		DocumentSerialNumber     string     `json:"document_serial_number"`
		CTRNumber                string     `json:"ctr_number"`
		WBS                      string     `json:"wbs"`
		CompanyDocumentNumber    string     `json:"company_document_number"`
		ContractorDocumentNumber string     `json:"contractor_document_number"`
		DocumentTitle            string     `json:"document_title"`
		Discipline               string     `json:"discipline"`
		SubDiscipline            *string    `json:"sub_discipline"`
		DocumentType             string     `json:"document_type"`
		DocumentCategory         string     `json:"document_category"`
		Package                  string     `json:"package"`
		DueDate                  *time.Time `json:"due_date"`
		Status                   string     `json:"status"`
	}
)
