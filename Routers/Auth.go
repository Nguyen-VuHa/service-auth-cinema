package routers

import (
	controllers "service-auth/Controllers"
	data_layers "service-auth/DataLayers"
	initializers "service-auth/Initializers"
	repositories "service-auth/Repositories"
	auth_services "service-auth/Services/AuthServices"

	"github.com/gin-gonic/gin"
)

func AuthRoutes(routes *gin.RouterGroup) {

	var userDataLayer = data_layers.NewUserDataLayer(initializers.DB)
	var userRepository = repositories.NewIntanceUserDataLayer(userDataLayer)
	var authService = auth_services.NewAuthService(userRepository)
	var authController = controllers.NewAuthController(authService)

	authGroup := routes.Group("/auth")
	{
		authGroup.POST("/sign-in")
		authGroup.POST("/sign-up", authController.SignUpController)
		authGroup.POST("/facebook")
		authGroup.POST("/google")

		authCallBackGroup := authGroup.Group("/callback")
		{
			authCallBackGroup.GET("/facebook")
			authCallBackGroup.GET("/google")
		}
	}
}
