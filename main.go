package main

import (
	"auth-service/apis/routes"
	"auth-service/bootstrap"
	"auth-service/constants"
	"net/http"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func init() {
	bootstrap.LoadEnvVariables()   // Khởi tạo các biến trong file .env
	bootstrap.ConnectDatabase()    // Connect với Database
	bootstrap.MigrateDatabase()    // Khởi tạo database trong Models
	bootstrap.InitLogger()         //  Khởi tạo logger cho service
	bootstrap.ConfigFacebookAuth() //  Khởi tạo xác thực facebook cho service
	bootstrap.ConfigGoogleAuth()   //  Khởi tạo xác thực google cho service

	// Kết nối Redis Client
	// Lấy thông tin kết nối của Redis từ biến môi trường
	var redisIPAddr = os.Getenv(constants.REDIS_HOST)
	var userName = os.Getenv(constants.REDIS_USERNAME)
	var redisPassword = os.Getenv(constants.REDIS_PASSWORD)

	// Khởi tạo kết nối Redis cho RedisAuth với database 1
	bootstrap.RedisAuth = bootstrap.InitRedis(redisIPAddr, userName, redisPassword, 01) // #01
	// Khởi tạo kết nối Redis cho RedisUser với database 2
	bootstrap.RedisUser = bootstrap.InitRedis(redisIPAddr, userName, redisPassword, 02) // #02
}

func main() {
	// khởi tạo mặt định Gin framework
	r := gin.Default()

	// config CORS
	config := cors.DefaultConfig()

	config.AllowOrigins = []string{"http://localhost:3000"} // Có thể thay đổi "*" với domain khi triển khai lên môi trường production
	// config CORS allow các method request và các params trên header
	config.AllowCredentials = true
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Content-Type", "Accept", "Authorization"}

	r.Use(cors.New(config)) // set router đi vào CORS kiểm tra trước khi server chuyển hướng router để xử lý

	// API default -> kiểm tra server khi run
	r.GET("/hello-world", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"code":    http.StatusOK,
			"status":  "OK",
			"message": "Hello World API",
		})
	})

	apiGroup := r.Group("api")
	routes.MainRoute(apiGroup)

	// Kiểm tra các routes không hợp lệ không nằm trong khai báo sẽ trả kết quả về ở đây - tương đương với Page Not Found
	r.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{
			"code":    http.StatusNotFound,
			"status":  "ERROR",
			"message": "API NOT FOUND.",
		})
	})

	// run server với default port 5100 hoặc biến PORT trong .env
	r.Run(":5100")
}
