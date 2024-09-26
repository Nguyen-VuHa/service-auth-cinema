package auth_routes

import (
	"auth-service/apis/controllers"
	"auth-service/bootstrap"
	"auth-service/usecases"

	"github.com/gin-gonic/gin"
)

func NewSignInGoogleRouter(group *gin.RouterGroup) {
	sign_in_google_usecase := usecases.NewSignInGoogleUsecase(bootstrap.GoogleConfig)

	sign_in_google := controllers.SignInGoogleController{
		SignInGoogleUsecase: sign_in_google_usecase,
	}

	group.GET("/google", sign_in_google.SignInWithGoogle)
}
