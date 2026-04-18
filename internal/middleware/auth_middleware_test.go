package middleware

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"context"

	"github.com/codebayu/account-service/internal/utils"
	"github.com/labstack/echo/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockTokenRepository struct {
	mock.Mock
}

func (m *MockTokenRepository) SaveRefreshToken(ctx context.Context, userID, tokenID string, expiresIn time.Duration) error {
	return nil
}

func (m *MockTokenRepository) GetRefreshToken(ctx context.Context, userID, tokenID string) (string, error) {
	return "", nil
}

func (m *MockTokenRepository) DeleteRefreshToken(ctx context.Context, userID, tokenID string) error {
	return nil
}

func (m *MockTokenRepository) BlacklistAccessToken(ctx context.Context, tokenID string, expiresIn time.Duration) error {
	return nil
}

func (m *MockTokenRepository) IsAccessTokenBlacklisted(ctx context.Context, tokenID string) (bool, error) {
	args := m.Called(ctx, tokenID)
	return args.Bool(0), args.Error(1)
}

func TestAuthMiddleware(t *testing.T) {
	e := echo.New()
	os.Setenv("JWT_SECRET", "secret")

	mockTokenRepo := new(MockTokenRepository)
	
	mw := AuthMiddleware(mockTokenRepo)
	handler := mw(func(c *echo.Context) error {
		return (*c).String(http.StatusOK, "passed")
	})

	t.Run("Success Auth", func(t *testing.T) {
		token, _, _ := utils.GenerateToken("user-uuid", 1*time.Hour, false)

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		mockTokenRepo.On("IsAccessTokenBlacklisted", mock.Anything, mock.Anything).Return(false, nil).Once()

		err := handler(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, "user-uuid", (*c).Get("user_uuid"))
	})

	t.Run("Missing Auth Header", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := handler(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusUnauthorized, rec.Code)
	})

	t.Run("Invalid Auth Header Format", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.Header.Set("Authorization", "InvalidFormat token")
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := handler(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusUnauthorized, rec.Code)
	})

	t.Run("Invalid Token", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.Header.Set("Authorization", "Bearer invalid-token")
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := handler(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusUnauthorized, rec.Code)
	})

	t.Run("Refresh Token Not Allowed", func(t *testing.T) {
		token, _, _ := utils.GenerateToken("user-uuid", 1*time.Hour, true)

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := handler(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusUnauthorized, rec.Code)
	})

	t.Run("Token Blacklisted", func(t *testing.T) {
		token, _, _ := utils.GenerateToken("user-uuid", 1*time.Hour, false)

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		mockTokenRepo.On("IsAccessTokenBlacklisted", mock.Anything, mock.Anything).Return(true, nil).Once()

		err := handler(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusUnauthorized, rec.Code)
		assert.Contains(t, rec.Body.String(), "token has been revoked")
	})
}
