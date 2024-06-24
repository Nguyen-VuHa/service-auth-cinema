package helpers

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"

	"golang.org/x/crypto/bcrypt"
)

// function này thường sử dụng tạo mới user hoặc thay đổi mật khẩu cần mã hoá password
func HashPasswordWithBcrypt(password string) (string, error) {
	// này là Cost Factor - thường được gọi là yếu tố chi phí xác định độ phức tạp của quá trình hash
	// Giá trị này càng cao thì việc tính toán hash càng tốn nhiều thời gian và tài nguyên.
	cost := 10 // Cost factor mặc định thường được khuyến nghị là 10. Đây là một giá trị cân bằng tốt giữa bảo mật và hiệu suất.

	// hash password thông qua function GenerateFromPassword của bcrypt
	hash, err := bcrypt.GenerateFromPassword([]byte(password), cost)

	if err != nil {
		return string(hash), err
	}

	return string(hash), nil
}

// passwordHash: là password được hash từ funtion HashPasswordWithBcrypt - (lưu trữ trong Database)
// passwordInput: là password lấy từ request gửi lên để so sánh - (password user nhập)
func ComparePasswordByBcrypt(passwordHash string, passwordInput string) error {
	// Kiểm tra tính hợp lệ giữa mật khẩu người dùng và mật khẩu lưu trữ
	// thông qua function CompareHashAndPassword của thư viện bcrypt
	return bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(passwordInput))
}

// Hàm tạo HMAC hash
func GenerateHMAC(secret, message string) string {
	key := []byte(secret) // Chuyển đổi khóa bí mật từ chuỗi sang byte slice
	// Tạo một đối tượng HMAC mới sử dụng thuật toán SHA-256 và khóa bí mật
	h := hmac.New(sha256.New, key)
	// Ghi thông điệp vào đối tượng HMAC
	h.Write([]byte(message))
	// Tính toán giá trị HMAC và chuyển đổi thành chuỗi hexadecimal
	return hex.EncodeToString(h.Sum(nil))
}

// Hàm kiểm tra HMAC hash
func ValidateHMAC(secret, message, hash string) bool {
	// Tạo lại HMAC hash từ thông điệp và khóa bí mật
	generatedHash := GenerateHMAC(secret, message)
	// So sánh hash được tạo ra với hash được cung cấp
	// hmac.Equal đảm bảo so sánh an toàn, tránh các tấn công timing attacks
	return hmac.Equal([]byte(generatedHash), []byte(hash))
}
