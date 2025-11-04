package entity

import "github.com/google/uuid"

type UserDiscipline struct {
	ID      uuid.UUID `json:"id" gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	Name    string    `json:"name" gorm:"not null"`
	Initial string    `json:"initial" gorm:"not null"`

	PackageID *uuid.UUID `json:"package_id" gorm:""`

	Timestamp

	Package *Package `json:"package,omitempty" gorm:"foreignKey:PackageID"`
}
