package auth_routes

import (
	"auth-service/apis/controllers"
	"auth-service/bootstrap"
	"auth-service/repository"
	"auth-service/usecases"

	"github.com/gin-gonic/gin"
)

func NewRefreshTokenRouter(group *gin.RouterGroup) {
	redis_repo := repository.NewRedisRepository()
	third_party_repo := repository.NewThirdPartyRepository(bootstrap.DB)
	refresh_token_usecase := usecases.NewRefreshTokenUsecase(redis_repo, third_party_repo, bootstrap.GoogleConfig)

	rfc := controllers.RefreshTokenController{
		RefreshTokenUsecase: refresh_token_usecase,
	}

	group.GET("/refresh", rfc.RefreshToken)
}
