package auth_routes

import (
	"auth-service/apis/controllers"
	"auth-service/bootstrap"
	"auth-service/repository"
	"auth-service/usecases"

	"github.com/gin-gonic/gin"
)

func NewSignInRouter(group *gin.RouterGroup) {
	userRepo := repository.NewUserRepository(bootstrap.DB)
	validateRepo := repository.NewValidation()
	redisRepo := repository.NewRedisRepository()
	serviceMailRepo := repository.NewServiceMailRepository()

	signInUsercase := usecases.NewSignInUsecase(userRepo, validateRepo, redisRepo, serviceMailRepo)
	sc := controllers.SignInController{
		SignInUsecase: signInUsercase,
	}

	group.POST("/sign-in", sc.SignIn)
}
