package entity

import "github.com/google/uuid"

type AreaOfConcernGroup struct {
	ID             uuid.UUID `json:"id" gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	ReviewFocus    string    `json:"review_focus" gorm:"not null"`
	UserDiscipline string    `json:"user_discipline" gorm:""`

	PackageID uuid.UUID `json:"package_id" gorm:"not null"`

	Timestamp

	// UserDiscipline *UserDiscipline `json:"user_discipline,omitempty" gorm:"foreignKey:UserDisciplineID"`
	Package        *Package        `json:"package,omitempty" gorm:"foreignKey:PackageID"`
	AreaOfConcerns []AreaOfConcern `json:"area_of_concerns,omitempty" gorm:"foreignKey:AreaOfConcernGroupID"`
}
