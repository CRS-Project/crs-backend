package entity

import "github.com/google/uuid"

type DisciplineListDocumentConsolidator struct {
	ID uuid.UUID `json:"id" gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`

	DisciplineGroupConsolidatorID uuid.UUID `json:"discipline_group_consolidator_id" gorm:"not null"`
	DisciplineListDocumentID      uuid.UUID `json:"discipline_list_document_id" gorm:"not null"`

	DeletedBy uuid.UUID `json:"deleted_by" gorm:"type:uuid"`
	UpdatedBy uuid.UUID `json:"updated_by" gorm:"type:uuid"`
	Timestamp

	DisciplineGroupConsolidator *DisciplineGroupConsolidator `json:"discipline_group_consolidator,omitempty" gorm:"foreignKey:DisciplineGroupConsolidatorID"`
	DisciplineListDocument      *DisciplineListDocument      `json:"discipline_list_document,omitempty" gorm:"foreignKey:DisciplineListDocumentID"`
}
