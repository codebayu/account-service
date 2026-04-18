package middleware

import (
	"net/http"
	"strings"

	"github.com/codebayu/account-service/common/response"
	"github.com/codebayu/account-service/internal/repository"
	"github.com/codebayu/account-service/internal/utils"
	"github.com/labstack/echo/v5"
)

func AuthMiddleware(tokenRepo repository.TokenRepository) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c *echo.Context) error {
		authHeader := c.Request().Header.Get("Authorization")
		if authHeader == "" {
			return c.JSON(http.StatusUnauthorized, response.Error(http.StatusUnauthorized, 401000, "unauthorized", []string{"missing authorization header"}))
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			return c.JSON(http.StatusUnauthorized, response.Error(http.StatusUnauthorized, 401000, "unauthorized", []string{"invalid authorization header format"}))
		}

		tokenString := parts[1]
		claims, err := utils.ParseToken(tokenString)
		if err != nil {
			return c.JSON(http.StatusUnauthorized, response.Error(http.StatusUnauthorized, 401000, "unauthorized", []string{err.Error()}))
		}

		if claims.IsRefresh {
			return (*c).JSON(http.StatusUnauthorized, response.Error(http.StatusUnauthorized, 401000, "unauthorized", []string{"invalid token type"}))
		}

		ctx := (*c).Request().Context()
		isBlacklisted, err := tokenRepo.IsAccessTokenBlacklisted(ctx, claims.ID)
		if err != nil {
			return (*c).JSON(http.StatusInternalServerError, response.Error(http.StatusInternalServerError, 500000, "internal server error", []string{"failed to verify token status"}))
		}
		if isBlacklisted {
			return (*c).JSON(http.StatusUnauthorized, response.Error(http.StatusUnauthorized, 401000, "unauthorized", []string{"token has been revoked"}))
		}

		// Set user_uuid to context
		c.Set("user_uuid", claims.UUID)

		return next(c)
		}
	}
}
