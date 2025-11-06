package service

import (
	"context"
	"net/http"

	"github.com/CRS-Project/crs-backend/internal/api/repository"
	"github.com/CRS-Project/crs-backend/internal/dto"
	"github.com/CRS-Project/crs-backend/internal/entity"
	myerror "github.com/CRS-Project/crs-backend/internal/pkg/error"
	"github.com/CRS-Project/crs-backend/internal/pkg/meta"
	"github.com/CRS-Project/crs-backend/internal/utils"
	"gorm.io/gorm"
)

type (
	UserService interface {
		Create(ctx context.Context, req dto.CreateUserRequest) (dto.CreateUserResponse, error)
		GetAll(ctx context.Context, metaReq meta.Meta) ([]dto.UserNonAdminDetailResponse, meta.Meta, error)
		GetById(ctx context.Context, userId string) (dto.UserNonAdminDetailResponse, error)
		Update(ctx context.Context, userId string, req dto.UpdateUserRequest) (dto.UserNonAdminDetailResponse, error)
		Delete(ctx context.Context, userId string) error
	}

	userService struct {
		userRepository           repository.UserRepository
		userDisciplineRepository repository.UserDisciplineRepository
		packageRepository        repository.PackageRepository
		db                       *gorm.DB
	}
)

func NewUser(userRepository repository.UserRepository,
	userDisciplineRepository repository.UserDisciplineRepository,
	packageRepository repository.PackageRepository,
	db *gorm.DB) UserService {
	return &userService{
		userRepository:           userRepository,
		userDisciplineRepository: userDisciplineRepository,
		packageRepository:        packageRepository,
		db:                       db,
	}
}

func (s *userService) Create(ctx context.Context, req dto.CreateUserRequest) (dto.CreateUserResponse, error) {
	disciplineId := ""

	if req.Role == "CONTRACTOR" && req.DisciplineID == nil {
		contractorDisc, err := s.userDisciplineRepository.GetContractorDiscipline(ctx, nil)
		if err != nil {
			return dto.CreateUserResponse{}, err
		}

		userDiscipline, err := s.userDisciplineRepository.GetByID(ctx, nil, disciplineId, "Users")
		if err != nil {
			return dto.CreateUserResponse{}, err
		}

		for _, u := range userDiscipline.Users {
			if u.PackageID.String() == req.PackageID {
				return dto.CreateUserResponse{}, myerror.New("this package has already contractor", http.StatusBadRequest)
			}
		}

		disciplineId = contractorDisc.ID.String()
	} else {
		disciplineId = *req.DisciplineID
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
		Name:             req.Name,
		Email:            req.Email,
		Password:         hashPassword,
		IsVerified:       true,
		Role:             entity.Role(req.Role),
		Initial:          req.Initial,
		Institution:      req.Institution,
		PhotoProfile:     req.PhotoProfile,
		DisciplineNumber: req.DisciplineNumber,
		UserDisciplineID: discipline.ID,
		PackageID:        &pkg.ID,
	})
	if err != nil {
		return dto.CreateUserResponse{}, err
	}

	return dto.CreateUserResponse{
		ID:               userCreated.ID.String(),
		Name:             userCreated.Name,
		Email:            userCreated.Email,
		Initial:          userCreated.Initial,
		Institution:      userCreated.Institution,
		DisciplineNumber: userCreated.DisciplineNumber,
		PhotoProfile:     userCreated.PhotoProfile,
		IsVerified:       true,
		Role:             string(userCreated.Role),
		Package:          pkg.Name,
		Discipline:       discipline.Name,
	}, nil
}

func (s *userService) GetAll(ctx context.Context, metaReq meta.Meta) ([]dto.UserNonAdminDetailResponse, meta.Meta, error) {
	users, metaRes, err := s.userRepository.GetAll(ctx, nil, metaReq, "UserDiscipline", "Package")
	if err != nil {
		return nil, metaReq, err
	}

	var res []dto.UserNonAdminDetailResponse

	for _, user := range users {
		pkg := "all"
		if user.Package != nil {
			pkg = user.Package.Name
		}
		res = append(res, dto.UserNonAdminDetailResponse{
			ID:               user.ID.String(),
			Name:             user.Name,
			Email:            user.Email,
			Initial:          user.Initial,
			Institution:      user.Institution,
			PhotoProfile:     user.PhotoProfile,
			Role:             string(user.Role),
			DisciplineNumber: user.DisciplineNumber,
			Discipline:       user.UserDiscipline.Name,
			Package:          pkg,
		})
	}

	return res, metaRes, nil
}

func (s *userService) GetById(ctx context.Context, userId string) (dto.UserNonAdminDetailResponse, error) {
	user, err := s.userRepository.GetById(ctx, nil, userId, "UserDiscipline", "Package")
	if err != nil {
		return dto.UserNonAdminDetailResponse{}, err
	}
	pkg := "all"
	if user.Package != nil {
		pkg = user.Package.Name
	}
	return dto.UserNonAdminDetailResponse{
		ID:               userId,
		Name:             user.Name,
		Email:            user.Email,
		Initial:          user.Initial,
		Institution:      user.Institution,
		PhotoProfile:     user.PhotoProfile,
		Role:             string(user.Role),
		DisciplineNumber: user.DisciplineNumber,
		Discipline:       user.UserDiscipline.Name,
		Package:          pkg,
	}, nil
}

func (s *userService) Update(ctx context.Context, userId string, req dto.UpdateUserRequest) (dto.UserNonAdminDetailResponse, error) {
	user, err := s.userRepository.GetById(ctx, nil, userId, "UserDiscipline")
	if err != nil {
		return dto.UserNonAdminDetailResponse{}, err
	}

	hashPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return dto.UserNonAdminDetailResponse{}, err
	}

	pkg, err := s.packageRepository.GetByID(ctx, nil, req.PackageID)
	if err != nil {
		return dto.UserNonAdminDetailResponse{}, err
	}

	var disciplineID string
	if req.DisciplineID == nil {
		disciplineID = user.UserDisciplineID.String()
	} else {
		disciplineID = *req.DisciplineID
	}

	discipline, err := s.userDisciplineRepository.GetByID(ctx, nil, disciplineID)
	if err != nil {
		return dto.UserNonAdminDetailResponse{}, err
	}

	user.Name = req.Name
	user.Email = req.Email
	user.Password = hashPassword
	user.Initial = req.Initial
	user.Institution = req.Institution
	user.PackageID = &pkg.ID
	user.DisciplineNumber = req.DisciplineNumber
	if req.DisciplineID != nil {
		user.UserDisciplineID = discipline.ID
	}

	_, err = s.userRepository.Update(ctx, nil, user)
	if err != nil {
		return dto.UserNonAdminDetailResponse{}, err
	}

	return dto.UserNonAdminDetailResponse{
		ID:               user.ID.String(),
		Name:             user.Name,
		Email:            user.Email,
		Initial:          user.Initial,
		Institution:      user.Institution,
		PhotoProfile:     user.PhotoProfile,
		Role:             string(user.Role),
		DisciplineNumber: user.DisciplineNumber,
		Package:          pkg.Name,
		Discipline:       discipline.Name,
	}, nil
}

func (s *userService) Delete(ctx context.Context, userId string) error {
	user, err := s.userRepository.GetById(ctx, nil, userId)
	if err != nil {
		return err
	}

	return s.userRepository.Delete(ctx, nil, user)
}
