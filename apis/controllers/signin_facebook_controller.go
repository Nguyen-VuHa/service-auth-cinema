package controllers

import (
	"auth-service/domains"
	"net/http"

	"github.com/gin-gonic/gin"
)

type SignInFacebookController struct {
	SignInFacebookUsecase domains.SignInFacebookUsecase
}

func (sfc *SignInFacebookController) SignInWithFacebook(c *gin.Context) {

	auth_url, err := sfc.SignInFacebookUsecase.AuthFacebookURL()

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Sign-in with facebook failed"})
		return
	}

	c.JSON(http.StatusOK, auth_url)
}
