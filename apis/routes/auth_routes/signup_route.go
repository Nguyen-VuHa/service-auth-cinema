package auth_routes

import (
	"auth-service/apis/controllers"
	"auth-service/bootstrap"
	"auth-service/repository"
	"auth-service/usecases"

	"github.com/gin-gonic/gin"
)

func NewSignUpRouter(group *gin.RouterGroup) {
	useRepo := repository.NewUserRepository(bootstrap.DB)
	userProfileRepo := repository.NewUserProfileRepository(bootstrap.DB)
	validateRepo := repository.NewValidation()

	signUpUseCase := usecases.NewSignUpUsecase(useRepo, validateRepo, userProfileRepo)
	sc := controllers.SignupController{
		SignUpUseCase: signUpUseCase,
	}

	group.POST("/sign-up", sc.SignUp)
}
