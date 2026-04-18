package middleware

import (
	"testing"

	"github.com/codebayu/account-service/internal/config"
	"github.com/labstack/echo/v5"
	"github.com/stretchr/testify/assert"
)

func TestRegisterGlobal(t *testing.T) {
	e := echo.New()
	cfg := &config.Config{
		AllowedOrigins: []string{"http://localhost:3000"},
	}

	RegisterGlobal(e, cfg)
	
	// Simply ensuring no panic and middleware was registered without crashing
	assert.NotNil(t, e)
}
