package utils

import (
	"time"

	"github.com/pquerna/otp/totp"
)

// Hàm tạo mã OTP
func GenerateOTP(secretKey string) (string, string, error) {
	// Tạo khóa bí mật dựa trên userID, device và IP

	// Sinh mã OTP
	otp, err := totp.Generate(totp.GenerateOpts{
		Issuer:      "CinemaBooking",
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
func VerifyOTP(otpCode, secretKey string) bool {
	// Xác minh mã OTP nhập vào
	return totp.Validate(otpCode, secretKey)
}
