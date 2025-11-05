package entity

import (
	"time"

	"github.com/google/uuid"
)

type Document struct {
	ID                       uuid.UUID `json:"id" gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	DocumentUrl              *string   `json:"document_url" gorm:""`
	WBS                      string    `json:"wp" gorm:"not null"`
	DocumentSerialNumber     string    `json:"document_serial_number" gorm:"not null"`
	CompanyDocumentNumber    string    `json:""`
	ContractorDocumentNumber string    `json:"document_number" gorm:"not null"`
	DocumentTitle            string    `json:"document_title" gorm:"not null"`
	Discipline               string    `json:"discipline" gorm:"not null"`
	SubDiscipline            *string   `json:"sub_discipline" gorm:""`
	DocumentType             string    `json:"document_type" gorm:"not null"`
	DocumentCategory         string    `json:"document_category" gorm:"not null"`
	DocumentDate             string    `json:"document_date" gorm:"not null"`

	Project string `json:"project" gorm:"not null"`
	Lookup  string `json:"lookup" gorm:"not null"`
	Status  string `json:"status" gorm:"not null"`

	Deadline time.Time `json:"deadline" gorm:"not null"`

	UserID    uuid.UUID `json:"user_id" gorm:"not null"`
	PackageID uuid.UUID `json:"package_id" gorm:"not null"`

	Timestamp

	User    *User    `json:"user,omitempty" gorm:"foreignKey:UserID"`
	Package *Package `json:"package,omitempty" gorm:"foreignKey:PackageID"`
}
