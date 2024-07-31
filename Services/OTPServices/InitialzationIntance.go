package otp_services

type OTPService struct {
}

// khởi tạo intance NewAuthService định nghĩa struct AuthService
func NewOTPService() *OTPService {
	return &OTPService{}
}

func (otpService *OTPService) SendOTPViaMail(cacheKey, email, otp string) error {

	return nil
}
