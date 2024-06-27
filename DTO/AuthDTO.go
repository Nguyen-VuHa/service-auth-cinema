package DTO

import (
	"time"

	"github.com/google/uuid"
)

// Type request sign up account
type SignUp_Request struct {
	Email       string `json:"email"` // trả ra JSON với định dạng email thay vì Email
	Password    string `json:"password"`
	FullName    string `json:"full_name"`
	BirthDay    string `json:"birth_day"`
	PhoneNumber string `json:"phone_number"`
	Address     string `json:"address"`
}

// type response từ function trả về trong service Auth
type AuthService_SignUp_Response struct {
}

// Type request sign up account
type SignIn_Request struct {
	Email     string `json:"email"` // trả ra JSON với định dạng email thay vì Email
	Password  string `json:"password"`
	IPAddress string `json:"ip_address"`
	Device    string `json:"device"`
}

// Type request sign up account
type AuthService_SignIn_Response struct {
	UserID       uuid.UUID `json:"u_id"`
	AccessToken  string    `json:"acc_k"`
	RefreshToken string    `json:"ref_k"`
}

// Type request sign up with facebook
type SignInFacebook_Request struct {
	AuthToken string `json:"auth_token"`
}

// Type request sign up with facebook
type AuthService_SignInFacebook_Response struct {
	URL string `json:"url"`
}

type Callback_SignIn_Facebook struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Email       string    `json:"email"`
	AccessToken string    `json:"access_token"`
	TokenType   string    `json:"token_type"`
	Expiry      time.Time `json:"expiry"`
}

// Type request sign up account
type AuthService_Callback_Facebook_Response struct {
	UserID      uuid.UUID `json:"u_id"`
	AccessToken string    `json:"acc_k"`
}
