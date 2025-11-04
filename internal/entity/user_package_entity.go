package entity

import "github.com/google/uuid"

type UserPackage struct {
	ID        uuid.UUID `json:"id" gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	UserID    uuid.UUID `json:"user_id" gorm:"not null"`
	PackageID uuid.UUID `json:"package_id" gorm:"not null"`

	Timestamp

	Package *Package `json:"package,omitempty" gorm:"foreignKey:PackageID"`
}
