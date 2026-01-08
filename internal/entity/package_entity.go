package entity

import (
	"github.com/CRS-Project/crs-backend/internal/dto"
	"github.com/google/uuid"
)

type Package struct {
	ID          uuid.UUID `json:"id" gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	Name        string    `json:"name" gorm:"not null;"`
	Description string    `json:"description" gorm:""`

	DisciplineGroups []DisciplineGroup `json:"discipline_groups,omitempty" gorm:"foreignKey:PackageID"`

	DeletedBy uuid.UUID `json:"deleted_by" gorm:"type:uuid"`
	UpdatedBy uuid.UUID `json:"updated_by" gorm:"type:uuid"`
	Timestamp
}

func (p *Package) ToInfo() dto.PackageInfo {
	return dto.PackageInfo{
		ID:          p.ID.String(),
		Name:        p.Name,
		Description: p.Description,
	}
}
