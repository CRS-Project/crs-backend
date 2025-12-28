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
		userRepository                               repository.UserRepository
		userDisciplineRepository                     repository.UserDisciplineRepository
		disciplineGroupConsolidatorRepository        repository.DisciplineGroupConsolidatorRepository
		disciplineListDocumentConsolidatorRepository repository.DisciplineListDocumentConsolidatorRepository
		packageRepository                            repository.PackageRepository
		db                                           *gorm.DB
	}
)

func NewUser(userRepository repository.UserRepository,
	userDisciplineRepository repository.UserDisciplineRepository,
	disciplineGroupConsolidatorRepository repository.DisciplineGroupConsolidatorRepository,
	disciplineListDocumentConsolidatorRepository repository.DisciplineListDocumentConsolidatorRepository,
	packageRepository repository.PackageRepository,
	db *gorm.DB) UserService {
	return &userService{
		userRepository:                               userRepository,
		userDisciplineRepository:                     userDisciplineRepository,
		disciplineGroupConsolidatorRepository:        disciplineGroupConsolidatorRepository,
		disciplineListDocumentConsolidatorRepository: disciplineListDocumentConsolidatorRepository,
		packageRepository:                            packageRepository,
		db:                                           db,
	}
}

func (s *userService) Create(ctx context.Context, req dto.CreateUserRequest) (dto.CreateUserResponse, error) {
	disciplineId := ""

	if req.Role == "CONTRACTOR" && req.DisciplineID == nil {
		contractorDisc, err := s.userDisciplineRepository.GetContractorDiscipline(ctx, nil)
		if err != nil {
			return dto.CreateUserResponse{}, err
		}

		userDiscipline, err := s.userDisciplineRepository.GetByID(ctx, nil, contractorDisc.ID.String(), "Users")
		if err != nil {
			return dto.CreateUserResponse{}, err
		}

		for _, u := range userDiscipline.Users {
			if u.PackageID.String() == req.PackageID {
				return dto.CreateUserResponse{}, myerror.New("this package has already contractor", http.StatusBadRequest)
			}
		}

		disciplineId = contractorDisc.ID.String()
	} else if req.DisciplineID != nil {
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

	var pkgId *string
	if userCreated.Package != nil {
		pkgs := userCreated.PackageID.String()
		pkgId = &pkgs
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
		PackageID:        pkgId,
		DisciplineID:     disciplineId,
	}, nil
}

func (s *userService) GetAll(ctx context.Context, metaReq meta.Meta) ([]dto.UserNonAdminDetailResponse, meta.Meta, error) {
	users, metaRes, err := s.userRepository.GetAll(ctx, nil, metaReq, "UserDiscipline", "Package")
	if err != nil {
		return nil, metaReq, err
	}

	var res []dto.UserNonAdminDetailResponse

	for _, user := range users {
		var pkgId *string
		pkg := "All Access"
		if user.Package != nil {
			pkg = user.Package.Name
			pkgs := user.PackageID.String()
			pkgId = &pkgs
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
			PackageID:        pkgId,
			DisciplineID:     user.UserDisciplineID.String(),
		})
	}

	return res, metaRes, nil
}

func (s *userService) GetById(ctx context.Context, userId string) (dto.UserNonAdminDetailResponse, error) {
	user, err := s.userRepository.GetById(ctx, nil, userId, "UserDiscipline", "Package")
	if err != nil {
		return dto.UserNonAdminDetailResponse{}, err
	}
	var pkgId *string
	pkg := "all"
	if user.Package != nil {
		pkg = user.Package.Name
		pkgs := user.PackageID.String()
		pkgId = &pkgs
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
		PackageID:        pkgId,
		DisciplineID:     user.UserDisciplineID.String(),
	}, nil
}

func (s *userService) Update(ctx context.Context, userId string, req dto.UpdateUserRequest) (dto.UserNonAdminDetailResponse, error) {
	curUser, err := s.userRepository.GetById(ctx, nil, userId, "UserDiscipline")
	if err != nil {
		return dto.UserNonAdminDetailResponse{}, err
	}

	user, err := s.userRepository.GetById(ctx, nil, req.ID, "UserDiscipline")
	if err != nil {
		return dto.UserNonAdminDetailResponse{}, err
	}

	if curUser.Role != entity.RoleSuperAdmin && user.ID.String() != curUser.ID.String() {
		return dto.UserNonAdminDetailResponse{}, myerror.New("role not allowed", http.StatusUnauthorized)
	}

	if req.Password != nil {
		if curUser.Role != entity.RoleSuperAdmin || curUser.ID.String() == user.ID.String() {
			return dto.UserNonAdminDetailResponse{}, myerror.New("role not allowed to change password", http.StatusUnauthorized)
		}
		hashPassword, err := utils.HashPassword(*req.Password)
		if err != nil {
			return dto.UserNonAdminDetailResponse{}, err
		}
		user.Password = hashPassword
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
	user.Initial = req.Initial
	user.Institution = req.Institution
	user.DisciplineNumber = req.DisciplineNumber
	user.PhotoProfile = req.PhotoProfile
	if req.DisciplineID != nil {
		user.UserDisciplineID = discipline.ID
		user.UserDiscipline = nil
	}
	user.UpdatedBy = uuid.MustParse(userId)

	_, err = s.userRepository.Update(ctx, nil, user)
	if err != nil {
		return dto.UserNonAdminDetailResponse{}, err
	}

	var pkgId *string
	pkgres := "all"
	if user.Package != nil {
		pkgres = user.Package.Name
		pkgs := user.PackageID.String()
		pkgId = &pkgs
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
		Package:          pkgres,
		Discipline:       discipline.Name,
		PackageID:        pkgId,
		DisciplineID:     user.UserDisciplineID.String(),
	}, nil
}

func (s *userService) Delete(ctx context.Context, userId string) error {
	user, err := s.userRepository.GetById(ctx, nil, userId)
	if err != nil {
		return err
	}

	disciplineGroupConsolidators, err := s.disciplineGroupConsolidatorRepository.GetByUserID(ctx, nil, userId, "DisciplineListDocumentConsolidators")
	if err != nil {
		return err
	}

	var disciplineGroupConsolidatorIds []string
	// var disciplineListDocumentConsolidatorIds []string
	for _, c := range disciplineGroupConsolidators {
		disciplineGroupConsolidatorIds = append(disciplineGroupConsolidatorIds, c.ID.String())
		// for _, dld := range c.DisciplineListDocumentConsolidators {
		// 	disciplineListDocumentConsolidatorIds = append(disciplineListDocumentConsolidatorIds, dld.ID.String())
		// }
	}

	if err := s.disciplineListDocumentConsolidatorRepository.DeleteByDisciplineGroupConsolidatorID(ctx, nil, disciplineGroupConsolidatorIds); err != nil {
		return err
	}

	if err := s.disciplineGroupConsolidatorRepository.DeleteBulk(ctx, nil, disciplineGroupConsolidatorIds); err != nil {
		return err
	}

	user.DeletedBy = uuid.MustParse(userId)
	return s.userRepository.Delete(ctx, nil, user)
}
