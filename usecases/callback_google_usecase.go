package usecases

import (
	"auth-service/constants"
	"auth-service/domains"
	"auth-service/models"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"golang.org/x/oauth2"
)

type callbackGoogleUsecase struct {
	googleConfig    *oauth2.Config
	userRepo        domains.UserRepository
	userProfileRepo domains.UserProfileRepository
	thirdPartyRepo  domains.ThirdPartyRepository
}

func NewCallbackGoogleUsecase(
	googleConfig *oauth2.Config,
	userRepo domains.UserRepository,
	userProfileRepo domains.UserProfileRepository,
	thirdPartyRepo domains.ThirdPartyRepository,
) domains.CallbackGoogleUsecase {
	return &callbackGoogleUsecase{
		googleConfig,
		userRepo,
		userProfileRepo,
		thirdPartyRepo,
	}
}

func (cgu *callbackGoogleUsecase) GetDetailUserWithCodeGoogle(code string) (domains.DataCallbackSignInGoogle, error) {
	var data_google domains.DataCallbackSignInGoogle

	token, err := cgu.googleConfig.Exchange(context.Background(), code)

	if err != nil {
		fmt.Println(err)
		return data_google, err
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

	data_google.AccessToken = token.AccessToken
	data_google.RefreshToken = token.RefreshToken
	data_google.Expiry = token.Expiry

	return data_google, nil
}

func (cgu *callbackGoogleUsecase) CreateUserLoginWithGoogle(data domains.DataCallbackSignInGoogle) (domains.DataSignInGoogleResponse, error) {
	var data_response domains.DataSignInGoogleResponse

	email_google := data.ID + "@google.com"
	// kiểm tra email google tồn tại chưa
	user_google, err := cgu.userRepo.GetByEmail(email_google)

	data_response.AccessToken = data.AccessToken
	data_response.Method = constants.LOGIN_GOOGLE_ID

	if err == nil {
		data_response.UserID = user_google.UserID

		return data_response, nil
	}

	var user_google_create models.User

	user_google_create.Email = email_google
	// thêm một số trường với rule khi tạo mới tài khoản
	user_google_create.UserStatus = constants.USER_STATUS_ACTIVE
	user_google_create.LoginMethodID = constants.LOGIN_GOOGLE_ID

	// create user
	err = cgu.userRepo.Create(&user_google_create)

	fmt.Println(err)
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

	return data_response, nil
}
