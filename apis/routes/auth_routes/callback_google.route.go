package auth_routes

import (
	"auth-service/apis/controllers"
	"auth-service/bootstrap"
	"auth-service/repository"
	"auth-service/usecases"

	"github.com/gin-gonic/gin"
)

func NewCallbackGoogleRouter(group *gin.RouterGroup) {
	user_repo := repository.NewUserRepository(bootstrap.DB)
	user_profile_repo := repository.NewUserProfileRepository(bootstrap.DB)
	third_party_repo := repository.NewThirdPartyRepository(bootstrap.DB)
	redis_repo := repository.NewRedisRepository()

	callback_google_usecase := usecases.NewCallbackGoogleUsecase(
		bootstrap.GoogleConfig, user_repo,
		user_profile_repo, third_party_repo, redis_repo,
	)

	callback_google := controllers.CallbackGoogleController{
		CallbackGoogleUsecase: callback_google_usecase,
	}

	group.GET("/callback/google", callback_google.CallbackSignInWithGoogle)
}
