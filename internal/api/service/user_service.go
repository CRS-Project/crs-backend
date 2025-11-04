package service

import (
	"context"

	"github.com/CRS-Project/crs-backend/internal/api/repository"
	"github.com/CRS-Project/crs-backend/internal/dto"
	"gorm.io/gorm"
)

type (
	UserService interface {
		Create(ctx context.Context, req dto.CreateUserRequest) (dto.CreateUserResponse, error)
		GetById(ctx context.Context, userId string) (dto.UserResponse, error)
	}

	userService struct {
		userRepository       repository.UserRepository
		userDisciplineNumber repository.UserDisciplineNumberRepository
		db                   *gorm.DB
	}
)

func NewUser(userRepository repository.UserRepository,
	userDisciplineNumber repository.UserDisciplineNumberRepository,
	db *gorm.DB) UserService {
	return &userService{
		userRepository:       userRepository,
		userDisciplineNumber: userDisciplineNumber,
		db:                   db,
	}
}

func (s *userService) Create(ctx context.Context, req dto.CreateUserRequest) (dto.CreateUserResponse, error) {
	return dto.CreateUserResponse{}, nil
	// countUserDiscipline, err := s.userDisciplineNumber.CountByUserDisciplineID(ctx, nil,req.DisciplineID)
}

func (s *userService) GetById(ctx context.Context, userId string) (dto.UserResponse, error) {
	// user, err := s.userRepository.GetByIdWithFilmList(ctx, nil, userId)
	// if err != nil {
	// 	return dto.UserResponse{}, err
	// }

	// return dto.UserResponse{
	// 	ID:          user.ID.String(),
	// 	Name:    user.Name,
	// 	PhoneNumber: us,
	// }, nil
	return dto.UserResponse{}, nil
}
