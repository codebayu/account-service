package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/codebayu/account-service/common/apperror"
	"github.com/codebayu/account-service/internal/dto"
	"github.com/labstack/echo/v5"
	"github.com/stretchr/testify/assert"
)

func TestAuthHandler_RefreshToken(t *testing.T) {
	mockSvc := new(MockAuthService)
	h := NewAuthHandler(mockSvc)
	e := echo.New()

	t.Run("Success Refresh Token", func(t *testing.T) {
		body, _ := json.Marshal(map[string]string{"refreshToken": "valid-refresh-token"})
		req := httptest.NewRequest(http.MethodPost, "/auth/refresh-token", bytes.NewReader(body))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		mockSvc.On("RefreshToken", context.Background(), "valid-refresh-token").
			Return(&dto.AuthResponseData{AccessToken: "new-access-token"}, nil).Once()

		err := h.RefreshToken(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
		mockSvc.AssertExpectations(t)
	})

	t.Run("Missing refreshToken", func(t *testing.T) {
		body, _ := json.Marshal(map[string]string{})
		req := httptest.NewRequest(http.MethodPost, "/auth/refresh-token", bytes.NewReader(body))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := h.RefreshToken(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("Invalid Token", func(t *testing.T) {
		body, _ := json.Marshal(map[string]string{"refreshToken": "bad-token"})
		req := httptest.NewRequest(http.MethodPost, "/auth/refresh-token", bytes.NewReader(body))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		mockSvc.On("RefreshToken", context.Background(), "bad-token").
			Return(nil, apperror.ErrInvalidToken).Once()

		err := h.RefreshToken(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusUnauthorized, rec.Code)
		mockSvc.AssertExpectations(t)
	})
}

func TestAuthHandler_Logout(t *testing.T) {
	mockSvc := new(MockAuthService)
	h := NewAuthHandler(mockSvc)
	e := echo.New()

	t.Run("Success Logout", func(t *testing.T) {
		body, _ := json.Marshal(map[string]string{"refreshToken": "valid-refresh-token"})
		req := httptest.NewRequest(http.MethodPost, "/auth/logout", bytes.NewReader(body))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		req.Header.Set("Authorization", "Bearer valid-access-token")
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		mockSvc.On("Logout", context.Background(), "valid-refresh-token", "valid-access-token").
			Return(nil).Once()

		err := h.Logout(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
		mockSvc.AssertExpectations(t)
	})

	t.Run("Missing refreshToken", func(t *testing.T) {
		body, _ := json.Marshal(map[string]string{})
		req := httptest.NewRequest(http.MethodPost, "/auth/logout", bytes.NewReader(body))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := h.Logout(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("Service Error", func(t *testing.T) {
		body, _ := json.Marshal(map[string]string{"refreshToken": "valid-refresh-token"})
		req := httptest.NewRequest(http.MethodPost, "/auth/logout", bytes.NewReader(body))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		mockSvc.On("Logout", context.Background(), "valid-refresh-token", "").
			Return(errors.New("service error")).Once()

		err := h.Logout(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		mockSvc.AssertExpectations(t)
	})

	t.Run("Logout No Auth Header", func(t *testing.T) {
		body, _ := json.Marshal(map[string]string{"refreshToken": "valid-refresh-token"})
		req := httptest.NewRequest(http.MethodPost, "/auth/logout", bytes.NewReader(body))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		// No Authorization header
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		mockSvc.On("Logout", context.Background(), "valid-refresh-token", "").
			Return(nil).Once()

		err := h.Logout(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
		mockSvc.AssertExpectations(t)
	})
}
