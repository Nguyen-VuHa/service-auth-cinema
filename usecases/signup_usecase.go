package usecases

import (
	"auth-service/domains"
	"auth-service/models"
)

type signupUsecase struct {
	userRepository domains.UserRepository
	valRepository  domains.ValidateRepository
}

func NewSignupUsecase(userRepository domains.UserRepository, validate domains.ValidateRepository) domains.SignUpUsecase {
	return &signupUsecase{
		userRepository: userRepository,
		valRepository:  validate,
	}
}

func (su *signupUsecase) ValidateDataRequest(signUpRequest domains.SignUpRequest) error {
	// Kiểm tra Email
	// 1. Email is require
	errIsRequire := su.valRepository.IsRequireString(signUpRequest.Email)

	if errIsRequire != nil {
		return errIsRequire
	}

	// 2. Email hợp lệ
	errIsEmail := su.valRepository.IsEmail(signUpRequest.Email)

	if errIsEmail != nil {
		return errIsEmail
	}

	return nil
}

func (su *signupUsecase) GetUserByEmail(email string) (models.User, error) {
	return su.userRepository.GetByEmail(email)
}

func (su *signupUsecase) Create(user *models.User) error {
	return su.userRepository.Create(user)
}
