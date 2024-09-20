package domains

import "auth-service/models"

type VerifyOTPUsecase interface {
	ValidateEmail(email string) (models.User, bool, error)
	CheckOTPValid(data_request VerifyOTPRequest, user_data models.User) (bool, error)
	UpdateUser(data_request VerifyOTPRequest, user_data models.User) (VerifyOTPDataResponse, error)
}

type VerifyOTPRequest struct {
	OTP       string `json:"otp"`
	Email     string `json:"email"`
	IpAddress string `json:"ip_address"`
	Device    string `json:"device"`
}

type VerifyOTPResponse struct {
	ResponseBasic
	Data VerifyOTPDataResponse `json:"data"`
}

type VerifyOTPDataResponse struct {
	UserID       string `json:"user_id"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}
