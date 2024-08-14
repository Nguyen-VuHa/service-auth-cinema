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
	validateRepo := repository.NewValidation()

	SignUpUseCase := usecases.NewSignupUsecase(useRepo, validateRepo)
	sc := controllers.SignupController{
		SignUpUseCase: SignUpUseCase,
	}

	group.POST("/sign-up", sc.SignUp)
}
