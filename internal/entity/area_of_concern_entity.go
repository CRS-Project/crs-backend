package entity

import "github.com/google/uuid"

type AreaOfConcern struct {
	ID              uuid.UUID `json:"id" gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	Description     string    `json:"description" gorm:"not null"`
	AreaOfConcernId string    `json:"area_of_concern_id" gorm:"not null"`

	AreaOfConcernGroupID uuid.UUID `json:"area_of_concern_group_id" gorm:"not null"`
	PackageID            uuid.UUID `json:"package_id" gorm:"not null"`

	Timestamp

	AreaOfConcernGroup *AreaOfConcernGroup         `json:"area_of_concern_group,omitempty" gorm:"foreignKey:AreaOfConcernGroupID"`
	Package            *Package                    `json:"package,omitempty" gorm:"foreignKey:PackageID"`
	Consolidators      []AreaOfConcernConsolidator `json:"consolidators,omitempty" gorm:"foreignKey:AreaOfConcernID"`
}
