package auth_services

import (
	"fmt"
	"service-auth/DTO"
	repositories "service-auth/Repositories"
)

// AuthService định nghĩa các interface ở bên trên
type AuthService struct {
	userRepository repositories.UserRepository
}

func NewAuthService(userRepository repositories.UserRepository) *AuthService {
	return &AuthService{userRepository}
}

func (repo *AuthService) SignInAccount(dataRequest DTO.SignUpRequest) interface{} {
	// Logic đăng ký
	// 1. Kiểm tra tồn tại của email
	user, err := repo.userRepository.GetUserByEmail(dataRequest.Email)

	fmt.Println(user)
	fmt.Println(err)
	if err == nil { // email tồn tại -> thông báo mã lỗi và trả về kết quả failed

	}

	// 2. hash password với passhash: email + uuid

	// 3. insert thông tin vào database tương ứng

	// 4. trả về kết quả
	return user
}
