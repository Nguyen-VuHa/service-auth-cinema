package bootstrap

import (
	"auth-service/constants"
	"os"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/facebook"
)

var (
	FacebookConfig *oauth2.Config
	GoogleConfig   *oauth2.Config
)

func ConfigFacebookAuth() {
	// get thông tin ứng dụng xác thực facebook trong env
	facebookAppID := os.Getenv(constants.FACEBOOK_APP_ID)
	facebookSecret := os.Getenv(constants.FACEBOOK_SECRET)
	urlCallBack := os.Getenv(constants.URL_CALLBACK)
	urlCallBack += constants.FACEBOOK_PATH_CALLBACK

	// set facebookConfig
	FacebookConfig = &oauth2.Config{
		ClientID:     facebookAppID,
		ClientSecret: facebookSecret,
		RedirectURL:  urlCallBack,
		// Scopes:       []string{"public_profile", "email"}, // field email khi nào được facebook phê duyệt mới sử dụng đc
		Scopes:   []string{"public_profile"},
		Endpoint: facebook.Endpoint,
	}
}
