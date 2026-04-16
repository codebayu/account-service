package handlers

import (
	"net/http"

	"github.com/labstack/echo/v5"
)

func (h *Handler) HealthCheck(c *echo.Context) error {
	return c.JSON(http.StatusOK, map[string]string{
		"status": "ok",
	})
}
