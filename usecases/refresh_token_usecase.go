package usecases

import (
	"auth-service/constants"
	"auth-service/domains"
	"auth-service/internal/jwt_util"
	"auth-service/utils"
	"errors"
	"fmt"
	"os"
	"time"

	"golang.org/x/oauth2"
)

type refreshTokenUsecase struct {
	redisRepo      domains.RedisRepository
	thirdPartyRepo domains.ThirdPartyRepository
	googleConfig   *oauth2.Config
}

func NewRefreshTokenUsecase(
	redisRepo domains.RedisRepository,
	thirdPartyRepo domains.ThirdPartyRepository,
	googleConfig *oauth2.Config,
) domains.RefreshTokenUsecase {
	return &refreshTokenUsecase{
		redisRepo,
		thirdPartyRepo,
		googleConfig,
	}
}

func (rfu *refreshTokenUsecase) ValidateRefreshToken(data domains.RefreshTokenRequest) error {
	// tạo redis key bằng keyDEVICE_SECRET_KEY + UserID + Device
	signRedis := os.Getenv(constants.DEVICE_SECRET_KEY) + data.UserID + data.Device

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
	sign_redis := os.Getenv(constants.DEVICE_SECRET_KEY) + data.UserID + data.Device

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

func (rfu *refreshTokenUsecase) ValidateTokenGoogle(user_id, access_token string) (bool, error) {
	// Kiểm tra xem có đúng user đó không
	fields := []string{"AccessToken"}
	data_redis, err := rfu.redisRepo.RedisAuthHMGetFields(user_id, fields)

	if err != nil {
		return false, err
	}

	if data_redis["AccessToken"] != access_token {
		return false, errors.New("token mismatch")
	}

	return true, nil
}

func (rfu *refreshTokenUsecase) CreateRefreshTokenGoogle(user_id string) (string, error) {
	// tạo token và thông tin user trả về cho người dùng
	var tokenData domains.JWTToken

	tokenData.UserID = fmt.Sprint(user_id)
	// lấy secret key trong .env
	sign_access := os.Getenv(constants.JWT_ACCESS_SECRET)
	access_token, err := jwt_util.CreateAccessToken(tokenData, sign_access)

	if err != nil { // create Token failed.
		return "", err
	}

	var redisDataJWT domains.RedisDataJWT

	redisDataJWT.AccessToken = access_token
	redisDataJWTMapString := utils.StructureToMapString(redisDataJWT)

	timeToLiveJWTData := time.Hour * 24 * 30 // Thời gian cache là 30 ngay
	rfu.redisRepo.RedisAuthHMSet(fmt.Sprint(user_id), redisDataJWTMapString, timeToLiveJWTData)

	return access_token, nil
}
