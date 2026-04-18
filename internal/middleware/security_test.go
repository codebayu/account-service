package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/codebayu/account-service/internal/config"
	"github.com/labstack/echo/v5"
	"github.com/stretchr/testify/assert"
)

func TestSecurityMiddleware(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	mw := SecurityMiddleware()
	handler := mw(func(c *echo.Context) error {
		return c.String(http.StatusOK, "ok")
	})

	err := handler(c)
	assert.NoError(t, err)
	assert.Equal(t, "1; mode=block", rec.Header().Get(echo.HeaderXXSSProtection))
	assert.Equal(t, "nosniff", rec.Header().Get(echo.HeaderXContentTypeOptions))
}

func TestCORSMiddleware(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodOptions, "/", nil)
	req.Header.Set(echo.HeaderOrigin, "http://localhost:3000")
	req.Header.Set(echo.HeaderAccessControlRequestMethod, http.MethodPost)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	cfg := &config.Config{
		AllowedOrigins: []string{"http://localhost:3000"},
	}

	mw := CORSMiddleware(cfg)
	handler := mw(func(c *echo.Context) error {
		return c.String(http.StatusOK, "ok")
	})

	err := handler(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusNoContent, rec.Code) // Preflight returns 204
	assert.Equal(t, "http://localhost:3000", rec.Header().Get(echo.HeaderAccessControlAllowOrigin))
}
