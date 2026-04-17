package middleware

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/codebayu/account-service/internal/config"
	"github.com/labstack/echo/v5"
	"github.com/stretchr/testify/assert"
)

func TestSignatureMiddleware(t *testing.T) {
	e := echo.New()
	cfg := &config.Config{
		APIKey:    "test_key",
		APISecret: "test_secret",
		ChannelID: "WEB",
	}

	mw := SignatureMiddleware(cfg)
	handler := mw(func(c *echo.Context) error {
		return (*c).String(http.StatusOK, "passed")
	})

	t.Run("Success Validation", func(t *testing.T) {
		datetime := "1713280000"
		stringToHash := cfg.APIKey + datetime
		h := hmac.New(sha256.New, []byte(cfg.APISecret))
		h.Write([]byte(stringToHash))
		signature := hex.EncodeToString(h.Sum(nil))

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.Header.Set("x-signature", signature)
		req.Header.Set("x-datetime", datetime)
		req.Header.Set("x-channel", "WEB")
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := handler(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, "passed", rec.Body.String())
	})

	t.Run("Missing Headers", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := handler(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.Contains(t, rec.Body.String(), "are required")
	})

	t.Run("Invalid Signature", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.Header.Set("x-signature", "wrong_signature")
		req.Header.Set("x-datetime", "1713280000")
		req.Header.Set("x-channel", "WEB")
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := handler(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusUnauthorized, rec.Code)
		assert.Contains(t, rec.Body.String(), "invalid signature")
	})

	t.Run("Invalid Channel", func(t *testing.T) {
		datetime := "1713280000"
		stringToHash := cfg.APIKey + datetime
		h := hmac.New(sha256.New, []byte(cfg.APISecret))
		h.Write([]byte(stringToHash))
		signature := hex.EncodeToString(h.Sum(nil))

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.Header.Set("x-signature", signature)
		req.Header.Set("x-datetime", datetime)
		req.Header.Set("x-channel", "MOBILE") // Should be WEB
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := handler(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusUnauthorized, rec.Code)
		assert.Contains(t, rec.Body.String(), "invalid channel")
	})

	t.Run("Bypass Swagger Path", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/swagger/index.html", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := handler(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, "passed", rec.Body.String())
	})

	t.Run("Bypass Health Path", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/health", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := handler(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, "passed", rec.Body.String())
	})
}
