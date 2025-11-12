package service

import (
	"context"

	"github.com/CRS-Project/crs-backend/internal/api/repository"
	"github.com/CRS-Project/crs-backend/internal/entity"
	"github.com/CRS-Project/crs-backend/internal/pkg/meta"
	"gorm.io/gorm"
)

type (
	UserDisciplineService interface {
		GetAll(ctx context.Context, metaReq meta.Meta) ([]entity.UserDiscipline, error)
	}

	userDisciplineService struct {
		userDisciplineRepository repository.UserDisciplineRepository
		db                       *gorm.DB
	}
)

func NewUserDiscipline(userDisciplineRepository repository.UserDisciplineRepository,
	db *gorm.DB) UserDisciplineService {
	return &userDisciplineService{
		userDisciplineRepository: userDisciplineRepository,
		db:                       db,
	}
}

func (s *userDisciplineService) GetAll(ctx context.Context, metaReq meta.Meta) ([]entity.UserDiscipline, error) {
	return s.userDisciplineRepository.GetAllNotAdminAndContractor(ctx, nil, metaReq)
}
