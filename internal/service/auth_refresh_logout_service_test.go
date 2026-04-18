package service

import (
	"context"
	"errors"
	"os"
	"testing"
	"time"

	"github.com/codebayu/account-service/common/apperror"
	"github.com/codebayu/account-service/internal/config"
	"github.com/codebayu/account-service/internal/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestAuthService_RefreshToken(t *testing.T) {
	os.Setenv("JWT_SECRET", "test_secret")
	defer os.Unsetenv("JWT_SECRET")

	ctx := context.Background()
	mockRepo := new(MockUserRepository)
	mockTokenRepo := new(MockTokenRepository)
	mockCfg := &config.Config{
		JWTAccessExp:  1,
		JWTRefreshExp: 2,
	}
	svc := NewAuthService(mockRepo, mockTokenRepo, mockCfg)

	// Generate a valid refresh token for use in tests
	validRefreshToken, _, _ := utils.GenerateToken("user-uuid", 5*time.Minute, true)
	accessToken, _, _ := utils.GenerateToken("user-uuid", 5*time.Minute, false)

	t.Run("Success Refresh Token", func(t *testing.T) {
		claims, _ := utils.ParseToken(validRefreshToken)
		mockTokenRepo.On("GetRefreshToken", ctx, "user-uuid", claims.ID).Return("valid", nil).Once()
		mockTokenRepo.On("DeleteRefreshToken", ctx, "user-uuid", claims.ID).Return(nil).Once()
		mockTokenRepo.On("SaveRefreshToken", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil).Once()

		res, err := svc.RefreshToken(ctx, validRefreshToken)

		assert.NoError(t, err)
		assert.NotNil(t, res)
		assert.NotEmpty(t, res.AccessToken)
		mockTokenRepo.AssertExpectations(t)
	})

	t.Run("Invalid Token String", func(t *testing.T) {
		res, err := svc.RefreshToken(ctx, "not-a-valid-token")

		assert.Error(t, err)
		assert.Nil(t, res)
		assert.Equal(t, apperror.ErrInvalidToken, err)
	})

	t.Run("Access Token Used as Refresh Token", func(t *testing.T) {
		res, err := svc.RefreshToken(ctx, accessToken)

		assert.Error(t, err)
		assert.Nil(t, res)
		assert.Equal(t, apperror.ErrInvalidToken, err)
	})

	t.Run("Token Not Found in Redis (Revoked)", func(t *testing.T) {
		claims, _ := utils.ParseToken(validRefreshToken)
		mockTokenRepo.On("GetRefreshToken", ctx, "user-uuid", claims.ID).Return("", errors.New("redis: nil")).Once()

		res, err := svc.RefreshToken(ctx, validRefreshToken)

		assert.Error(t, err)
		assert.Nil(t, res)
		assert.Equal(t, apperror.ErrInvalidToken, err)
		mockTokenRepo.AssertExpectations(t)
	})

	t.Run("Redis Error on SaveRefreshToken", func(t *testing.T) {
		// Need a fresh valid refresh token
		freshToken, _, _ := utils.GenerateToken("user-uuid-2", 5*time.Minute, true)
		freshClaims, _ := utils.ParseToken(freshToken)

		mockTokenRepo.On("GetRefreshToken", ctx, "user-uuid-2", freshClaims.ID).Return("valid", nil).Once()
		mockTokenRepo.On("DeleteRefreshToken", ctx, "user-uuid-2", freshClaims.ID).Return(nil).Once()
		mockTokenRepo.On("SaveRefreshToken", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(errors.New("redis error")).Once()

		res, err := svc.RefreshToken(ctx, freshToken)

		assert.Error(t, err)
		assert.Nil(t, res)
		assert.Equal(t, apperror.ErrInternalServer, err)
		mockTokenRepo.AssertExpectations(t)
	})
}

func TestAuthService_Logout(t *testing.T) {
	os.Setenv("JWT_SECRET", "test_secret")
	defer os.Unsetenv("JWT_SECRET")

	ctx := context.Background()
	mockRepo := new(MockUserRepository)
	mockTokenRepo := new(MockTokenRepository)
	mockCfg := &config.Config{
		JWTAccessExp:  1,
		JWTRefreshExp: 2,
	}
	svc := NewAuthService(mockRepo, mockTokenRepo, mockCfg)

	validRefreshToken, _, _ := utils.GenerateToken("user-uuid", 5*time.Minute, true)
	validAccessToken, _, _ := utils.GenerateToken("user-uuid", 5*time.Minute, false)

	t.Run("Success Logout with Access Token Blacklisting", func(t *testing.T) {
		refreshClaims, _ := utils.ParseToken(validRefreshToken)
		accessClaims, _ := utils.ParseToken(validAccessToken)

		mockTokenRepo.On("DeleteRefreshToken", ctx, "user-uuid", refreshClaims.ID).Return(nil).Once()
		mockTokenRepo.On("BlacklistAccessToken", ctx, accessClaims.ID, mock.Anything).Return(nil).Once()

		err := svc.Logout(ctx, validRefreshToken, validAccessToken)

		assert.NoError(t, err)
		mockTokenRepo.AssertExpectations(t)
	})

	t.Run("Success Logout without Access Token", func(t *testing.T) {
		freshToken, _, _ := utils.GenerateToken("user-uuid-3", 5*time.Minute, true)
		freshClaims, _ := utils.ParseToken(freshToken)

		mockTokenRepo.On("DeleteRefreshToken", ctx, "user-uuid-3", freshClaims.ID).Return(nil).Once()

		err := svc.Logout(ctx, freshToken, "")

		assert.NoError(t, err)
		mockTokenRepo.AssertExpectations(t)
	})

	t.Run("Invalid Refresh Token", func(t *testing.T) {
		err := svc.Logout(ctx, "invalid-token", "")

		assert.Error(t, err)
		assert.Equal(t, apperror.ErrInvalidToken, err)
	})

	t.Run("Redis Error on Delete Refresh Token", func(t *testing.T) {
		freshToken, _, _ := utils.GenerateToken("user-uuid-4", 5*time.Minute, true)
		freshClaims, _ := utils.ParseToken(freshToken)

		mockTokenRepo.On("DeleteRefreshToken", ctx, "user-uuid-4", freshClaims.ID).Return(errors.New("redis error")).Once()

		err := svc.Logout(ctx, freshToken, "")

		assert.Error(t, err)
		assert.Equal(t, apperror.ErrInternalServer, err)
		mockTokenRepo.AssertExpectations(t)
	})

	t.Run("Expired Access Token is Skipped Gracefully", func(t *testing.T) {
		freshToken, _, _ := utils.GenerateToken("user-uuid-5", 5*time.Minute, true)
		freshClaims, _ := utils.ParseToken(freshToken)

		mockTokenRepo.On("DeleteRefreshToken", ctx, "user-uuid-5", freshClaims.ID).Return(nil).Once()

		// Pass an expired access token string — the service should just skip blacklisting
		err := svc.Logout(ctx, freshToken, "this-is-a-bad-access-token")

		assert.NoError(t, err)
		mockTokenRepo.AssertExpectations(t)
	})
}
