package service

import (
	"context"

	"github.com/CRS-Project/crs-backend/internal/api/repository"
	"github.com/CRS-Project/crs-backend/internal/dto"
	"github.com/CRS-Project/crs-backend/internal/entity"
	"github.com/CRS-Project/crs-backend/internal/pkg/meta"
	"github.com/CRS-Project/crs-backend/internal/utils"
	"github.com/google/uuid"
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
		userRepository                 repository.UserRepository
		userDisciplineRepository       repository.UserDisciplineRepository
		userDisciplineNumberRepository repository.UserDisciplineNumberRepository
		userPackageRepository          repository.UserPackageRepository
		packageRepository              repository.PackageRepository
		db                             *gorm.DB
	}
)

func NewUser(userRepository repository.UserRepository,
	userDisciplineRepository repository.UserDisciplineRepository,
	userDisciplineNumberRepository repository.UserDisciplineNumberRepository,
	userPackageRepository repository.UserPackageRepository,
	packageRepository repository.PackageRepository,
	db *gorm.DB) UserService {
	return &userService{
		userRepository:                 userRepository,
		userDisciplineRepository:       userDisciplineRepository,
		userDisciplineNumberRepository: userDisciplineNumberRepository,
		userPackageRepository:          userPackageRepository,
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
		ID:          userCreated.ID.String(),
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

func (s *userService) GetAll(ctx context.Context, metaReq meta.Meta) ([]dto.UserNonAdminDetailResponse, meta.Meta, error) {
	users, metaRes, err := s.userRepository.GetAll(ctx, nil, metaReq, "UserDisciplineNumber.UserDiscipline", "UserPackage.Package")
	if err != nil {
		return nil, metaReq, err
	}

	var res []dto.UserNonAdminDetailResponse

	for _, user := range users {
		var pkgAccess []*dto.PackageInfo
		for _, pkg := range user.UserPackage {
			pkgAccess = append(pkgAccess, &dto.PackageInfo{
				ID:   pkg.ID.String(),
				Name: pkg.Package.Name,
			})
		}

		pkg := "no package"
		if len(pkgAccess) > 0 && pkgAccess[0] != nil {
			pkg = pkgAccess[0].Name
		}

		res = append(res, dto.UserNonAdminDetailResponse{
			ID:           user.ID.String(),
			Name:         user.Name,
			Email:        user.Email,
			Initial:      user.Initial,
			Institution:  user.Institution,
			PhotoProfile: user.PhotoProfile,
			Role:         string(user.Role),
			Discipline:   user.UserDisciplineNumber.UserDiscipline.Name,
			Package:      pkg,
		})
	}

	return res, metaRes, nil
}

func (s *userService) GetById(ctx context.Context, userId string) (dto.UserNonAdminDetailResponse, error) {
	user, err := s.userRepository.GetById(ctx, nil, userId, "UserDisciplineNumber.UserDiscipline", "UserPackage.Package")
	if err != nil {
		return dto.UserNonAdminDetailResponse{}, err
	}

	var pkgAccess []*dto.PackageInfo
	for _, pkg := range user.UserPackage {
		pkgAccess = append(pkgAccess, &dto.PackageInfo{
			ID:   pkg.ID.String(),
			Name: pkg.Package.Name,
		})
	}

	pkg := "no package"
	if pkgAccess[0] != nil {
		pkg = pkgAccess[0].Name
	}

	return dto.UserNonAdminDetailResponse{
		ID:           userId,
		Name:         user.Name,
		Email:        user.Email,
		Initial:      user.Initial,
		Institution:  user.Institution,
		PhotoProfile: user.PhotoProfile,
		Role:         string(user.Role),
		Discipline:   user.UserDisciplineNumber.UserDiscipline.Name,
		Package:      pkg,
	}, nil
}

func (s *userService) Update(ctx context.Context, userId string, req dto.UpdateUserRequest) (dto.UserNonAdminDetailResponse, error) {
	tx := s.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	user, err := s.userRepository.GetById(ctx, tx, userId, "UserDisciplineNumber", "UserPackage")
	if err != nil {
		tx.Rollback()
		return dto.UserNonAdminDetailResponse{}, err
	}

	hashPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		tx.Rollback()
		return dto.UserNonAdminDetailResponse{}, err
	}

	pkg, err := s.packageRepository.GetByID(ctx, tx, req.PackageID)
	if err != nil {
		tx.Rollback()
		return dto.UserNonAdminDetailResponse{}, err
	}

	var disciplineID string
	if req.DisciplineID == nil {
		disciplineID = user.UserDisciplineNumber.UserDisciplineID.String()
	} else {
		disciplineID = *req.DisciplineID
	}

	discipline, err := s.userDisciplineRepository.GetByID(ctx, tx, disciplineID)
	if err != nil {
		tx.Rollback()
		return dto.UserNonAdminDetailResponse{}, err
	}

	user.Name = req.Name
	user.Email = req.Email
	user.Password = hashPassword
	user.Initial = req.Initial
	user.Institution = req.Institution

	_, err = s.userRepository.Update(ctx, tx, user)
	if err != nil {
		tx.Rollback()
		return dto.UserNonAdminDetailResponse{}, err
	}

	if user.UserDisciplineNumber != nil {
		user.UserDisciplineNumber.UserDisciplineID = uuid.MustParse(disciplineID)
		_, err := s.userDisciplineNumberRepository.Update(ctx, tx, *user.UserDisciplineNumber)
		if err != nil {
			tx.Rollback()
			return dto.UserNonAdminDetailResponse{}, err
		}
	}

	if len(user.UserPackage) > 0 {
		user.UserPackage[0].PackageID = pkg.ID
		_, err := s.userPackageRepository.Save(ctx, tx, user.UserPackage[0])
		if err != nil {
			tx.Rollback()
			return dto.UserNonAdminDetailResponse{}, err
		}
	} else {
		_, err := s.userPackageRepository.Create(ctx, tx, entity.UserPackage{
			UserID:    user.ID,
			PackageID: pkg.ID,
		})
		if err != nil {
			tx.Rollback()
			return dto.UserNonAdminDetailResponse{}, err
		}
	}

	if err := tx.Commit().Error; err != nil {
		return dto.UserNonAdminDetailResponse{}, err
	}

	return dto.UserNonAdminDetailResponse{
		ID:           user.ID.String(),
		Name:         user.Name,
		Email:        user.Email,
		Initial:      user.Initial,
		Institution:  user.Institution,
		PhotoProfile: user.PhotoProfile,
		Role:         string(user.Role),
		Discipline:   discipline.Name,
		Package:      pkg.Name,
	}, nil
}

func (s *userService) Delete(ctx context.Context, userId string) error {
	user, err := s.userRepository.GetById(ctx, nil, userId)
	if err != nil {
		return err
	}

	return s.userRepository.Delete(ctx, nil, user)
}
