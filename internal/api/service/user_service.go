package service

import (
	"context"

	"github.com/CRS-Project/crs-backend/internal/api/repository"
	"github.com/CRS-Project/crs-backend/internal/dto"
	"github.com/CRS-Project/crs-backend/internal/entity"
	"github.com/CRS-Project/crs-backend/internal/utils"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type (
	UserService interface {
		Create(ctx context.Context, req dto.CreateUserRequest) (dto.CreateUserResponse, error)
		GetById(ctx context.Context, userId string) (dto.UserDetailResponse, error)
	}

	userService struct {
		userRepository                 repository.UserRepository
		userDisciplineRepository       repository.UserDisciplineRepository
		userDisciplineNumberRepository repository.UserDisciplineNumberRepository
		packageRepository              repository.PackageRepository
		db                             *gorm.DB
	}
)

func NewUser(userRepository repository.UserRepository,
	userDisciplineRepository repository.UserDisciplineRepository,
	userDisciplineNumberRepository repository.UserDisciplineNumberRepository,
	packageRepository repository.PackageRepository,
	db *gorm.DB) UserService {
	return &userService{
		userRepository:                 userRepository,
		userDisciplineRepository:       userDisciplineRepository,
		userDisciplineNumberRepository: userDisciplineNumberRepository,
		packageRepository:              packageRepository,
		db:                             db,
	}
}

func (s *userService) Create(ctx context.Context, req dto.CreateUserRequest) (dto.CreateUserResponse, error) {
	disciplineId := ""

	if req.Role == "CONTRACTOR" && req.DisciplineID == nil {
		contractorDisc, err := s.userDisciplineRepository.GetContractorDiscipline(ctx, nil)
		if err != nil {
			return dto.CreateUserResponse{}, err
		}

		disciplineId = contractorDisc.ID.String()
	} else {
		disciplineId = *req.DisciplineID
	}

	countUserDiscipline, err := s.userDisciplineNumberRepository.CountByUserDisciplineID(ctx, nil, disciplineId)
	if err != nil {
		return dto.CreateUserResponse{}, err
	}

	discipline, err := s.userDisciplineRepository.GetByID(ctx, nil, disciplineId)
	if err != nil {
		return dto.CreateUserResponse{}, err
	}

	pkg, err := s.packageRepository.GetByID(ctx, nil, req.PackageID)
	if err != nil {
		return dto.CreateUserResponse{}, err
	}

	hashPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return dto.CreateUserResponse{}, err
	}

	userCreated, err := s.userRepository.Create(ctx, nil, entity.User{
		Name:        req.Name,
		Email:       req.Email,
		Password:    hashPassword,
		IsVerified:  true,
		Role:        entity.Role(req.Role),
		Initial:     req.Initial,
		Institution: req.Institution,
		UserDisciplineNumber: &entity.UserDisciplineNumber{
			Number:           countUserDiscipline + 1,
			UserDisciplineID: uuid.MustParse(disciplineId),
			PackageID:        &pkg.ID,
		},
		UserPackage: []entity.UserPackage{
			{
				PackageID: pkg.ID,
			},
		},
	})
	if err != nil {
		return dto.CreateUserResponse{}, err
	}

	return dto.CreateUserResponse{
		Name:        userCreated.Name,
		Email:       userCreated.Email,
		Initial:     userCreated.Initial,
		Institution: userCreated.Institution,
		IsVerified:  true,
		Role:        string(userCreated.Role),
		Package:     pkg.Name,
		Discipline:  discipline.Name,
	}, nil
}

func (s *userService) GetById(ctx context.Context, userId string) (dto.UserDetailResponse, error) {
	// user, err := s.userRepository.GetByIdWithFilmList(ctx, nil, userId)
	// if err != nil {
	// 	return dto.UserResponse{}, err
	// }

	// return dto.UserResponse{
	// 	ID:          user.ID.String(),
	// 	Name:    user.Name,
	// 	PhoneNumber: us,
	// }, nil
	return dto.UserDetailResponse{}, nil
}
