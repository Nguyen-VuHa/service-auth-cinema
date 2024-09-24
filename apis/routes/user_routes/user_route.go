package user_routes

import (
	"auth-service/apis/middlewares"

	"github.com/gin-gonic/gin"
)

func UserMainRouter(group *gin.RouterGroup) {
	userGroup := group.Group("user").Use(middlewares.VerifyTokenProfile())
	{
		NewGetDetailRouter(userGroup)
	}
}
