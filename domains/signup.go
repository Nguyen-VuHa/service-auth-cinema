package domains

import "auth-service/models"

type SignUpRequest struct {
	Email       string `json:"email"`
	Password    string `json:"password"`
	FullName    string `json:"full_name"`
	BirthDay    string `json:"birth_day"`
	PhoneNumber string `json:"phone_number"`
	Address     string `json:"address"`
}

type SignUpResponse struct {
	ResponseBasic
	Data DataSignUpResponse `json:"data"`
}

type DataSignUpResponse struct {
}

type SignUpUsecase interface {
	ValidateDataRequest(signUpRequest SignUpRequest) error
	GetUserByEmail(email string) (models.User, error)
	CreateUser(user *models.User) error
	CreateUserProfile(userProfile *models.UserProfile) error
}
