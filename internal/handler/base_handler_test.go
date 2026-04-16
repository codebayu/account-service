package handler

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/codebayu/account-service/common/apperror"
	"github.com/labstack/echo/v5"
	"github.com/stretchr/testify/assert"
)

func TestHandleServiceError(t *testing.T) {
	e := echo.New()

	t.Run("AppError Mapping", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := handleServiceError(c, apperror.ErrWrongPassword)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusUnauthorized, rec.Code)
		assert.Contains(t, rec.Body.String(), "wrong password")
	})

	t.Run("Default Error Mapping", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := handleServiceError(c, errors.New("unknown error"))

		assert.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		assert.Contains(t, rec.Body.String(), "internal server error")
	})
}
