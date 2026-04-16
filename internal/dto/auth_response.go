package dto

import "time"

type AuthResponseData struct {
	AccessToken        string    `json:"accessToken"`
	AccessTokenExpire  time.Time `json:"accessTokenExpire"`
	RefreshToken       string    `json:"refreshToken"`
	RefreshTokenExpire time.Time `json:"refreshTokenExpire"`
}
