package usecases

import (
	"auth-service/domains"
	"auth-service/models"
	"auth-service/utils"
	"errors"
	"time"
)

type signupUsecase struct {
	userRepository       domains.UserRepository
	userProfileReposiory domains.UserProfileRepository
	valRepository        domains.ValidateRepository
}

func NewSignUpUsecase(userRepository domains.UserRepository, validate domains.ValidateRepository, userProfileRepository domains.UserProfileRepository) domains.SignUpUsecase {
	return &signupUsecase{
		userRepository:       userRepository,
		valRepository:        validate,
		userProfileReposiory: userProfileRepository,
	}
}

func (su *signupUsecase) ValidateDataRequest(signUpRequest domains.SignUpRequest) error {
	validations := []func() error{
		func() error { return su.isEmailValid(signUpRequest.Email) },
		func() error { return su.isPasswordValid(signUpRequest.Password) },
		func() error { return su.isFullnameValid(signUpRequest.FullName) },
		func() error { return su.isBirthDayValid(signUpRequest.BirthDay) },
		func() error { return su.isPhoneNumberValid(signUpRequest.PhoneNumber) },
	}

	for _, validation := range validations {
		if err := validation(); err != nil {
			return err
		}
	}

	return nil
}

func (su *signupUsecase) GetUserByEmail(email string) (models.User, error) {
	return su.userRepository.GetByEmail(email)
}

func (su *signupUsecase) CreateUser(user *models.User) error {
	return su.userRepository.Create(user)
}

func (su *signupUsecase) CreateUserProfile(userProfile *models.UserProfile) error {
	return su.userProfileReposiory.Create(userProfile)
}

// function validation email
func (su *signupUsecase) isEmailValid(email string) error {
	// Kiểm tra Email
	// 1. Email is require
	errIsRequire := su.valRepository.IsRequireString(email)

	if errIsRequire != nil {
		return errIsRequire
	}

	// 2. Email hợp lệ
	errIsEmail := su.valRepository.IsEmail(email)

	if errIsEmail != nil {
		return errIsEmail
	}

	// 3. Độ dài chuỗi cho phép
	errIsLength := su.valRepository.IsMaxLengthString(email, 100) // tối đa 100 ký tự

	if errIsLength != nil {
		return errIsLength
	}

	return nil
}

// function validation password
func (su *signupUsecase) isPasswordValid(password string) error {
	// Kiểm tra password
	// 1. password is require
	errIsRequire := su.valRepository.IsRequireString(password)

	if errIsRequire != nil {
		return errIsRequire
	}

	// 2. kiểm tra range ký tự
	errIsLength := su.valRepository.IsRangeLength(password, 8, 50) // tối thiểu 8 ký tự và tối đa 50 ký tự.

	if errIsLength != nil {
		return errIsLength
	}

	return nil
}

// function validation fullname
func (su *signupUsecase) isFullnameValid(fullName string) error {
	// Kiểm tra fullName
	// 1. fullName is require
	errIsRequire := su.valRepository.IsRequireString(fullName)

	if errIsRequire != nil {
		return errIsRequire
	}

	// 2. kiểm tra range ký tự
	errIsLength := su.valRepository.IsMaxLengthString(fullName, 100) // tối đa 100 ký tự

	if errIsLength != nil {
		return errIsLength
	}

	return nil
}

// function validation fullname
func (su *signupUsecase) isBirthDayValid(birthday string) error {
	// Kiểm tra birthday
	// 1. birthday is require
	errIsRequire := su.valRepository.IsRequireString(birthday)

	if errIsRequire != nil {
		return errIsRequire
	}

	// 2. convert string to date

	timeConvert, errConvert := utils.StringToDate(birthday)

	if errConvert != nil {
		return errors.New("dữ liệu không đúng định dạng")
	}

	// Check if birthday is not in the future
	if timeConvert.After(time.Now()) {
		return errors.New("ngày sinh không thể ở tương lai")
	}

	// Check if year of birthday is >= 1900
	if timeConvert.Year() < 1900 {
		return errors.New("ngày sinh nhật phải lớn năm 1900")
	}

	return nil
}

// function validation fullname
func (su *signupUsecase) isPhoneNumberValid(phoneNumber string) error {
	// Kiểm tra phoneNumber
	// 1. phoneNumber is require
	errIsRequire := su.valRepository.IsRequireString(phoneNumber)

	if errIsRequire != nil {
		return errIsRequire
	}

	// 2. số điện thoại hợp lệ
	errIsPhoneNumber := su.valRepository.IsPhoneNumber(phoneNumber)

	if errIsPhoneNumber != nil {
		return errIsPhoneNumber
	}

	return nil
}
