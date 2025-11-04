package entity

import "github.com/google/uuid"

type UserDisciplineNumber struct {
	ID     uuid.UUID `json:"id" gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	Number int       `json:"number" gorm:""`

	UserDisciplineID uuid.UUID  `json:"user_discipline_id" gorm:"not null"`
	PackageID        *uuid.UUID `json:"package_id" gorm:""`

	Timestamp

	UserDiscipline *UserDiscipline `json:"user_discipline,omitempty" gorm:"foreignKey:UserDisciplineID"`
	Package        *Package        `json:"package,omitempty" gorm:"foreignKey:PackageID"`
}
