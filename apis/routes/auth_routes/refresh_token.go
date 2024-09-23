package auth_routes

import (
	"auth-service/apis/controllers"
	"auth-service/repository"
	"auth-service/usecases"

	"github.com/gin-gonic/gin"
)

func NewRefreshTokenRouter(group *gin.RouterGroup) {
	redis_repo := repository.NewRedisRepository()
	refresh_token_usecase := usecases.NewRefreshTokenUsecase(redis_repo)

	rfc := controllers.RefreshTokenController{
		RefreshTokenUsecase: refresh_token_usecase,
	}

	group.GET("/refresh", rfc.RefreshToken)
}
