package handlers

import (
	"fmt"
	"net/http"

	"github.com/codebayu/account-service/cmd/api/requests"
	"github.com/codebayu/account-service/common/response"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v5"
)

func (h *Handler) Register(c *echo.Context) error {
	var req requests.RegisterRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, response.Error(http.StatusBadRequest, 400000, "bad request", []string{err.Error()}))
	}

	// Validation
	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		var errors []string
		for _, err := range err.(validator.ValidationErrors) {
			errors = append(errors, fmt.Sprintf("%s must be a valid %s", err.Field(), err.Tag()))
		}
		return c.JSON(http.StatusBadRequest, response.Error(http.StatusBadRequest, 400000, "bad request", errors))
	}

	// Service call
	data, err := h.AuthService.Register(req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.Error(http.StatusBadRequest, 400000, "bad request", []string{err.Error()}))
	}

	return c.JSON(http.StatusCreated, response.Success(http.StatusCreated, 201000, "created", data))
}

func (h *Handler) Login(c *echo.Context) error {
	var req requests.LoginRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, response.Error(http.StatusBadRequest, 400000, "bad request", []string{err.Error()}))
	}

	// Validation
	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		var errors []string
		for _, err := range err.(validator.ValidationErrors) {
			errors = append(errors, fmt.Sprintf("%s must be a valid %s", err.Field(), err.Tag()))
		}
		return c.JSON(http.StatusBadRequest, response.Error(http.StatusBadRequest, 400000, "bad request", errors))
	}

	// Service call
	data, err := h.AuthService.Login(req)
	if err != nil {
		if err.Error() == "wrong password" {
			return c.JSON(http.StatusUnauthorized, response.Error(http.StatusUnauthorized, 401008, "wrong password", nil))
		}
		return c.JSON(http.StatusBadRequest, response.Error(http.StatusBadRequest, 400000, "bad request", []string{err.Error()}))
	}

	return c.JSON(http.StatusOK, response.Success(http.StatusOK, 200000, "ok", data))
}
