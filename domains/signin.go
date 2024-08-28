package domains

import (
	"auth-service/models"

	"github.com/google/uuid"
)

type SignInRequest struct {
	Email     string `json:"email"`
	Password  string `json:"password"`
	IpAddress string `json:"ip_address"`
	Device    string `json:"device"`
}

type SignInResponse struct {
	ResponseBasic
	Data DataSignInResponse `json:"data"`
}

type DataSignInResponse struct {
	UserID       uuid.UUID `json:"user_id"`
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
}

type SignInUsecase interface {
	ValidateDataRequest(signInRequest SignInRequest) error
	GetUserByEmail(email string) (models.User, error)
	ComparePasswordUser(passwordHash, passwordInput string) error
	CheckAccountVerification(userData models.User, signInRequest SignInRequest) error
	CreateTokenAndDataResponse(userData models.User, signInRequest SignInRequest) (DataSignInResponse, error)
}
