package auth_routes

import (
	"auth-service/apis/controllers"
	"auth-service/bootstrap"
	"auth-service/repository"
	"auth-service/usecases"

	"github.com/gin-gonic/gin"
)

func NewCallbackFacebookRouter(group *gin.RouterGroup) {
	user_repo := repository.NewUserRepository(bootstrap.DB)
	user_profile_repo := repository.NewUserProfileRepository(bootstrap.DB)
	third_party_repo := repository.NewThirdPartyRepository(bootstrap.DB)

	callback_facebook_usecase := usecases.NewCallbackFacebookUsecase(bootstrap.FacebookConfig, user_repo, user_profile_repo, third_party_repo)

	callback_facebook := controllers.CallbackFacebookController{
		CallbackFacebookUsecase: callback_facebook_usecase,
	}

	group.GET("/callback/facebook", callback_facebook.CallbackSignInWithFacebook)
}
