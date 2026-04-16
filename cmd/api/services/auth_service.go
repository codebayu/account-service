package services

import (
	"errors"
	"time"

	"github.com/codebayu/account-service/cmd/api/requests"
	"github.com/codebayu/account-service/internal/models"
	"github.com/codebayu/account-service/internal/utils"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthService struct {
	DB *gorm.DB
}

type RegisterResponseData struct {
	AccessToken        string    `json:"accessToken"`
	AccessTokenExpire  time.Time `json:"accessTokenExpire"`
	RefreshToken       string    `json:"refreshToken"`
	RefreshTokenExpire time.Time `json:"refreshTokenExpire"`
}

func (s *AuthService) Register(req requests.RegisterRequest) (*RegisterResponseData, error) {
	// Check if email already exists
	var existingUser models.User
	err := s.DB.Where("email = ?", req.Email).First(&existingUser).Error
	if err == nil {
		return nil, errors.New("email already registered")
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
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

	if err := s.DB.Create(&user).Error; err != nil {
		return nil, err
	}

	// Reload to get UUID and other fields
	s.DB.First(&user, user.ID)

	// Generate tokens
	accessToken, accessExp, err := utils.GenerateToken(user.UUID, 1*time.Hour, false)
	if err != nil {
		return nil, err
	}

	refreshToken, refreshExp, err := utils.GenerateToken(user.UUID, 7*24*time.Hour, true)
	if err != nil {
		return nil, err
	}

	return &RegisterResponseData{
		AccessToken:        accessToken,
		AccessTokenExpire:  accessExp,
		RefreshToken:       refreshToken,
		RefreshTokenExpire: refreshExp,
	}, nil
}
