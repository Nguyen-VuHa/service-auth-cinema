package routes

import (
	"auth-service/apis/routes/auth_routes"
	"auth-service/apis/routes/user_routes"

	"github.com/gin-gonic/gin"
)

func MainRoute(group *gin.RouterGroup) {
	auth_routes.AuthMainRouter(group)
	user_routes.UserMainRouter(group)
}
