package usecases

import (
	"auth-service/constants"
	"auth-service/domains"
	"os"

	"golang.org/x/oauth2"
)

type signInFacebookUsecase struct {
	facebookConfig *oauth2.Config
}

func NewSignInFacebookUsecase(facebookConfig *oauth2.Config) domains.SignInFacebookUsecase {
	return &signInFacebookUsecase{
		facebookConfig,
	}
}

func (sfu *signInFacebookUsecase) AuthFacebookURL() (string, error) {
	url := ""

	sign_state_facebook := os.Getenv(constants.FACEBOOK_SIGN_KEY)
	url = sfu.facebookConfig.AuthCodeURL(sign_state_facebook)

	return url, nil
}
