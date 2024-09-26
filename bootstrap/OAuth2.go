package bootstrap

import (
	"auth-service/constants"
	"os"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/facebook"
	"golang.org/x/oauth2/google"
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

func ConfigGoogleAuth() {
	// get thông tin ứng dụng xác thực facebook trong env
	google_client_id := os.Getenv(constants.GOOGLE_CLIENT_ID)
	google_secret_id := os.Getenv(constants.GOOGLE_SECRET_ID)
	urlCallBack := os.Getenv(constants.URL_CALLBACK)
	urlCallBack += constants.GOOGLE_PATH_CALLBACK

	// set facebookConfig
	GoogleConfig = &oauth2.Config{
		ClientID:     google_client_id, // Thay bằng Client ID của bạn
		ClientSecret: google_secret_id, // Thay bằng Client Secret của bạn
		RedirectURL:  urlCallBack,      // URL nhận callback từ Google
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.profile", "https://www.googleapis.com/auth/userinfo.email"},
		Endpoint:     google.Endpoint,
	}
}
