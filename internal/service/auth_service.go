package service

import (
	"context"
	"errors"
	"time"

	"github.com/codebayu/account-service/internal/config"

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
	RefreshToken(ctx context.Context, tokenString string) (*dto.AuthResponseData, error)
	Logout(ctx context.Context, refreshTokenString string, accessTokenString string) error
}

type authService struct {
	repo           repository.UserRepository
	tokenRepo      repository.TokenRepository
	cfg            *config.Config
	tokenGenerator func(uuid string, duration time.Duration, isRefresh bool) (string, time.Time, error)
}

func NewAuthService(repo repository.UserRepository, tokenRepo repository.TokenRepository, cfg *config.Config) AuthService {
	return &authService{
		repo:           repo,
		tokenRepo:      tokenRepo,
		cfg:            cfg,
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
	accessDuration := time.Duration(s.cfg.JWTAccessExp) * time.Minute
	refreshDuration := time.Duration(s.cfg.JWTRefreshExp) * time.Minute

	accessToken, accessExp, err := s.tokenGenerator(uuid, accessDuration, false)
	if err != nil {
		return nil, apperror.ErrInternalServer
	}

	refreshToken, refreshExp, err := s.tokenGenerator(uuid, refreshDuration, true)
	if err != nil {
		return nil, apperror.ErrInternalServer
	}

	// Parse refresh token to get JTI to store in Redis
	claims, err := utils.ParseToken(refreshToken)
	if err != nil {
		return nil, apperror.ErrInternalServer
	}

	// Store refresh token info in Redis
	ctx := context.Background()
	err = s.tokenRepo.SaveRefreshToken(ctx, uuid, claims.ID, refreshDuration)
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

func (s *authService) RefreshToken(ctx context.Context, tokenString string) (*dto.AuthResponseData, error) {
	// 1. Parse & validate token signature and expiration
	claims, err := utils.ParseToken(tokenString)
	if err != nil {
		return nil, apperror.ErrInvalidToken
	}

	if !claims.IsRefresh {
		return nil, apperror.ErrInvalidToken
	}

	// 2. Check token in Redis
	_, err = s.tokenRepo.GetRefreshToken(ctx, claims.UUID, claims.ID)
	if err != nil {
		// Redis error or key not found (revoked)
		return nil, apperror.ErrInvalidToken
	}

	// 3. Delete old token from Redis (Token rotation)
	s.tokenRepo.DeleteRefreshToken(ctx, claims.UUID, claims.ID)

	// 4. Generate new tokens
	return s.generateAuthTokens(claims.UUID)
}

func (s *authService) Logout(ctx context.Context, refreshTokenString string, accessTokenString string) error {
	refreshClaims, err := utils.ParseToken(refreshTokenString)
	if err != nil {
		return apperror.ErrInvalidToken
	}

	if refreshClaims.IsRefresh {
		err = s.tokenRepo.DeleteRefreshToken(ctx, refreshClaims.UUID, refreshClaims.ID)
		if err != nil {
			return apperror.ErrInternalServer
		}
	}

	// Process Access Token Blocklisting if provided
	if accessTokenString != "" {
		accessClaims, err := utils.ParseToken(accessTokenString)
		if err == nil && !accessClaims.IsRefresh {
			// Calculate remaining duration
			expiresIn := time.Until(accessClaims.ExpiresAt.Time)
			if expiresIn > 0 {
				_ = s.tokenRepo.BlacklistAccessToken(ctx, accessClaims.ID, expiresIn)
			}
		}
	}

	return nil
}
