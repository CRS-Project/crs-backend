package entity

import "github.com/google/uuid"

type UserDicipline struct {
	ID     uuid.UUID `json:"id" gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	Name   string    `json:"name" gorm:"not null"`
	Number int       `json:"number" gorm:""`

	Initial string `json:"initial" gorm:"not null"`

	PackageID uuid.UUID `json:"package_id" gorm:"not null"`

	Timestamp
}
