package service

import (
	"errors"
	"os"
	"testing"
	"time"

	"github.com/codebayu/account-service/common/apperror"
	"github.com/codebayu/account-service/internal/dto"
	"github.com/codebayu/account-service/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func TestAuthService_Register(t *testing.T) {
	os.Setenv("JWT_SECRET", "test_secret")
	defer os.Unsetenv("JWT_SECRET")

	mockRepo := new(MockUserRepository)
	svc := NewAuthService(mockRepo)

	req := dto.RegisterRequest{
		Name:     "Test User",
		Email:    "test@example.com",
		Password: "password123",
		Gender:   "male",
	}

	t.Run("Success Register", func(t *testing.T) {
		mockRepo.On("FindByEmail", req.Email).Return(nil, gorm.ErrRecordNotFound).Once()
		mockRepo.On("Create", mock.AnythingOfType("*models.User")).Return(nil).Once()

		res, err := svc.Register(req)

		assert.NoError(t, err)
		assert.NotNil(t, res)
		assert.NotEmpty(t, res.AccessToken)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Email Already Exists", func(t *testing.T) {
		mockRepo.On("FindByEmail", req.Email).Return(&models.User{Email: req.Email}, nil).Once()

		res, err := svc.Register(req)

		assert.Error(t, err)
		assert.Nil(t, res)
		assert.Equal(t, apperror.ErrEmailRegistered, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Internal Server Error on Create", func(t *testing.T) {
		mockRepo.On("FindByEmail", req.Email).Return(nil, gorm.ErrRecordNotFound).Once()
		mockRepo.On("Create", mock.AnythingOfType("*models.User")).Return(errors.New("db error")).Once()

		res, err := svc.Register(req)

		assert.Error(t, err)
		assert.Nil(t, res)
		assert.Equal(t, apperror.ErrInternalServer, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Internal Server Error on Token Generation", func(t *testing.T) {
		mockRepo.On("FindByEmail", req.Email).Return(nil, gorm.ErrRecordNotFound).Once()
		mockRepo.On("Create", mock.AnythingOfType("*models.User")).Return(nil).Once()

		s := NewAuthService(mockRepo).(*authService)
		s.tokenGenerator = func(uuid string, duration time.Duration, isRefresh bool) (string, time.Time, error) {
			return "", time.Time{}, errors.New("token error")
		}

		res, err := s.Register(req)

		assert.Error(t, err)
		assert.Nil(t, res)
		assert.Equal(t, apperror.ErrInternalServer, err)
		mockRepo.AssertExpectations(t)
	})
}

func TestAuthService_Login(t *testing.T) {
	os.Setenv("JWT_SECRET", "test_secret")
	defer os.Unsetenv("JWT_SECRET")

	mockRepo := new(MockUserRepository)
	svc := NewAuthService(mockRepo)

	password := "password123"
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	user := &models.User{
		UUID:     "user-uuid",
		Email:    "test@example.com",
		Password: string(hashedPassword),
	}

	t.Run("Success Login", func(t *testing.T) {
		mockRepo.On("FindByEmail", user.Email).Return(user, nil).Once()

		res, err := svc.Login(dto.LoginRequest{Email: user.Email, Password: password})

		assert.NoError(t, err)
		assert.NotNil(t, res)
		assert.NotEmpty(t, res.AccessToken)
		mockRepo.AssertExpectations(t)
	})

	t.Run("User Not Found", func(t *testing.T) {
		mockRepo.On("FindByEmail", "notfound@example.com").Return(nil, gorm.ErrRecordNotFound).Once()

		res, err := svc.Login(dto.LoginRequest{Email: "notfound@example.com", Password: password})

		assert.Error(t, err)
		assert.Nil(t, res)
		assert.Equal(t, apperror.ErrUserNotFound, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Wrong Password", func(t *testing.T) {
		mockRepo.On("FindByEmail", user.Email).Return(user, nil).Once()

		res, err := svc.Login(dto.LoginRequest{Email: user.Email, Password: "wrongpassword"})

		assert.Error(t, err)
		assert.Nil(t, res)
		assert.Equal(t, apperror.ErrWrongPassword, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Internal Server Error on Find", func(t *testing.T) {
		mockRepo.On("FindByEmail", user.Email).Return(nil, errors.New("db error")).Once()

		res, err := svc.Login(dto.LoginRequest{Email: user.Email, Password: password})

		assert.Error(t, err)
		assert.Nil(t, res)
		assert.Equal(t, apperror.ErrInternalServer, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Internal Server Error on Refresh Token Generation", func(t *testing.T) {
		mockRepo.On("FindByEmail", user.Email).Return(user, nil).Once()

		s := NewAuthService(mockRepo).(*authService)
		s.tokenGenerator = func(uuid string, duration time.Duration, isRefresh bool) (string, time.Time, error) {
			if isRefresh {
				return "", time.Time{}, errors.New("token error")
			}
			return "access-token", time.Now(), nil
		}

		res, err := s.Login(dto.LoginRequest{Email: user.Email, Password: password})

		assert.Error(t, err)
		assert.Nil(t, res)
		assert.Equal(t, apperror.ErrInternalServer, err)
		mockRepo.AssertExpectations(t)
	})
}
