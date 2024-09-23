package usecases

import (
	"auth-service/constants"
	"auth-service/domains"
	"auth-service/internal/jwt_util"
	"errors"
	"fmt"
	"os"
)

type refreshTokenUsecase struct {
	redisRepo domains.RedisRepository
}

func NewRefreshTokenUsecase(redisRepo domains.RedisRepository) domains.RefreshTokenUsecase {
	return &refreshTokenUsecase{
		redisRepo,
	}
}

func (rfu *refreshTokenUsecase) ValidateRefreshToken(data domains.RefreshTokenRequest) error {
	// tạo redis key bằng keyDEVICE_SECRET_KEY + UserID + Device
	signRedis := os.Getenv(constants.DEVICE_SECRET_KEY) + fmt.Sprint(data.UserID) + data.Device

	fields := []string{"RefreshToken"}

	data_redis, err := rfu.redisRepo.RedisAuthHMGetFields(signRedis, fields)

	if err != nil {
		return err
	}

	if data_redis[fields[0]] == nil {
		return errors.New("data is empty")
	}

	sign_refresh := os.Getenv(constants.JWT_REFRESH_SECRET)

	sign_hash_refresh := sign_refresh + data.Device

	_, err = jwt_util.VerifyJWTToken(data_redis[fields[0]].(string), sign_hash_refresh)

	if err != nil {
		return err
	}

	return nil
}

func (rfu *refreshTokenUsecase) CreateRefreshToken(data domains.RefreshTokenRequest) (string, error) {
	// tạo redis key bằng keyDEVICE_SECRET_KEY + UserID + Device
	sign_redis := os.Getenv(constants.DEVICE_SECRET_KEY) + fmt.Sprint(data.UserID) + data.Device

	// tạo token và thông tin user trả về cho người dùng
	var token_data domains.JWTToken
	token_data.UserID = data.UserID

	// lấy secret key trong .env
	sign_access := os.Getenv(constants.JWT_ACCESS_SECRET)
	sign_hash_access := sign_access + data.Device

	access_token, err := jwt_util.CreateAccessToken(token_data, sign_hash_access)

	if err != nil { // create Token failed.
		return "", err
	}

	err = rfu.redisRepo.RedisAuthHSetUpdateField(sign_redis, "AccessToken", access_token)

	if err != nil {
		fmt.Println("Save Redis failed", err)
		return "", err
	}

	return access_token, nil
}
