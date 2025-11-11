package service

import (
	"github.com/CRS-Project/crs-backend/internal/api/repository"
	"gorm.io/gorm"
)

type (
	StatisticService interface {
	}

	statisticService struct {
		commentRepository       repository.CommentRepository
		documentRepository      repository.DocumentRepository
		areaOfConcernRepository repository.AreaOfConcernRepository
		userRepository          repository.UserRepository
		db                      *gorm.DB
	}
)

func NewStatistic(commentRepository repository.CommentRepository,
	documentRepository repository.DocumentRepository,
	areaOfConcernRepository repository.AreaOfConcernRepository,
	userRepository repository.UserRepository,
	db *gorm.DB) StatisticService {
	return &statisticService{
		commentRepository:       commentRepository,
		documentRepository:      documentRepository,
		areaOfConcernRepository: areaOfConcernRepository,
		userRepository:          userRepository,
		db:                      db,
	}
}
