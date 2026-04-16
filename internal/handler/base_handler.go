package handler

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/codebayu/account-service/common/apperror"
	"github.com/codebayu/account-service/common/response"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v5"
)

func validateRequest(c *echo.Context, req interface{}) *apperror.AppError {
	if err := (*c).Bind(req); err != nil {
		return apperror.New(http.StatusBadRequest, 400000, "bad request", []string{err.Error()})
	}

	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		var errMessages []string
		for _, err := range err.(validator.ValidationErrors) {
			errMessages = append(errMessages, fmt.Sprintf("%s must be a valid %s", err.Field(), err.Tag()))
		}
		return apperror.New(http.StatusBadRequest, 400000, "bad request", errMessages)
	}

	return nil
}

func handleServiceError(c *echo.Context, err error) error {
	var appErr *apperror.AppError
	if errors.As(err, &appErr) {
		return (*c).JSON(appErr.StatusCode, response.Error(appErr.StatusCode, appErr.Code, appErr.Message, appErr.Errors))
	}
	return (*c).JSON(http.StatusInternalServerError, response.Error(http.StatusInternalServerError, 500000, "internal server error", []string{err.Error()}))
}
