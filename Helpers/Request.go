package helpers

import (
	"net/http"
	"strings"
)

// lấy thông tin địa chỉ IP từ client request
func GetIPClient(r *http.Request) string {
	// Kiểm tra các header X-Forwarded-For và X-Real-IP nếu server nằm sau proxy
	forwarded := r.Header.Get("X-Forwarded-For")

	if forwarded != "" {
		// Trường hợp có nhiều IP trong X-Forwarded-For, lấy IP đầu tiên
		ip := strings.Split(forwarded, ",")[0]
		return strings.TrimSpace(ip)
	}

	realIP := r.Header.Get("X-Real-IP")
	if realIP != "" {
		return realIP
	}

	// Lấy IP từ RemoteAddr nếu không có các header trên
	ip := r.RemoteAddr
	// Trường hợp RemoteAddr có chứa port, chỉ lấy phần IP
	if strings.Contains(ip, ":") {
		ip = strings.Split(ip, ":")[0]
	}
	return ip
}

// lấy thông tin device gửi lên từ header request
func GetUserAgent(r *http.Request) string {
	return r.Header.Get("User-Agent")
}
