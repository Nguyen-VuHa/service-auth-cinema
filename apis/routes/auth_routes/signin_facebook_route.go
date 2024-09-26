package auth_routes

import (
	"auth-service/apis/controllers"
	"auth-service/bootstrap"
	"auth-service/usecases"

	"github.com/gin-gonic/gin"
)

func NewSignInFacebookRouter(group *gin.RouterGroup) {
	sign_in_facebook_usecase := usecases.NewSignInFacebookUsecase(bootstrap.FacebookConfig)
	sign_in_facebook := controllers.SignInFacebookController{
		SignInFacebookUsecase: sign_in_facebook_usecase,
	}

	group.GET("/facebook", sign_in_facebook.SignInWithFacebook)
}
