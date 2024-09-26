package domains

import (
	"time"

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
	UserID       string `json:"user_id"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type SignInUsecase interface {
	ValidateDataRequest(signInRequest SignInRequest) error
	GetUserByEmail(email string) (UserDTO, string, error) // string là password lấy ra để so sánh với password nhập vào
	ComparePasswordUser(passwordHash, passwordInput string) error
	CheckAccountVerification(userData UserDTO, signInRequest SignInRequest) error
	CreateTokenAndDataResponse(userData UserDTO, signInRequest SignInRequest) (DataSignInResponse, error)
}

type SignInFacebookUsecase interface {
	AuthFacebookURL() (string, error)
}

type SignInGoogleUsecase interface {
	AuthGoogleURL() (string, error)
}

type CallbackFacebookUsecase interface {
	GetDetailUserWithCodeFacebook(code string) (DataCallbackSignInFacebook, error)
	CreateUserLoginWithFacebook(data DataCallbackSignInFacebook) (DataSignInFacebookResponse, error)
}

type CallbackSignInFacebookResponse struct {
	ResponseBasic
	Data DataSignInFacebookResponse `json:"data"`
}

type DataSignInFacebookResponse struct {
	UserID      uuid.UUID `json:"uid"`
	AccessToken string    `json:"access_token"`
	Method      uint      `json:"method"`
}

type DataCallbackSignInFacebook struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Email       string    `json:"email"`
	AccessToken string    `json:"accessToken"`
	TokenType   string    `json:"tokenType"`
	Expiry      time.Time `json:"expiry"`
}
