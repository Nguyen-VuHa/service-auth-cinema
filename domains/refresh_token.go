package domains

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
}
