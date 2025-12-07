package entity

import "github.com/google/uuid"

type DisciplineGroupConsolidator struct {
	ID uuid.UUID `json:"id" gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`

	UserID            uuid.UUID `json:"user_id" gorm:"not null"`
	DisciplineGroupID uuid.UUID `json:"discipline_group_id" gorm:"not null"`

	Timestamp

	User                                *User                                `json:"user,omitempty" gorm:"foreignKey:UserID"`
	DisciplineGroup                     *DisciplineGroup                     `json:"discipline_group,omitempty" gorm:"foreignKey:DisciplineGroupID"`
	DisciplineListDocumentConsolidators []DisciplineListDocumentConsolidator `json:"discipline_list_document_consolidators,omitempty" gorm:"foreignKey:DisciplineGroupConsolidatorID"`
}
