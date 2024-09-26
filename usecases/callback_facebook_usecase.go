package usecases

import (
	"auth-service/constants"
	"auth-service/domains"
	"auth-service/models"
	"encoding/json"
	"fmt"

	"golang.org/x/oauth2"
)

type callbackFacebookUsecase struct {
	facebookConfig  *oauth2.Config
	userRepo        domains.UserRepository
	userProfileRepo domains.UserProfileRepository
	thirdPartyRepo  domains.ThirdPartyRepository
}

func NewCallbackFacebookUsecase(
	facebookConfig *oauth2.Config,
	useRepo domains.UserRepository,
	userProfileRepo domains.UserProfileRepository,
	thirdPartyRepo domains.ThirdPartyRepository,
) domains.CallbackFacebookUsecase {
	return &callbackFacebookUsecase{
		facebookConfig,
		useRepo,
		userProfileRepo,
		thirdPartyRepo,
	}
}

func (cfu *callbackFacebookUsecase) GetDetailUserWithCodeFacebook(code string) (domains.DataCallbackSignInFacebook, error) {
	var data_facebook domains.DataCallbackSignInFacebook

	token, err := cfu.facebookConfig.Exchange(oauth2.NoContext, code)
	if err != nil {
		// http.Error(w, "Failed to exchange token: "+err.Error(), http.StatusInternalServerError)
		return data_facebook, err
	}

	client := cfu.facebookConfig.Client(oauth2.NoContext, token)
	response, err := client.Get("https://graph.facebook.com/me?fields=id,name,email")
	if err != nil {
		// http.Error(w, "Failed to get user info: "+err.Error(), http.StatusInternalServerError)
		return data_facebook, err
	}

	defer response.Body.Close()

	if err := json.NewDecoder(response.Body).Decode(&data_facebook); err != nil {
		// http.Error(w, "Failed to decode user info: "+err.Error(), http.StatusInternalServerError)
		return data_facebook, err
	}

	data_facebook.AccessToken = token.AccessToken
	data_facebook.TokenType = token.TokenType
	data_facebook.Expiry = token.Expiry

	return data_facebook, nil
}

func (cfu *callbackFacebookUsecase) CreateUserLoginWithFacebook(data domains.DataCallbackSignInFacebook) (domains.DataSignInFacebookResponse, error) {
	var data_response domains.DataSignInFacebookResponse

	email_facebook := data.ID + "@facebook.com"

	// kiểm tra email facebook tồn tại chưa
	user_facebook, err := cfu.userRepo.GetByEmail(email_facebook)

	data_response.AccessToken = data.AccessToken
	data_response.Method = constants.LOGIN_FACEBOOK_ID

	if err == nil {
		data_response.UserID = user_facebook.UserID

		return data_response, nil
	}

	var user_facebook_create models.User

	user_facebook_create.Email = email_facebook
	// thêm một số trường với rule khi tạo mới tài khoản
	user_facebook_create.UserStatus = constants.USER_STATUS_ACTIVE
	user_facebook_create.LoginMethodID = constants.LOGIN_FACEBOOK_ID

	// create user nếu chưa có trong hệ thống
	err = cfu.userRepo.Create(&user_facebook_create)

	if err != nil {
		return data_response, err
	}

	// set data cho barng profile
	var user_facebook_profile = make(map[string]interface{})
	user_facebook_profile[constants.USER_PROFILE_FULLNAME] = data.Name

	var profile_keys = []string{constants.USER_PROFILE_FULLNAME}

	for _, key := range profile_keys { // 1 số biến object cần lưu vào user profile (FullName)

		var userProfileData models.UserProfile // Khai báo biến để chứa thông tin detail user hợp lệ
		userProfileData.ProfileKey = key
		userProfileData.ProfileValue = user_facebook_profile[key].(string)
		userProfileData.UserID = user_facebook_create.UserID // Gán UserID khoá ngoại trong UserProfile

		err := cfu.userProfileRepo.Create(&userProfileData)

		if err != nil {
			fmt.Println(err)
		}
	}

	// insert thông tin vào auth third party
	var authThirdParty models.AuthThirdParty

	authThirdParty.AccessToken = data.AccessToken
	authThirdParty.ProviderID = data.ID
	authThirdParty.Provider = "Facebook"
	authThirdParty.ExpiredTime = data.Expiry
	authThirdParty.UserID = user_facebook_create.UserID

	err = cfu.thirdPartyRepo.Create(&authThirdParty)

	if err != nil {
		return data_response, err
	}

	data_response.UserID = user_facebook_create.UserID

	return data_response, nil
}
