package entity

import (
	"time"

	"github.com/google/uuid"
)

type Document struct {
	ID             uuid.UUID `json:"id" gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	FilePath       *string   `json:"file_path" gorm:""`
	Discipline     string    `json:"discipline" gorm:"not null"`
	DocumentNumber string    `json:"document_number" gorm:"not null"`
	DocumentTitle  string    `json:"document_title" gorm:"not null"`
	DocumentDate   string    `json:"document_date" gorm:"not null"`

	Project string `json:"project" gorm:"not null"`
	WP      string `json:"wp" gorm:"not null"`
	Lookup  string `json:"lookup" gorm:"not null"`
	Status  string `json:"status" gorm:"not null"`

	Deadline time.Time `json:"deadline" gorm:"not null"`

	PackageID uuid.UUID `json:"package_id" gorm:"not null"`

	Timestamp

	Package *Package `json:"package,omitempty" gorm:"foreignKey:PackageID"`
}
