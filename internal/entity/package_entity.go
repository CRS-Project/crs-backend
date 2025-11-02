package entity

import "github.com/google/uuid"

type Package struct {
	ID   uuid.UUID `json:"id" gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	Name string    `json:"name" gorm:"not null"`

	Timestamp
}
