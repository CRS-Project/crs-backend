package entity

import "github.com/google/uuid"

type UserDiscipline struct {
	ID      uuid.UUID `json:"id" gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	Name    string    `json:"name" gorm:"not null"`
	Initial string    `json:"initial" gorm:"not null"`

	DeletedBy uuid.UUID `json:"deleted_by" gorm:"type:uuid"`
	UpdatedBy uuid.UUID `json:"updated_by" gorm:"type:uuid"`
	Timestamp

	Users []User `json:"user,omitempty" gorm:"foreignKey:UserDisciplineID"`
}
