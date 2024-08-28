package domains

import (
	"time"

	"github.com/google/uuid"
)

type JWTToken struct {
	UserID uuid.UUID
	Exp    time.Duration
}

type RedisDataJWT struct {
	Device       string `json:"device"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}
