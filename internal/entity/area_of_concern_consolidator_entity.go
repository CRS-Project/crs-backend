package entity

import "github.com/google/uuid"

type AreaOfConcernConsolidator struct {
	ID uuid.UUID `json:"id" gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`

	UserID          uuid.UUID `json:"user_id" gorm:"not null"`
	AreaOfConcernID uuid.UUID `json:"area_of_concern_id" gorm:"not null"`

	Timestamp

	User          *User          `json:"user,omitempty" gorm:"foreignKey:UserID"`
	AreaOfConcern *AreaOfConcern `json:"area_of_concern,omitempty" gorm:"foreignKey:AreaOfConcernID"`
}
