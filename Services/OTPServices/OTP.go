package otp_services

import (
	"time"

	"github.com/pquerna/otp/totp"
)

func (otpService *OTPService) GenerationOTP(secretKey string) (string, string, error) {
	// secretKey dựa trên email, device và IP

	// Sinh mã OTP
	otp, err := totp.Generate(totp.GenerateOpts{
		Issuer:      "Cinema Auth",
		AccountName: secretKey,
		Period:      60,
		SecretSize:  20, // Độ dài khóa bí mật tính bằng byte
	})

	if err != nil {
		return "", "", err
	}

	// Tạo OTP từ khóa bí mật
	otpCode, err := totp.GenerateCode(otp.Secret(), time.Now())
	if err != nil {
		return "", "", err
	}

	return otpCode, otp.Secret(), nil
}

// Hàm xác minh mã OTP
func (otpService *OTPService) AuthorizationOTP(otpCode, secretKey string) bool {
	// Xác minh mã OTP nhập vào
	return totp.Validate(otpCode, secretKey)
}
