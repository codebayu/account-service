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

// GetCurrentUser godoc
// @Summary      Get current user profile
// @Description  Get information of the currently authenticated user
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        Authorization header    string  true  "Bearer {token}"
// @Param        X-Signature   header    string  true  "Digital Signature"
// @Param        X-Datetime    header    string  true  "Unix Timestamp"
// @Param        X-Channel     header    string  true  "Channel ID (WEB/MOBILE)"
// @Success      200  {object}  response.Response{data=dto.UserResponse}
// @Failure      401  {object}  response.Response
// @Router       /user/current [get]
// @Security     BearerAuth
// @Security     SignatureAuth
// @Security     DatetimeAuth
// @Security     ChannelAuth
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
