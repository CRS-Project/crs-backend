package service

import (
	"context"

	"github.com/CRS-Project/crs-backend/internal/api/repository"
	"github.com/CRS-Project/crs-backend/internal/dto"
	"gorm.io/gorm"
)

type (
	StatisticService interface {
		GetAOCAndCommentChart(ctx context.Context, packageId string) ([]dto.StatisticAOCAndCommentChart, error)
		GetCommentCard(ctx context.Context, packageId string) (dto.StatisticAOCAndCommentCard, error)
		GetCommentUserChart(ctx context.Context, packageId string) ([]dto.StatisticCommentUsersChart, error)
	}

	statisticService struct {
		statisticRepository     repository.StatisticRepository
		commentRepository       repository.CommentRepository
		documentRepository      repository.DocumentRepository
		areaOfConcernRepository repository.AreaOfConcernRepository
		userRepository          repository.UserRepository
		db                      *gorm.DB
	}
)

func NewStatistic(statisticRepository repository.StatisticRepository,
	commentRepository repository.CommentRepository,
	documentRepository repository.DocumentRepository,
	areaOfConcernRepository repository.AreaOfConcernRepository,
	userRepository repository.UserRepository,
	db *gorm.DB) StatisticService {
	return &statisticService{
		statisticRepository:     statisticRepository,
		commentRepository:       commentRepository,
		documentRepository:      documentRepository,
		areaOfConcernRepository: areaOfConcernRepository,
		userRepository:          userRepository,
		db:                      db,
	}
}

func (s *statisticService) GetAOCAndCommentChart(ctx context.Context, packageId string) ([]dto.StatisticAOCAndCommentChart, error) {
	return s.statisticRepository.GetAOCAndCommentChart(ctx, nil, packageId)
}

func (s *statisticService) GetCommentCard(ctx context.Context, packageId string) (dto.StatisticAOCAndCommentCard, error) {
	return s.statisticRepository.GetCommentCard(ctx, nil, packageId)
}

func (s *statisticService) GetCommentUserChart(ctx context.Context, packageId string) ([]dto.StatisticCommentUsersChart, error) {
	return s.statisticRepository.GetCommentUserChart(ctx, nil, packageId)
}
