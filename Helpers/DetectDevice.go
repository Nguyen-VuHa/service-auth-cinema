package helpers

import (
	"strings"

	"github.com/gin-gonic/gin"
)

func GetClientIP(c *gin.Context) string {
	// Lấy địa chỉ IP từ header X-Forwarded-For trước tiên
	ip := c.GetHeader("X-Forwarded-For")
	if ip != "" {
		// Nếu có nhiều địa chỉ IP được liệt kê, chỉ lấy cái đầu tiên
		ip = strings.TrimSpace(strings.Split(ip, ",")[0])
		return ip
	}

	// Nếu không có header X-Forwarded-For, sử dụng RemoteIP của Gin
	ip = c.RemoteIP()

	return ip
}

func GetDevice(c *gin.Context) string {
	device := c.Request.UserAgent() // lấy thông tin device trong request header.

	return device
}
