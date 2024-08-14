package utils

import "golang.org/x/crypto/bcrypt"

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
