package controllers

import (
	"auth-service/domains"
	"net/http"

	"github.com/gin-gonic/gin"
)

type SignInGoogleController struct {
	SignInGoogleUsecase domains.SignInGoogleUsecase
}

func (sfc *SignInGoogleController) SignInWithGoogle(c *gin.Context) {
	auth_url, err := sfc.SignInGoogleUsecase.AuthGoogleURL()

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Sign-in with google failed"})

		return
	}

	c.JSON(http.StatusOK, auth_url)
}
