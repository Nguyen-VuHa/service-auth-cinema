package main

import (
	"net/http"
	initializers "service-auth/Initializers"
	routers "service-auth/Routers"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func init() {
	initializers.LoadEnvVariables()
	initializers.ConnectDatabase()
	initializers.MigrateDatabase()
}

func main() {
	r := gin.Default()

	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"*"} // You can change "*" with your specific domain
	config.AllowCredentials = true
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Content-Type", "Accept", "Authorization"}

	r.Use(cors.New(config))

	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"code":    http.StatusOK,
			"status":  "OK",
			"message": "Hello World API",
		})
	})

	apiGroup := r.Group("/api")
	routers.MainRoutes(apiGroup)

	r.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{
			"code":    http.StatusNotFound,
			"status":  "ERROR",
			"message": "NETWORK ERROR.",
		})
	})

	r.Run()
}
