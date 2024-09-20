package auth_routes

import (
	"auth-service/apis/controllers"
	"auth-service/bootstrap"
	"auth-service/repository"
	"auth-service/usecases"

	"github.com/gin-gonic/gin"
)

func NewVerifyOTPRouter(group *gin.RouterGroup) {
	userRepo := repository.NewUserRepository(bootstrap.DB)
	redisRepo := repository.NewRedisRepository()

	verify_otp_usecase := usecases.NewVerifyOTPUsecase(userRepo, redisRepo)

	votp := controllers.VerifyOTPController{
		VerifyOTPUsecase: verify_otp_usecase,
	}

	group.GET("/verify-otp", votp.VerifyOTP)
}
