package middleware

import (
	"github.com/codebayu/account-service/internal/config"
	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
)

func RegisterGlobal(e *echo.Echo, cfg *config.Config) {
	e.Use(middleware.RequestLogger())
	e.Use(middleware.Recover())
	e.Use(middleware.RequestID())
	e.Use(CORSMiddleware(cfg))
	e.Use(SecurityMiddleware())
	e.Use(GlobalRateLimiter())
	e.Use(SignatureMiddleware(cfg))
}
