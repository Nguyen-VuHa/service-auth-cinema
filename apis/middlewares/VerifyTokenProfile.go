package middlewares

import (
	"auth-service/constants"
	"auth-service/internal/jwt_util"
	"auth-service/repository"
	"auth-service/utils"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

func VerifyTokenProfile() gin.HandlerFunc {
	return func(c *gin.Context) {
		user_id_params := c.Query("_user_id")

		tokenString := c.GetHeader("Authorization")
		tokenString = strings.TrimPrefix(tokenString, "Bearer ")

		if tokenString == "" {
			c.JSON(http.StatusForbidden, gin.H{"error": "Permission Denied"})
			c.Abort()
			return
		}

		tokenArr := strings.Split(tokenString, ".")

		if len(tokenArr) > 1 {
			methodID := tokenArr[0]
			tokenPlatform := tokenArr[1]

			if methodID == fmt.Sprint(constants.LOGIN_FACEBOOK_ID) {
				facebook_app_id := os.Getenv(constants.FACEBOOK_APP_ID)
				facebook_app_secret := os.Getenv(constants.FACEBOOK_SECRET)

				is_valid, err := VerifyTokenFacebook(tokenPlatform, facebook_app_id+"|"+facebook_app_secret)

				if err != nil {
					c.JSON(http.StatusForbidden, gin.H{"error": "Permission Denied"})
					c.Abort()
					return
				}

				if !is_valid {
					c.JSON(http.StatusForbidden, gin.H{"error": "Permission Denied"})
					c.Abort()
					return
				}
			}

			c.Next()
			return
		}

		// kiểm tra access token hợp lệ không
		device := utils.GetDevice(c)
		sign_access := os.Getenv(constants.JWT_ACCESS_SECRET)
		sign_hash_access := sign_access + device
		user_id, err := jwt_util.VerifyJWTToken(tokenString, sign_hash_access)

		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			c.Abort()

			return
		}

		if user_id_params != user_id {
			c.JSON(http.StatusForbidden, gin.H{"error": "Permission Denied"})
			c.Abort()
			return
		}

		// tạo redis key bằng keyDEVICE_SECRET_KEY + UserID + Device
		sign_redis := os.Getenv(constants.DEVICE_SECRET_KEY) + user_id_params + device

		// Kiểm tra xem có đúng user đó không
		fields := []string{"AccessToken"}

		redis_repo := repository.NewRedisRepository()

		data_redis, err := redis_repo.RedisAuthHMGetFields(sign_redis, fields)

		if err != nil {
			c.JSON(http.StatusForbidden, gin.H{"error": "Permission Denied"})
			c.Abort()
			return
		}

		if data_redis[fields[0]] == nil {
			c.JSON(http.StatusForbidden, gin.H{"error": "Permission Denied"})
			c.Abort()
			return
		}

		if data_redis[fields[0]].(string) != tokenString {
			c.JSON(http.StatusForbidden, gin.H{"error": "Permission Denied"})
			c.Abort()
			return
		}

		c.Next()
	}
}
