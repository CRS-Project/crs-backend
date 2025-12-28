package entity

import "github.com/google/uuid"

type DisciplineGroup struct {
	ID                uuid.UUID `json:"id" gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	ReviewFocus       string    `json:"review_focus" gorm:"not null"`
	UserDiscipline    string    `json:"user_discipline" gorm:""`
	DisciplineInitial string    `json:"discipline_initial" gorm:""`

	PackageID uuid.UUID `json:"package_id" gorm:"not null"`

	DeletedBy uuid.UUID `json:"deleted_by" gorm:"type:uuid"`
	UpdatedBy uuid.UUID `json:"updated_by" gorm:"type:uuid"`
	Timestamp

	Package                      *Package                      `json:"package,omitempty" gorm:"foreignKey:PackageID"`
	DisciplineListDocuments      []DisciplineListDocument      `json:"discipline_list_documents,omitempty" gorm:"foreignKey:DisciplineGroupID"`
	DisciplineGroupConsolidators []DisciplineGroupConsolidator `json:"discipline_group_consolidator,omitempty" gorm:"foreignKey:DisciplineGroupID"`
}
