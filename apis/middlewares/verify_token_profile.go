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
		methodID := ""

		if len(tokenArr) > 1 {
			methodID = tokenArr[0]
			tokenPlatform := strings.Join(tokenArr[1:], ".")

			if methodID == fmt.Sprint(constants.LOGIN_GOOGLE_ID) {
				tokenString = tokenPlatform
			}

			// kiểm tra chỉ method thuộc các phương thức đặt biệt mới xác thực.
			arr_method := []string{fmt.Sprint(constants.LOGIN_FACEBOOK_ID)}

			mothod_is_valid := utils.ItemIsArrayString(methodID, arr_method)

			if mothod_is_valid {
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
		}

		// kiểm tra access token hợp lệ không
		device := utils.GetDevice(c)
		sign_hash_access := os.Getenv(constants.JWT_ACCESS_SECRET)

		if methodID != fmt.Sprint(constants.LOGIN_GOOGLE_ID) {
			sign_hash_access = sign_hash_access + device
		}

		user_id, err := jwt_util.VerifyJWTToken(tokenString, sign_hash_access)

		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			c.Abort()

			return
		}

		if user_id_params != "" {
			if user_id_params != user_id {
				c.JSON(http.StatusForbidden, gin.H{"error": "Permission Denied"})
				c.Abort()
				return
			}
		}

		// tạo redis key bằng keyDEVICE_SECRET_KEY + UserID + Device
		sign_redis := os.Getenv(constants.DEVICE_SECRET_KEY) + user_id_params + device

		if methodID == fmt.Sprint(constants.LOGIN_GOOGLE_ID) {
			sign_redis = user_id
		}

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
