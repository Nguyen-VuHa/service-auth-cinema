package usecases

import (
	"auth-service/constants"
	"auth-service/domains"
	"auth-service/internal/jwt_util"
	"auth-service/models"
	"auth-service/utils"
	"errors"
	"fmt"
	"os"
	"time"
)

type verifyOTPUsecase struct {
	userRepo  domains.UserRepository
	redisRepo domains.RedisRepository
}

func NewVerifyOTPUsecase(userRepo domains.UserRepository, redisRepo domains.RedisRepository) domains.VerifyOTPUsecase {
	return &verifyOTPUsecase{
		userRepo,
		redisRepo,
	}
}

// Validate data (email có thuộc hệ thống và chưa active)
func (votp *verifyOTPUsecase) ValidateEmail(email string) (models.User, bool, error) {
	user_data, err := votp.userRepo.GetByEmail(email)

	if err != nil {
		return user_data, false, err
	}

	if user_data.UserStatus != constants.USER_STATUS_PENDING {
		return user_data, false, errors.New("tài khoản đã được xác thực")
	}

	return user_data, true, nil
}

// checking OTP gửi lên
func (votp *verifyOTPUsecase) CheckOTPValid(data_request domains.VerifyOTPRequest, user_data models.User) (bool, error) {
	otpData := os.Getenv(constants.DEVICE_SECRET_KEY) + fmt.Sprint(user_data.UserID) + data_request.Device

	fmt.Println(otpData)

	secret_otp, err := votp.redisRepo.RedisUserGet(otpData)

	if err != nil {
		return false, err
	}

	// checking otp
	is_valid_otp := utils.VerifyOTP(data_request.OTP, secret_otp)

	return is_valid_otp, nil
}

// check OTP thành công tiến hành cập nhật thông tin user tạo token gửi về cho user.
func (votp *verifyOTPUsecase) UpdateUser(data_request domains.VerifyOTPRequest, user_data models.User) (domains.VerifyOTPDataResponse, error) {
	var data_response domains.VerifyOTPDataResponse

	// thay đổi trạng thái account user
	user_data.UserStatus = constants.USER_STATUS_ACTIVE

	err := votp.userRepo.Update(&user_data)

	if err != nil {
		return data_response, err
	}

	// tạo token và thông tin user trả về cho người dùng
	var tokenData domains.JWTToken
	tokenData.UserID = fmt.Sprint(user_data.UserID)

	// lấy secret key trong .env
	signAccess := os.Getenv(constants.JWT_REFRESH_SECRET)
	signRefresh := os.Getenv(constants.JWT_REFRESH_SECRET)

	// divice + hash secret - để kiểm tra người dùng phải sử dụng đúng thiết bị mà user đăng nhập.
	signHashAccess := signAccess + data_request.Device
	signHashRefresh := signRefresh + data_request.Device

	accessToken, err := jwt_util.CreateAccessToken(tokenData, signHashAccess)

	if err != nil { // create Token failed.
		return data_response, err
	}

	refreshToken, err := jwt_util.CreateRefreshToken(tokenData, signHashRefresh)

	if err != nil { // create Token failed.
		return data_response, err
	}

	// set  token cho response
	data_response.AccessToken = accessToken
	data_response.RefreshToken = refreshToken
	data_response.UserID = fmt.Sprint(user_data.UserID)

	userDataMapString := utils.StructureToMapString(user_data)

	// lưu trữ thông tin user trên Redis với key là user_id
	timeToLiveUserData := time.Hour * 24 // Thời gian cache là 1 ngày
	votp.redisRepo.RedisUserHMSet(fmt.Sprint(user_data.UserID), userDataMapString, timeToLiveUserData)
	// lưu trữ thông tin user trên Redis với key là email để cache dữ liệu
	votp.redisRepo.RedisUserHMSet(fmt.Sprint(user_data.Email), userDataMapString, timeToLiveUserData)

	// tạo redis key bằng keyDEVICE_SECRET_KEY + UserID + Device
	signRedis := os.Getenv(constants.DEVICE_SECRET_KEY) + fmt.Sprint(user_data.UserID) + data_request.Device

	var redisDataJWT domains.RedisDataJWT

	redisDataJWT.Device = data_request.Device
	redisDataJWT.AccessToken = accessToken
	redisDataJWT.RefreshToken = refreshToken

	redisDataJWTMapString := utils.StructureToMapString(redisDataJWT)
	timeToLiveJWTData := time.Hour * 24 * 30 // Thời gian cache là 30 ngày = thời gian hết hạn của refresh token
	votp.redisRepo.RedisAuthHMSet(signRedis, redisDataJWTMapString, timeToLiveJWTData)

	return data_response, nil
}
