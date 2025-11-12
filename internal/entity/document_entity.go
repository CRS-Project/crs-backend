package entity

import (
	"github.com/google/uuid"
)

type StatusDocument string

const (
	StatusDocumentIFR StatusDocument = "IFR Comment"
	StatusDocumentIFU StatusDocument = "IFR Comment"
)

type Document struct {
	ID                       uuid.UUID      `json:"id" gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	DocumentUrl              *string        `json:"document_url" gorm:""`
	DocumentSerialNumber     string         `json:"document_serial_number" gorm:""`
	CTRNumber                string         `json:"ctr_number" gorm:""`
	WBS                      string         `json:"wbs" gorm:""`
	CompanyDocumentNumber    string         `json:"company_document_number" gorm:""`
	ContractorDocumentNumber string         `json:"contractor_document_number" gorm:""`
	DocumentTitle            string         `json:"document_title" gorm:"not null"`
	Discipline               string         `json:"discipline" gorm:""`
	SubDiscipline            *string        `json:"sub_discipline" gorm:""`
	DocumentType             string         `json:"document_type" gorm:""`
	DocumentCategory         string         `json:"document_category" gorm:""`
	Status                   StatusDocument `json:"status" gorm:"not null"`

	ContractorID uuid.UUID `json:"contractor_id" gorm:"not null"`
	PackageID    uuid.UUID `json:"package_id" gorm:"not null"`

	Timestamp

	Contractor *User    `json:"contractor,omitempty" gorm:"foreignKey:ContractorID"`
	Package    *Package `json:"package,omitempty" gorm:"foreignKey:PackageID"`
}
