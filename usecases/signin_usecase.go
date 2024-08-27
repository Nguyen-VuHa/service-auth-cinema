package usecases

import (
	"auth-service/bootstrap"
	"auth-service/constants"
	"auth-service/domains"
	"auth-service/models"
	"auth-service/utils"
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/go-redis/redis/v8"
)

type signInUsecase struct {
	userRepository domains.UserRepository
	valRepository  domains.ValidateRepository
	ctx            context.Context
	redisUser      *redis.Client
}

func NewSignInUsecase(userRepository domains.UserRepository, valRepository domains.ValidateRepository) domains.SignInUsecase {
	return &signInUsecase{
		userRepository,
		valRepository,
		context.Background(),
		bootstrap.RedisUser,
	}
}

func (su *signInUsecase) ValidateDataRequest(signInRequest domains.SignInRequest) error {
	validations := []func() error{
		func() error { return su.isEmailValid(signInRequest.Email) },
		func() error { return su.isPasswordValid(signInRequest.Password) },
	}

	for _, validation := range validations {
		if err := validation(); err != nil {
			return err
		}
	}

	return nil
}

func (su *signInUsecase) GetUserByEmail(email string) (models.User, error) {
	return su.userRepository.GetByEmail(email)
}

func (su *signInUsecase) ComparePasswordUser(passwordHash, passwordInput string) error {
	errValid := utils.ComparePasswordByBcrypt(passwordHash, passwordInput)

	return errValid
}

func (su *signInUsecase) CheckAccountVerification(userData models.User, signInRequest domains.SignInRequest) error {
	if userData.UserStatus == constants.USER_STATUS_PENDING { // chưa xác thực
		// tạo mã OTP bằng key DEVICE_SECRET_KEY + UserID + Device + IP
		otpData := os.Getenv(constants.DEVICE_SECRET_KEY) + fmt.Sprint(userData.UserID) + signInRequest.Device + signInRequest.IpAddress

		valueRedis, err := su.getDataRedisUser(otpData)

		if err != nil || valueRedis != "" {
			return nil
		}

		otp, secret, err := utils.GenerateOTP(otpData)

		if err != nil {
			return err
		}

		ttl := time.Minute // Thời gian cache là 1 phút
		// lưu redis hash với key là secret và value là true thời gian hết là 60 seconds
		errSaveRedis := su.redisUser.Set(su.ctx, otpData, secret, ttl).Err()

		if errSaveRedis != nil {
			fmt.Println("lỗi lưu data redis", errSaveRedis)
			return errSaveRedis
		}

		// send OTP qua mail của user
		errSentOTP := su.sendOtpViaMail(signInRequest, otp)

		if errSentOTP != nil {
			fmt.Println(errSentOTP)
			return errSentOTP
		}
	}

	return nil
}

// function validation email
func (su *signInUsecase) isEmailValid(email string) error {
	// Kiểm tra Email
	// 1. Email is require
	errIsRequire := su.valRepository.IsRequireString(email)

	if errIsRequire != nil {
		return errIsRequire
	}

	// 2. Email hợp lệ
	errIsEmail := su.valRepository.IsEmail(email)

	if errIsEmail != nil {
		return errIsEmail
	}

	// 3. Độ dài chuỗi cho phép
	errIsLength := su.valRepository.IsMaxLengthString(email, 100) // tối đa 100 ký tự

	if errIsLength != nil {
		return errIsLength
	}

	return nil
}

// function validation password
func (su *signInUsecase) isPasswordValid(password string) error {
	// Kiểm tra password
	// 1. password is require
	errIsRequire := su.valRepository.IsRequireString(password)

	if errIsRequire != nil {
		return errIsRequire
	}

	// 2. kiểm tra range ký tự
	errIsLength := su.valRepository.IsRangeLength(password, 8, 50) // tối thiểu 8 ký tự và tối đa 50 ký tự.

	if errIsLength != nil {
		return errIsLength
	}

	return nil
}

// function get data on Redis
func (su *signInUsecase) getDataRedisUser(keyRedis string) (string, error) {
	val, err := su.redisUser.Get(su.ctx, keyRedis).Result()

	if fmt.Sprint(err) == "redis: nil" {
		return val, nil
	}

	return val, err
}

// function get data on Redis
func (su *signInUsecase) sendOtpViaMail(signInRequest domains.SignInRequest, otp string) error {
	// URL của API bạn muốn gọi
	url := os.Getenv(constants.URL_API_SERVICE) + constants.PATH_SEND_EMAIL_OTP
	params := map[string]interface{}{
		"_mail":   signInRequest.Email,
		"_otp":    otp,
		"_ip":     signInRequest.IpAddress,
		"_device": signInRequest.Device,
		"_secret": os.Getenv(constants.SERVICE_SECRET),
	}

	urlSendOTP, err := utils.AddParamsToURL(url, params)

	if err != nil {
		return err
	}

	// Tạo một request mới
	req, err := http.NewRequest("GET", urlSendOTP, nil)

	if err != nil {
		return err
	}

	// Tạo một HTTP client
	client := &http.Client{}

	// Thực hiện request
	response, err := client.Do(req)

	if err != nil {
		fmt.Println("Error making GET request:", err)
		return err
	}

	defer response.Body.Close()

	// Đọc toàn bộ nội dung của response body
	body, err := io.ReadAll(response.Body)

	if err != nil {
		fmt.Println("Error reading response body:", err)
		return err
	}

	// In ra nội dung response
	fmt.Println("Response Body:", string(body))

	return nil
}
