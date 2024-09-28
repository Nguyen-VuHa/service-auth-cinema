package usecases

import (
	"auth-service/constants"
	"auth-service/domains"
	"auth-service/internal/jwt_util"
	"auth-service/utils"
	"context"
	"fmt"
	"os"
	"time"
)

type signInUsecase struct {
	userRepository  domains.UserRepository
	valRepository   domains.ValidateRepository
	redisRepository domains.RedisRepository
	serviceMailRepo domains.ServiceMailRepository
	ctx             context.Context
}

func NewSignInUsecase(
	userRepository domains.UserRepository, valRepository domains.ValidateRepository,
	redisRepository domains.RedisRepository, serviceMailRepository domains.ServiceMailRepository,
) domains.SignInUsecase {
	return &signInUsecase{
		userRepository,
		valRepository,
		redisRepository,
		serviceMailRepository,
		context.Background(),
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

func (su *signInUsecase) GetUserByEmail(email string) (domains.UserDTO, string, error) { // string là password lấy ra để so sánh khi đăng nhập
	var userData domains.UserDTO

	userModel, err := su.userRepository.GetByEmailPreload(email, "LoginMethod", "Profiles")

	if err != nil {
		return userData, "", err
	}

	userData.UserID = fmt.Sprint(userModel.UserID)
	userData.Email = userModel.Email
	userData.UserStatus = string(userModel.UserStatus)
	userData.LoginMethod = userModel.LoginMethod.LoginMethod
	userData.LoginMethodID = userModel.LoginMethodID
	userData.CreatedAt = userModel.CreatedAt
	userData.UpdatedAt = userModel.UpdatedAt

	for _, profile := range userModel.Profiles {
		switch profile.ProfileKey {
		case "full_name":
			userData.FullName = profile.ProfileValue
		case "birth_day":
			userData.BirthDay = profile.ProfileValue
		case "phone_number":
			userData.PhoneNumber = profile.ProfileValue
		}
	}

	return userData, userModel.Password, nil
}

func (su *signInUsecase) ComparePasswordUser(passwordHash, passwordInput string) error {
	errValid := utils.ComparePasswordByBcrypt(passwordHash, passwordInput)

	return errValid
}

func (su *signInUsecase) CheckAccountVerification(userData domains.UserDTO, signInRequest domains.SignInRequest) error {
	if userData.UserStatus == constants.USER_STATUS_PENDING { // chưa xác thực
		// tạo mã OTP bằng key UserID + Device
		otpData := userData.UserID + signInRequest.Device

		valueRedis, err := su.redisRepository.RedisUserGet(otpData)

		if (fmt.Sprint(err) != "redis: nil" && err != nil) || valueRedis != "" {
			return nil
		}

		otp, secret, err := utils.GenerateOTP(otpData)

		if err != nil {
			return err
		}

		ttl := time.Minute // Thời gian cache là 1 phút
		// lưu redis hash với key là secret và value là true thời gian hết là 60 seconds
		errSaveRedis := su.redisRepository.RedisUserSet(otpData, secret, ttl)

		if errSaveRedis != nil {
			fmt.Println("lỗi lưu data redis", errSaveRedis)
			return errSaveRedis
		}

		// fmt.Println(otp)
		// send OTP qua mail của user
		sendOTPParams := map[string]interface{}{
			"_mail":   signInRequest.Email,
			"_otp":    otp,
			"_ip":     signInRequest.IpAddress,
			"_device": signInRequest.Device,
			"_secret": os.Getenv(constants.SERVICE_SECRET),
		}

		errSentOTP := su.serviceMailRepo.SendOTPCodeToMail(sendOTPParams)

		if errSentOTP != nil {
			fmt.Println(errSentOTP)
			return errSentOTP
		}
	}

	return nil
}

func (su *signInUsecase) CreateTokenAndDataResponse(userData domains.UserDTO, signInRequest domains.SignInRequest) (domains.DataSignInResponse, error) {
	var dataReponse domains.DataSignInResponse

	// set user_id cho response
	dataReponse.UserID = userData.UserID

	if userData.UserStatus == constants.USER_STATUS_ACTIVE {
		// tạo token và thông tin user trả về cho người dùng
		var tokenData domains.JWTToken
		tokenData.UserID = userData.UserID

		// lấy secret key trong .env
		signAccess := os.Getenv(constants.JWT_ACCESS_SECRET)
		signRefresh := os.Getenv(constants.JWT_REFRESH_SECRET)

		// divice + hash secret - để kiểm tra người dùng phải sử dụng đúng thiết bị mà user đăng nhập.
		signHashAccess := signAccess + signInRequest.Device
		signHashRefresh := signRefresh + signInRequest.Device

		accessToken, err := jwt_util.CreateAccessToken(tokenData, signHashAccess)

		if err != nil { // create Token failed.
			return dataReponse, err
		}

		refreshToken, err := jwt_util.CreateRefreshToken(tokenData, signHashRefresh)

		if err != nil { // create Token failed.
			return dataReponse, err
		}

		// set  token cho response
		dataReponse.AccessToken = accessToken
		dataReponse.RefreshToken = refreshToken

		userDataMapString := utils.StructureToMapString(userData)

		// lưu trữ thông tin user trên Redis với key là user_id
		timeToLiveUserData := time.Hour * 24 // Thời gian cache là 1 ngày
		su.redisRepository.RedisUserHMSet(userData.UserID, userDataMapString, timeToLiveUserData)
		// lưu trữ thông tin user trên Redis với key là email để cache dữ liệu
		su.redisRepository.RedisUserHMSet(userData.Email, userDataMapString, timeToLiveUserData)
		// tạo redis key bằng keyDEVICE_SECRET_KEY + UserID + Device
		signRedis := os.Getenv(constants.DEVICE_SECRET_KEY) + userData.UserID + signInRequest.Device

		var redisDataJWT domains.RedisDataJWT

		redisDataJWT.Device = signInRequest.Device
		redisDataJWT.AccessToken = accessToken
		redisDataJWT.RefreshToken = refreshToken

		redisDataJWTMapString := utils.StructureToMapString(redisDataJWT)
		timeToLiveJWTData := time.Hour * 24 * 30 // Thời gian cache là 30 ngày = thời gian hết hạn của refresh token
		su.redisRepository.RedisAuthHMSet(signRedis, redisDataJWTMapString, timeToLiveJWTData)
	}

	return dataReponse, nil
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
