package handler

import (
	"net/http"

	"github.com/codebayu/account-service/common/response"
	"github.com/codebayu/account-service/internal/service"
	"github.com/labstack/echo/v5"
)

type UserHandler struct {
	userService service.UserService
}

func NewUserHandler(userService service.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

func (h *UserHandler) GetCurrentUser(c *echo.Context) error {
	userUUID, ok := (*c).Get("user_uuid").(string)
	if !ok {
		return (*c).JSON(http.StatusUnauthorized, response.Error(http.StatusUnauthorized, 401000, "unauthorized", []string{"invalid user context"}))
	}

	user, err := h.userService.GetProfile(userUUID)
	if err != nil {
		return handleServiceError(c, err)
	}

	return (*c).JSON(http.StatusOK, response.Success(http.StatusOK, 200000, "ok", user))
}
