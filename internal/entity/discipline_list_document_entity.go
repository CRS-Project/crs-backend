package entity

import "github.com/google/uuid"

type DisciplineListDocument struct {
	ID uuid.UUID `json:"id" gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`

	DocumentID        uuid.UUID `json:"document_id" gorm:"not null"`
	DisciplineGroupID uuid.UUID `json:"discipline_group_id" gorm:"not null"`
	PackageID         uuid.UUID `json:"package_id" gorm:"not null"`

	DeletedBy uuid.UUID `json:"deleted_by" gorm:"type:uuid"`
	UpdatedBy uuid.UUID `json:"updated_by" gorm:"type:uuid"`
	Timestamp

	DisciplineGroup *DisciplineGroup                     `json:"discipline_group,omitempty" gorm:"foreignKey:DisciplineGroupID"`
	Package         *Package                             `json:"package,omitempty" gorm:"foreignKey:PackageID"`
	Document        *Document                            `json:"document,omitempty" gorm:"foreignKey:DocumentID"`
	Consolidators   []DisciplineListDocumentConsolidator `json:"consolidators,omitempty" gorm:"foreignKey:DisciplineListDocumentID"`
	Comments        []Comment                            `json:"comments,omitempty" gorm:"foreignKey:DisciplineListDocumentID"`
}
