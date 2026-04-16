package handler

import (
	"net/http"

	"github.com/codebayu/account-service/common/response"
	"github.com/codebayu/account-service/internal/dto"
	"github.com/codebayu/account-service/internal/service"
	"github.com/labstack/echo/v5"
)

type AuthHandler struct {
	authService service.AuthService
}

func NewAuthHandler(authService service.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

// Register godoc
// @Summary      Register a new user
// @Description  Create a new user account with email and password
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        X-Signature  header    string  true  "Digital Signature"
// @Param        X-Datetime   header    string  true  "Unix Timestamp"
// @Param        X-Channel    header    string  true  "Channel ID (WEB/MOBILE)"
// @Param        request      body      dto.RegisterRequest  true  "Register Request Body"
// @Success      201  {object}  response.Response{data=dto.AuthResponseData}
// @Failure      400  {object}  response.Response
// @Failure      500  {object}  response.Response
// @Router       /auth/register [post]
// @Security     SignatureAuth
// @Security     DatetimeAuth
// @Security     ChannelAuth
func (h *AuthHandler) Register(c *echo.Context) error {
	var req dto.RegisterRequest
	if appErr := validateRequest(c, &req); appErr != nil {
		return (*c).JSON(appErr.StatusCode, response.Error(appErr.StatusCode, appErr.Code, appErr.Message, appErr.Errors))
	}

	data, err := h.authService.Register(req)
	if err != nil {
		return handleServiceError(c, err)
	}

	return (*c).JSON(http.StatusCreated, response.Success(http.StatusCreated, 201000, "created", data))
}

// Login godoc
// @Summary      Login as existing user
// @Description  Authenticate user and return access & refresh tokens
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        X-Signature  header    string  true  "Digital Signature"
// @Param        X-Datetime   header    string  true  "Unix Timestamp"
// @Param        X-Channel    header    string  true  "Channel ID (WEB/MOBILE)"
// @Param        request      body      dto.LoginRequest  true  "Login Request Body"
// @Success      200  {object}  response.Response{data=dto.AuthResponseData}
// @Failure      400  {object}  response.Response
// @Failure      401  {object}  response.Response
// @Failure      500  {object}  response.Response
// @Router       /auth/login [post]
// @Security     SignatureAuth
// @Security     DatetimeAuth
// @Security     ChannelAuth
func (h *AuthHandler) Login(c *echo.Context) error {
	var req dto.LoginRequest
	if appErr := validateRequest(c, &req); appErr != nil {
		return (*c).JSON(appErr.StatusCode, response.Error(appErr.StatusCode, appErr.Code, appErr.Message, appErr.Errors))
	}

	data, err := h.authService.Login(req)
	if err != nil {
		return handleServiceError(c, err)
	}

	return (*c).JSON(http.StatusOK, response.Success(http.StatusOK, 200000, "ok", data))
}
