package auth_routes

import (
	"auth-service/apis/controllers"

	"github.com/gin-gonic/gin"
)

func NewSignInRouter(group *gin.RouterGroup) {
	sc := controllers.SignInController{}

	group.POST("/sign-in", sc.SignIn)
}
