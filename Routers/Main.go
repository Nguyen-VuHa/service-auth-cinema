package routers

import (
	"github.com/gin-gonic/gin"
)

func MainRoutes(routes *gin.RouterGroup) {
	AuthRoutes(routes)
}
