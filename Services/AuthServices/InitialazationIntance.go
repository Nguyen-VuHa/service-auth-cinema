package auth_services

import (
	"context"
	repositories "service-auth/Repositories"
)

// khởi tạo context background để run Redis
var ctx = context.Background()

// Khai báo struct AuthService thông qua dependency injection (repositories.UserRepository)
type AuthService struct {
	userRepository repositories.UserRepository
}

// khởi tạo intance NewAuthService định nghĩa struct AuthService
func NewAuthService(userRepository repositories.UserRepository) *AuthService {
	return &AuthService{userRepository}
}
