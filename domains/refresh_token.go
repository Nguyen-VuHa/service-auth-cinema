package domains

import "time"

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token"`
	UserID       string `json:"user_id"`
	Device       string `json:"device"`
}

type RefreshTokenResponse struct {
	ResponseBasic
	Data string `json:"data"`
}

type RefreshTokenUsecase interface {
	ValidateRefreshToken(data RefreshTokenRequest) error
	CreateRefreshToken(data RefreshTokenRequest) (string, error)
	ValidateTokenGoogle(user_id, access_token string) (bool, error)
	CreateRefreshTokenGoogle(user_id string) (string, error)
}

type DataTokenGoogle struct {
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	ExpiredTime  time.Time `json:"expired_time"`
}
