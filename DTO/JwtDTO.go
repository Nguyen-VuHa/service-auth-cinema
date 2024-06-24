package DTO

import (
	"time"

	"github.com/google/uuid"
)

type JWTToken struct {
	UserID uuid.UUID
	Exp    time.Duration
}
