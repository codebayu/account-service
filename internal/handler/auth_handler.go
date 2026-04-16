package handler

import (
	"net/http"

	"github.com/codebayu/account-service/common/response"
	"github.com/codebayu/account-service/internal/dto"
	"github.com/codebayu/account-service/internal/service"
	"github.com/labstack/echo/v5"
)

type AuthHandler struct {
	authService service.AuthService
}

func NewAuthHandler(authService service.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

func (h *AuthHandler) Register(c *echo.Context) error {
	var req dto.RegisterRequest
	if appErr := validateRequest(c, &req); appErr != nil {
		return (*c).JSON(appErr.StatusCode, response.Error(appErr.StatusCode, appErr.Code, appErr.Message, appErr.Errors))
	}

	data, err := h.authService.Register(req)
	if err != nil {
		return handleServiceError(c, err)
	}

	return (*c).JSON(http.StatusCreated, response.Success(http.StatusCreated, 201000, "created", data))
}

func (h *AuthHandler) Login(c *echo.Context) error {
	var req dto.LoginRequest
	if appErr := validateRequest(c, &req); appErr != nil {
		return (*c).JSON(appErr.StatusCode, response.Error(appErr.StatusCode, appErr.Code, appErr.Message, appErr.Errors))
	}

	data, err := h.authService.Login(req)
	if err != nil {
		return handleServiceError(c, err)
	}

	return (*c).JSON(http.StatusOK, response.Success(http.StatusOK, 200000, "ok", data))
}
