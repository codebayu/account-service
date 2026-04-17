package service

import (
	"errors"
	"time"

	"github.com/codebayu/account-service/common/apperror"
	"github.com/codebayu/account-service/internal/dto"
	"github.com/codebayu/account-service/internal/models"
	"github.com/codebayu/account-service/internal/repository"
	"github.com/codebayu/account-service/internal/utils"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthService interface {
	Register(req dto.RegisterRequest) (*dto.AuthResponseData, error)
	Login(req dto.LoginRequest) (*dto.AuthResponseData, error)
}

type authService struct {
	repo           repository.UserRepository
	tokenGenerator func(uuid string, duration time.Duration, isRefresh bool) (string, time.Time, error)
}

func NewAuthService(repo repository.UserRepository) AuthService {
	return &authService{
		repo:           repo,
		tokenGenerator: utils.GenerateToken,
	}
}

func (s *authService) Register(req dto.RegisterRequest) (*dto.AuthResponseData, error) {
	// Check if email already exists
	_, err := s.repo.FindByEmail(req.Email)
	if err == nil {
		return nil, apperror.ErrEmailRegistered
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, apperror.ErrInternalServer
	}

	// Create user
	user := models.User{
		Name:      req.Name,
		Email:     req.Email,
		Password:  string(hashedPassword),
		Gender:    req.Gender,
		CreatedBy: "system",
		UpdatedBy: "system",
	}

	if err := s.repo.Create(&user); err != nil {
		return nil, apperror.ErrInternalServer
	}

	return s.generateAuthTokens(user.UUID)
}

func (s *authService) Login(req dto.LoginRequest) (*dto.AuthResponseData, error) {
	user, err := s.repo.FindByEmail(req.Email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperror.ErrUserNotFound
		}
		return nil, apperror.ErrInternalServer
	}

	// Compare password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, apperror.ErrWrongPassword
	}

	return s.generateAuthTokens(user.UUID)
}

func (s *authService) generateAuthTokens(uuid string) (*dto.AuthResponseData, error) {
	accessToken, accessExp, err := s.tokenGenerator(uuid, 1*time.Hour, false)
	if err != nil {
		return nil, apperror.ErrInternalServer
	}

	refreshToken, refreshExp, err := s.tokenGenerator(uuid, 7*24*time.Hour, true)
	if err != nil {
		return nil, apperror.ErrInternalServer
	}

	return &dto.AuthResponseData{
		AccessToken:        accessToken,
		AccessTokenExpire:  accessExp,
		RefreshToken:       refreshToken,
		RefreshTokenExpire: refreshExp,
	}, nil
}
