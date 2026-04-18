package handler

import (
	"net/http"
	"strings"

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

// RefreshToken godoc
// @Summary      Refresh an access token
// @Description  Get a new access token and refresh token via refresh token rotation
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        X-Signature  header    string  true  "Digital Signature"
// @Param        X-Datetime   header    string  true  "Unix Timestamp"
// @Param        X-Channel    header    string  true  "Channel ID (WEB/MOBILE)"
// @Param        request      body      dto.RefreshTokenRequest  true  "Refresh Token Request"
// @Success      200  {object}  response.Response{data=dto.AuthResponseData}
// @Failure      400  {object}  response.Response
// @Failure      401  {object}  response.Response
// @Router       /auth/refresh-token [post]
// @Security     SignatureAuth
// @Security     DatetimeAuth
// @Security     ChannelAuth
func (h *AuthHandler) RefreshToken(c *echo.Context) error {
	var body struct {
		RefreshToken string `json:"refreshToken" validate:"required"`
	}

	if err := (*c).Bind(&body); err != nil {
		return (*c).JSON(http.StatusBadRequest, response.Error(http.StatusBadRequest, 400000, "invalid body", nil))
	}
	if body.RefreshToken == "" {
		return (*c).JSON(http.StatusBadRequest, response.Error(http.StatusBadRequest, 400000, "refreshToken is required", nil))
	}

	ctx := (*c).Request().Context()
	data, err := h.authService.RefreshToken(ctx, body.RefreshToken)
	if err != nil {
		return handleServiceError(c, err)
	}

	return (*c).JSON(http.StatusOK, response.Success(http.StatusOK, 200000, "ok", data))
}

// Logout godoc
// @Summary      Logout user
// @Description  Logout user by revoking their refresh token
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        X-Signature  header    string  true  "Digital Signature"
// @Param        X-Datetime   header    string  true  "Unix Timestamp"
// @Param        X-Channel    header    string  true  "Channel ID (WEB/MOBILE)"
// @Param        request      body      dto.LogoutRequest  true  "Logout Request"
// @Success      200  {object}  response.Response
// @Failure      400  {object}  response.Response
// @Failure      401  {object}  response.Response
// @Router       /auth/logout [post]
// @Security     BearerAuth
// @Security     SignatureAuth
// @Security     DatetimeAuth
// @Security     ChannelAuth
func (h *AuthHandler) Logout(c *echo.Context) error {
	var body struct {
		RefreshToken string `json:"refreshToken" validate:"required"`
	}

	if err := (*c).Bind(&body); err != nil {
		return (*c).JSON(http.StatusBadRequest, response.Error(http.StatusBadRequest, 400000, "invalid body", nil))
	}
	if body.RefreshToken == "" {
		return (*c).JSON(http.StatusBadRequest, response.Error(http.StatusBadRequest, 400000, "refreshToken is required", nil))
	}

	authHeader := (*c).Request().Header.Get("Authorization")
	var accessToken string
	if authHeader != "" {
		parts := strings.Split(authHeader, " ")
		if len(parts) == 2 && parts[0] == "Bearer" {
			accessToken = parts[1]
		}
	}

	ctx := (*c).Request().Context()
	err := h.authService.Logout(ctx, body.RefreshToken, accessToken)
	if err != nil {
		return handleServiceError(c, err)
	}

	return (*c).JSON(http.StatusOK, response.Success(http.StatusOK, 200000, "logged out successfully", nil))
}
