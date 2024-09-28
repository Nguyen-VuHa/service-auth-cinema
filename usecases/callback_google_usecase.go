package usecases

import (
	"auth-service/constants"
	"auth-service/domains"
	"auth-service/internal/jwt_util"
	"auth-service/models"
	"auth-service/utils"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"golang.org/x/oauth2"
)

type callbackGoogleUsecase struct {
	googleConfig    *oauth2.Config
	userRepo        domains.UserRepository
	userProfileRepo domains.UserProfileRepository
	thirdPartyRepo  domains.ThirdPartyRepository
	redisRepo       domains.RedisRepository
}

func NewCallbackGoogleUsecase(
	googleConfig *oauth2.Config,
	userRepo domains.UserRepository,
	userProfileRepo domains.UserProfileRepository,
	thirdPartyRepo domains.ThirdPartyRepository,
	redisRepo domains.RedisRepository,
) domains.CallbackGoogleUsecase {
	return &callbackGoogleUsecase{
		googleConfig,
		userRepo,
		userProfileRepo,
		thirdPartyRepo,
		redisRepo,
	}
}

func (cgu *callbackGoogleUsecase) GetDetailUserWithCodeGoogle(code string) (domains.DataCallbackSignInGoogle, error) {
	var data_google domains.DataCallbackSignInGoogle

	token, err := cgu.googleConfig.Exchange(context.Background(), code)

	if err != nil {
		fmt.Println(err)
		return data_google, err
	}

	// Lấy id_token từ Extra
	id_token, ok := token.Extra("id_token").(string)

	if !ok {
		return data_google, fmt.Errorf("id_token not found in token response")
	}

	resp, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + token.AccessToken)

	if err != nil {
		return data_google, err
	}

	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	// Giải mã JSON từ []byte sang đối tượng User
	if err := json.Unmarshal(body, &data_google); err != nil {
		return data_google, err
	}

	data_google.AccessToken = id_token
	data_google.RefreshToken = token.RefreshToken
	data_google.Expiry = token.Expiry

	return data_google, nil
}

func (cgu *callbackGoogleUsecase) CreateUserLoginWithGoogle(data domains.DataCallbackSignInGoogle) (domains.DataSignInGoogleResponse, error) {
	var data_response domains.DataSignInGoogleResponse

	email_google := data.ID + "@google.com"
	// kiểm tra email google tồn tại chưa
	user_google, err := cgu.userRepo.GetByEmail(email_google)
	data_response.Method = constants.LOGIN_GOOGLE_ID

	if err == nil {
		// tạo token và thông tin user trả về cho người dùng
		var tokenData domains.JWTToken

		tokenData.UserID = fmt.Sprint(user_google.UserID)
		// lấy secret key trong .env
		sign_access := os.Getenv(constants.JWT_ACCESS_SECRET)
		access_token, err := jwt_util.CreateAccessToken(tokenData, sign_access)

		if err != nil { // create Token failed.
			return data_response, err
		}

		data_response.AccessToken = access_token
		data_response.UserID = user_google.UserID

		var redisDataJWT domains.RedisDataJWT

		redisDataJWT.AccessToken = access_token
		redisDataJWTMapString := utils.StructureToMapString(redisDataJWT)

		timeToLiveJWTData := time.Hour * 24 * 30 // Thời gian cache là 30 ngày = thời gian hết hạn của refresh token
		cgu.redisRepo.RedisAuthHMSet(fmt.Sprint(user_google.UserID), redisDataJWTMapString, timeToLiveJWTData)

		return data_response, nil
	}

	var tokenData domains.JWTToken
	tokenData.UserID = fmt.Sprint(user_google.UserID)
	// lấy secret key trong .env
	sign_access := os.Getenv(constants.JWT_ACCESS_SECRET)

	access_token, err := jwt_util.CreateAccessToken(tokenData, sign_access)

	if err != nil { // create Token failed.
		return data_response, err
	}

	data_response.AccessToken = access_token

	var user_google_create models.User

	user_google_create.Email = email_google
	// thêm một số trường với rule khi tạo mới tài khoản
	user_google_create.UserStatus = constants.USER_STATUS_ACTIVE
	user_google_create.LoginMethodID = constants.LOGIN_GOOGLE_ID

	// create user
	err = cgu.userRepo.Create(&user_google_create)

	if err != nil {
		return data_response, err
	}

	// set data cho barng profile
	var user_google_profile = make(map[string]interface{})
	user_google_profile[constants.USER_PROFILE_FULLNAME] = data.Name
	user_google_profile[constants.USER_PROFILE_EMAIL_PLATFORM] = data.Email

	var profile_keys = []string{constants.USER_PROFILE_FULLNAME, constants.USER_PROFILE_EMAIL_PLATFORM}

	for _, key := range profile_keys { // 1 số biến object cần lưu vào user profile (FullName)
		var userProfileData models.UserProfile // Khai báo biến để chứa thông tin detail user hợp lệ
		userProfileData.ProfileKey = key
		userProfileData.ProfileValue = user_google_profile[key].(string)
		userProfileData.UserID = user_google_create.UserID // Gán UserID khoá ngoại trong UserProfile

		err := cgu.userProfileRepo.Create(&userProfileData)

		if err != nil {
			fmt.Println(err)
		}
	}

	// insert thông tin vào auth third party
	var authThirdParty models.AuthThirdParty

	authThirdParty.AccessToken = data.AccessToken
	authThirdParty.RefreshToken = data.RefreshToken
	authThirdParty.ProviderID = data.ID
	authThirdParty.Provider = "Google"
	authThirdParty.ExpiredTime = data.Expiry
	authThirdParty.UserID = user_google_create.UserID

	err = cgu.thirdPartyRepo.Create(&authThirdParty)

	if err != nil {
		return data_response, err
	}

	data_response.UserID = user_google_create.UserID

	var redisDataJWT domains.RedisDataJWT

	redisDataJWT.AccessToken = access_token
	redisDataJWTMapString := utils.StructureToMapString(redisDataJWT)

	timeToLiveJWTData := time.Hour * 24 * 30 // Thời gian cache là 30 ngày = thời gian hết hạn của refresh token
	cgu.redisRepo.RedisAuthHMSet(fmt.Sprint(user_google_create.UserID), redisDataJWTMapString, timeToLiveJWTData)

	return data_response, nil
}
