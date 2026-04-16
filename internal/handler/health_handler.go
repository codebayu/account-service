package handler

import (
	"net/http"

	"github.com/labstack/echo/v5"
)

type HealthHandler struct{}

func NewHealthHandler() *HealthHandler {
	return &HealthHandler{}
}

// HealthCheck godoc
// @Summary      Health check
// @Description  Check if the server is up and running
// @Tags         health
// @Accept       json
// @Produce      json
// @Param        X-Signature  header    string  true  "Digital Signature"
// @Param        X-Datetime   header    string  true  "Unix Timestamp"
// @Param        X-Channel    header    string  true  "Channel ID (WEB/MOBILE)"
// @Success      200  {object}  map[string]string
// @Router       /health [get]
// @Security     SignatureAuth
// @Security     DatetimeAuth
// @Security     ChannelAuth
func (h *HealthHandler) HealthCheck(c *echo.Context) error {
	return (*c).JSON(http.StatusOK, map[string]string{
		"status": "ok",
	})
}
