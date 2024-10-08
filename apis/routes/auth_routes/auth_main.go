package auth_routes

import "github.com/gin-gonic/gin"

func AuthMainRouter(group *gin.RouterGroup) {
	authGroup := group.Group("auth")
	{
		NewSignUpRouter(authGroup)
		NewSignInRouter(authGroup)
		NewVerifyOTPRouter(authGroup)
	}
}
