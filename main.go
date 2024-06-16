package main

import (
	"net/http"
	initializers "service-auth/Initializers"
	routers "service-auth/Routers"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func init() {
	initializers.LoadEnvVariables() // Khởi tạo các biến trong file .env
	initializers.ConnectDatabase()  // Connect với Database
	initializers.MigrateDatabase()  // Khởi tạo database trong Models
	initializers.InitLogger()       //  Khởi tạo logger cho service
}

func main() {
	// Lấy logger đã khởi tạo từ init.go
	Logger := initializers.GetLogger()
	defer Logger.Sync() // Đảm bảo rằng các log sẽ được ghi lại

	// khởi tạo mặt định Gin framework
	r := gin.Default()

	// config CORS
	config := cors.DefaultConfig()

	config.AllowOrigins = []string{"*"} // Có thể thay đổi "*" với domain khi triển khai lên môi trường production
	// config CORS allow các method request và các params trên header
	config.AllowCredentials = true
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Content-Type", "Accept", "Authorization"}

	r.Use(cors.New(config)) // set router đi vào CORS kiểm tra trước khi server chuyển hướng router để xử lý

	// API default -> kiểm tra server khi run
	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"code":    http.StatusOK,
			"status":  "OK",
			"message": "Hello World API",
		})
	})

	// group API với đường dẫn /api
	apiGroup := r.Group("/api")
	routers.MainRoutes(apiGroup)

	// Kiểm tra các routes không hợp lệ không nằm trong khai báo sẽ trả kết quả về ở đây - tương đương với Page Not Found
	r.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{
			"code":    http.StatusNotFound,
			"status":  "ERROR",
			"message": "NETWORK ERROR.",
		})
	})

	// run server với default port 8080 hoặc biến PORT trong .env
	r.Run()
}
