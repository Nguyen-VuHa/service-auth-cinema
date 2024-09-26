package controllers

import (
	"auth-service/constants"
	"auth-service/domains"
	"auth-service/utils"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

type CallbackFacebookController struct {
	CallbackFacebookUsecase domains.CallbackFacebookUsecase
}

func (cfc *CallbackFacebookController) CallbackSignInWithFacebook(c *gin.Context) {
	url_callback := os.Getenv(constants.URL_CLIENT)

	code := c.Query("code")
	state := c.Query("state")

	signKeyFacebook := os.Getenv(constants.FACEBOOK_SIGN_KEY)

	// state không hợp lệ ~ request đáng nghi
	if signKeyFacebook != state {
		c.JSON(http.StatusForbidden, gin.H{"Message": "Xác thực thất bại"})
		return
	}

	data_callback_facebook, err := cfc.CallbackFacebookUsecase.GetDetailUserWithCodeFacebook(code)

	if err != nil {
		c.Redirect(http.StatusFound, url_callback)
		return
	}

	data_token, err := cfc.CallbackFacebookUsecase.CreateUserLoginWithFacebook(data_callback_facebook)

	if err != nil {
		c.Redirect(http.StatusFound, url_callback)
		return
	}

	var params = make(map[string]interface{})

	params["uid"] = data_token.UserID
	params["access_token"] = data_token.AccessToken
	params["method"] = data_token.Method

	url_redirect_client := url_callback + "/login/success"

	redirect_URL, err := utils.AddParamsToURL(url_redirect_client, params)

	if err != nil {
		c.Redirect(http.StatusFound, url_callback)
		return
	}

	c.Redirect(http.StatusFound, redirect_URL)
}
