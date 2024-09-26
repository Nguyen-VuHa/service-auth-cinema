package usecases

import (
	"auth-service/constants"
	"auth-service/domains"
	"os"

	"golang.org/x/oauth2"
)

type signInGoogleUsecase struct {
	googleConfig *oauth2.Config
}

func NewSignInGoogleUsecase(googleConfig *oauth2.Config) domains.SignInGoogleUsecase {
	return &signInGoogleUsecase{
		googleConfig,
	}
}

func (sfu *signInGoogleUsecase) AuthGoogleURL() (string, error) {
	url := ""

	sign_state_google := os.Getenv(constants.GOOGLE_SIGN_KEY)
	url = sfu.googleConfig.AuthCodeURL(sign_state_google, oauth2.AccessTypeOffline, oauth2.SetAuthURLParam("prompt", "consent"))

	return url, nil
}
