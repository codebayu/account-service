package handlers

import (
	"net/http"

	"github.com/codebayu/account-service/common/response"
	"github.com/labstack/echo/v5"
)

func (h *Handler) GetCurrentUser(c *echo.Context) error {
	userUUID, ok := (*c).Get("user_uuid").(string)
	if !ok {
		return (*c).JSON(http.StatusUnauthorized, response.Error(http.StatusUnauthorized, 401000, "unauthorized", []string{"invalid user context"}))
	}

	user, err := h.UserService.GetProfile(userUUID)
	if err != nil {
		if err.Error() == "user not found" {
			return (*c).JSON(http.StatusNotFound, response.Error(http.StatusNotFound, 404000, "user not found", nil))
		}
		return (*c).JSON(http.StatusInternalServerError, response.Error(http.StatusInternalServerError, 500000, "internal server error", []string{err.Error()}))
	}

	return (*c).JSON(http.StatusOK, response.Success(http.StatusOK, 200000, "ok", user))
}
