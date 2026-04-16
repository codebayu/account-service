package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/codebayu/account-service/internal/dto"
	"github.com/labstack/echo/v5"
	"github.com/stretchr/testify/assert"
)

func TestAuthHandler_Register(t *testing.T) {
	mockSvc := new(MockAuthService)
	h := NewAuthHandler(mockSvc)
	e := echo.New()

	reqBody := dto.RegisterRequest{
		Name:     "Test User",
		Email:    "test@example.com",
		Password: "password123",
		Gender:   "male",
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/auth/register", bytes.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	t.Run("Success Register", func(t *testing.T) {
		mockSvc.On("Register", reqBody).Return(&dto.AuthResponseData{AccessToken: "token"}, nil).Once()

		err := h.Register(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusCreated, rec.Code)
		mockSvc.AssertExpectations(t)
	})

	t.Run("Invalid Input", func(t *testing.T) {
		invalidBody, _ := json.Marshal(dto.RegisterRequest{Name: "Fail"})
		reqFail := httptest.NewRequest(http.MethodPost, "/auth/register", bytes.NewReader(invalidBody))
		reqFail.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		recFail := httptest.NewRecorder()
		cFail := e.NewContext(reqFail, recFail)

		err := h.Register(cFail)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, recFail.Code)
	})
}

func TestAuthHandler_Login(t *testing.T) {
	mockSvc := new(MockAuthService)
	h := NewAuthHandler(mockSvc)
	e := echo.New()

	reqBody := dto.LoginRequest{Email: "test@example.com", Password: "password123"}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/auth/login", bytes.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	t.Run("Success Login", func(t *testing.T) {
		mockSvc.On("Login", reqBody).Return(&dto.AuthResponseData{AccessToken: "token"}, nil).Once()

		err := h.Login(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
		mockSvc.AssertExpectations(t)
	})

	t.Run("Login Service Error", func(t *testing.T) {
		reqFail := httptest.NewRequest(http.MethodPost, "/auth/login", bytes.NewReader(body))
		reqFail.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		recFail := httptest.NewRecorder()
		cFail := e.NewContext(reqFail, recFail)

		mockSvc.On("Login", reqBody).Return(nil, errors.New("service error")).Once()

		err := h.Login(cFail)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, recFail.Code)
		mockSvc.AssertExpectations(t)
	})
}
