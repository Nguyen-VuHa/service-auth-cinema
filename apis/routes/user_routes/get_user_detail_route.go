package user_routes

import (
	"auth-service/apis/controllers"
	"auth-service/bootstrap"
	"auth-service/repository"
	"auth-service/usecases"

	"github.com/gin-gonic/gin"
)

func NewGetDetailRouter(group gin.IRoutes) {
	redis_repo := repository.NewRedisRepository()
	user_repo := repository.NewUserRepository(bootstrap.DB)
	get_detail_user_usecase := usecases.NewGetDetailUserUsecase(redis_repo, user_repo)

	gduc := controllers.GetDetailUserController{
		GetDetailUserUsecase: get_detail_user_usecase,
	}

	group.GET("/profile", gduc.GetDetailByID)
}
