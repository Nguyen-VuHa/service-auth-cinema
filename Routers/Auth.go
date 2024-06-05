package routers

import "github.com/gin-gonic/gin"

func AuthRoutes(routes *gin.RouterGroup) {
	authGroup := routes.Group("/auth")
	{
		authGroup.POST("/sign-in")
		authGroup.POST("/sign-up")
		authGroup.POST("/facebook")
		authGroup.POST("/google")

		authCallBackGroup := authGroup.Group("/callback")
		{
			authCallBackGroup.GET("/facebook")
			authCallBackGroup.GET("/google")
		}
	}
}
