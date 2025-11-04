package entity

import (
	"github.com/google/uuid"
)

type Role string

const (
	RoleSuperAdmin Role = "SUPER ADMIN"
	RoleContractor Role = "CONTRACTOR"
	RoleReviewer   Role = "REVIEWER"
)

type User struct {
	ID         uuid.UUID `json:"id" gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	Name       string    `json:"name" gorm:"not null"`
	Email      string    `json:"email" gorm:"uniqueIndex;not null"`
	Password   string    `json:"password" gorm:"not null"`
	IsVerified bool      `json:"is_verified" gorm:"default:false;not null"`
	Role       Role      `json:"role" gorm:"default:REVIEWER;not null"`

	Initial      string  `json:"initial" gorm:"not null"`
	Institution  string  `json:"institution" gorm:"not null"`
	PhotoProfile *string `json:"photo_profile" gorm:""`

	UserDisciplineNumberID uuid.UUID `json:"user_discipline_number_id" gorm:"not null"`

	Timestamp

	UserDisciplineNumber UserDisciplineNumber `json:"user_discipline_number" gorm:"foreignKey:UserDisciplineNumberID"`
}

func (u *User) TableName() string {
	return "users"
}
