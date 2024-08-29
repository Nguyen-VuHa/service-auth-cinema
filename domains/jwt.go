package domains

import (
	"time"
)

type JWTToken struct {
	UserID string
	Exp    time.Duration
}

type RedisDataJWT struct {
	Device       string `json:"device"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}
