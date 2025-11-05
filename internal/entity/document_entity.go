package entity

import (
	"time"

	"github.com/google/uuid"
)

type Document struct {
	ID                              uuid.UUID `json:"id" gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	DocumentUrl                     *string   `json:"document_url" gorm:""`
	DocumentSerialDisciplineNumber  string    `json:"document_serial_number" gorm:"not null"`
	CTRDisciplineNumber             string    `json:"ctr_number" gorm:"not null"`
	WBS                             string    `json:"wbs" gorm:"not null"`
	CompanyDocumentDisciplineNumber string    `json:"company_document_number" gorm:"not null"`
	ContractorDocumentNumber        string    `json:"contractor_document_number" gorm:"not null"`
	DocumentTitle                   string    `json:"document_title" gorm:"not null"`
	Discipline                      string    `json:"discipline" gorm:"not null"`
	SubDiscipline                   *string   `json:"sub_discipline" gorm:""`
	DocumentType                    string    `json:"document_type" gorm:"not null"`
	DocumentCategory                string    `json:"document_category" gorm:"not null"`

	Deadline time.Time `json:"deadline" gorm:"not null"`

	ContractorID uuid.UUID `json:"contractor_id" gorm:"not null"`
	PackageID    uuid.UUID `json:"package_id" gorm:"not null"`

	Timestamp

	Contractor *User    `json:"contractor,omitempty" gorm:"foreignKey:ContractorID"`
	Package    *Package `json:"package,omitempty" gorm:"foreignKey:PackageID"`
}
