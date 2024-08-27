package auth_routes

import (
	"auth-service/apis/controllers"
	"auth-service/bootstrap"
	"auth-service/repository"
	"auth-service/usecases"

	"github.com/gin-gonic/gin"
)

func NewSignInRouter(group *gin.RouterGroup) {
	useRepo := repository.NewUserRepository(bootstrap.DB)
	validateRepo := repository.NewValidation()

	signInUsercase := usecases.NewSignInUsecase(useRepo, validateRepo)
	sc := controllers.SignInController{
		SignInUsecase: signInUsercase,
	}

	group.POST("/sign-in", sc.SignIn)
}
