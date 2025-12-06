package entity

import (
	"github.com/CRS-Project/crs-backend/internal/dto"
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
	Email      string    `json:"email" gorm:"not null"`
	Password   string    `json:"password" gorm:"not null"`
	IsVerified bool      `json:"is_verified" gorm:"default:false;not null"`
	Role       Role      `json:"role" gorm:"default:REVIEWER;not null"`

	Initial          string  `json:"initial" gorm:"not null"`
	Institution      string  `json:"institution" gorm:"not null"`
	PhotoProfile     *string `json:"photo_profile" gorm:""`
	DisciplineNumber int     `json:"discipline_number" gorm:"not null"`

	UserDisciplineID uuid.UUID  `json:"user_discipline_id" gorm:"not null"`
	PackageID        *uuid.UUID `json:"package_id" gorm:""` // super admin akan null (berarti punya semua akses)

	Timestamp

	UserDiscipline *UserDiscipline `json:"user_discipline,omitempty" gorm:"foreignKey:UserDisciplineID"`
	Package        *Package        `json:"package,omitempty" gorm:"foreignKey:PackageID"`
}

func (u *User) TableName() string {
	return "users"
}

func (u *User) ToInfo() dto.PersonalInfo {
	return dto.PersonalInfo{
		ID:           u.ID.String(),
		Name:         u.Name,
		Email:        u.Email,
		Initial:      u.Initial,
		Institution:  u.Institution,
		PhotoProfile: u.PhotoProfile,
		Role:         string(u.Role),
	}
}
