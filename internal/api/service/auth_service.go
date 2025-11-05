package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/CRS-Project/crs-backend/internal/api/repository"
	"github.com/CRS-Project/crs-backend/internal/dto"
	mailer "github.com/CRS-Project/crs-backend/internal/pkg/email"
	myerror "github.com/CRS-Project/crs-backend/internal/pkg/error"
	"github.com/CRS-Project/crs-backend/internal/pkg/google/oauth"
	myjwt "github.com/CRS-Project/crs-backend/internal/pkg/jwt"
	"github.com/CRS-Project/crs-backend/internal/utils"

	"net/http"
	"os"
	"time"

	"gorm.io/gorm"
)

type (
	AuthService interface {
		Login(ctx context.Context, req dto.LoginRequest) (dto.LoginResponse, error)
		ForgetPassword(ctx context.Context, req dto.ForgetPasswordRequest) error
		ChangePassword(ctx context.Context, req dto.ChangePasswordRequest) error
		GetMe(ctx context.Context, userId string) (dto.GetMe, error)
	}

	authService struct {
		userRepository repository.UserRepository
		mailService    mailer.Mailer
		oauthService   oauth.Oauth
		db             *gorm.DB
	}
)

func NewAuth(userRepository repository.UserRepository,
	mailService mailer.Mailer,
	oauthService oauth.Oauth,
	db *gorm.DB) AuthService {
	return &authService{
		userRepository: userRepository,
		mailService:    mailService,
		oauthService:   oauthService,
		db:             db,
	}
}

func (s *authService) Login(ctx context.Context, req dto.LoginRequest) (dto.LoginResponse, error) {
	user, err := s.userRepository.GetByEmail(ctx, nil, req.Email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return dto.LoginResponse{}, myerror.New("email or password invalid", http.StatusBadRequest)
		}
		return dto.LoginResponse{}, err
	}

	if !user.IsVerified {
		return dto.LoginResponse{}, myerror.New("user is not verify", http.StatusUnauthorized)
	}

	checkPassword, err := utils.CheckPassword(user.Password, []byte(req.Password))
	if !checkPassword || err != nil {
		return dto.LoginResponse{}, myerror.New("email or password invalid", http.StatusBadRequest)
	}

	token, err := myjwt.GenerateToken(map[string]string{
		"user_id": user.ID.String(),
		"email":   user.Email,
		"role":    string(user.Role),
	}, 24*time.Hour)
	if err != nil {
		return dto.LoginResponse{}, err
	}

	return dto.LoginResponse{
		Token: token,
		Role:  string(user.Role),
	}, nil
}

func (s *authService) ForgetPassword(ctx context.Context, req dto.ForgetPasswordRequest) error {
	user, err := s.userRepository.GetByEmail(ctx, nil, req.Email)
	if err != nil {
		return err
	}

	if !user.IsVerified {
		return errors.New("user not verified")
	}

	token, err := myjwt.GenerateToken(map[string]string{
		"user_id": user.ID.String(),
		"email":   user.Email,
	}, 24*time.Hour)
	if err != nil {
		return err
	}

	// generate token
	token = fmt.Sprintf("%s/reset-password/%s", os.Getenv("APP_URL"), token)
	if err := s.mailService.MakeMail("./internal/pkg/email/template/forget_password_email.html", map[string]any{
		"Fullname": user.Name,
		"Link":     token,
	}).Send(user.Email, "Forget Password").Error; err != nil {
		return err
	}

	return nil
}

func (s *authService) ChangePassword(ctx context.Context, req dto.ChangePasswordRequest) error {
	user, err := s.userRepository.GetByEmail(ctx, nil, req.Email)
	if err != nil {
		return err
	}

	hashedPassword, err := utils.HashPassword(req.NewPassword)
	if err != nil {
		return err
	}

	user.Password = hashedPassword

	_, err = s.userRepository.Update(ctx, nil, user)
	if err != nil {
		return err
	}

	return nil
}

func (s *authService) GetMe(ctx context.Context, userId string) (dto.GetMe, error) {
	user, err := s.userRepository.GetById(ctx, nil, userId, "UserDisciplineNumber.UserDiscipline", "UserPackage.Package")
	if err != nil {
		return dto.GetMe{}, err
	}

	var pkgAccess []dto.PackageInfo
	for _, pkg := range user.UserPackage {
		pkgAccess = append(pkgAccess, dto.PackageInfo{
			ID:   pkg.ID.String(),
			Name: pkg.Package.Name,
		})
	}

	return dto.GetMe{
		PersonalInfo: dto.PersonalInfo{
			ID:           userId,
			Name:         user.Name,
			Email:        user.Email,
			Institution:  user.Institution,
			PhotoProfile: user.PhotoProfile,
			Initial:      user.Initial,
			Role:         string(user.Role),
		},
		UserDisciplineInfo: dto.UserDisciplineInfo{
			Discipline: user.UserDisciplineNumber.UserDiscipline.Name,
			Initial:    user.UserDisciplineNumber.UserDiscipline.Initial,
			Number:     user.UserDisciplineNumber.Number,
		},
		PackageAccess: pkgAccess,
	}, nil
}
