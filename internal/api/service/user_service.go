package service

import (
	"context"

	"github.com/CRS-Project/crs-backend/internal/api/repository"
	"github.com/CRS-Project/crs-backend/internal/dto"
	"gorm.io/gorm"
)

type (
	UserService interface {
		GetById(ctx context.Context, userId string) (dto.UserResponse, error)
	}

	userService struct {
		userRepository repository.UserRepository
		db             *gorm.DB
	}
)

func NewUser(userRepository repository.UserRepository,
	db *gorm.DB) UserService {
	return &userService{
		userRepository: userRepository,
		db:             db,
	}
}

func (s *userService) GetById(ctx context.Context, userId string) (dto.UserResponse, error) {
	// user, err := s.userRepository.GetByIdWithFilmList(ctx, nil, userId)
	// if err != nil {
	// 	return dto.UserResponse{}, err
	// }

	// return dto.UserResponse{
	// 	ID:          user.ID.String(),
	// 	Username:    user.Username,
	// 	PhoneNumber: us,
	// }, nil
	return dto.UserResponse{}, nil
}
