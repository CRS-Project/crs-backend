package entity

import (
	"github.com/CRS-Project/crs-backend/internal/dto"
	"github.com/google/uuid"
)

type Package struct {
	ID          uuid.UUID `json:"id" gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	Name        string    `json:"name" gorm:"not null;"`
	Description string    `json:"description" gorm:""`

	Timestamp
}

func (p *Package) ToInfo() dto.PackageInfo {
	return dto.PackageInfo{
		ID:          p.ID.String(),
		Name:        p.Name,
		Description: p.Description,
	}
}
