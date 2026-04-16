package middleware

import (
	"net/http"
	"strings"

	"github.com/codebayu/account-service/common/response"
	"github.com/codebayu/account-service/internal/utils"
	"github.com/labstack/echo/v5"
)

func AuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
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
			return c.JSON(http.StatusUnauthorized, response.Error(http.StatusUnauthorized, 401000, "unauthorized", []string{"invalid token type"}))
		}

		// Set user_uuid to context
		c.Set("user_uuid", claims.UUID)

		return next(c)
	}
}
