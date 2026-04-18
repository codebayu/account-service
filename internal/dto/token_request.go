package dto

// RefreshTokenRequest is the request body for the refresh token endpoint.
// swagger:model
type RefreshTokenRequest struct {
	RefreshToken string `json:"refreshToken" validate:"required"`
}

// LogoutRequest is the request body for the logout endpoint.
// swagger:model
type LogoutRequest struct {
	RefreshToken string `json:"refreshToken" validate:"required"`
}
