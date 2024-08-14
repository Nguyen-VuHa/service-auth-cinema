package repository

import (
	"auth-service/domains"
	"errors"
	"fmt"
	"regexp"
)

// Cấu trúc ValidationContext
type validateRepository struct{}

func NewValidation() domains.ValidateRepository {
	return &validateRepository{}
}

// Phương thức Validate cho Validate_IsRequire
func (v validateRepository) IsRequireString(data string) error {
	// kiểm tra dữ liệu rỗng
	if data == "" {
		return errors.New("data is empty")
	}

	return nil
}

// Phương thức Validate cho Validate_IsEmail
func (v validateRepository) IsEmail(data string) error {
	// kiểm tra dữ liệu rỗng
	re := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

	if !re.MatchString(data) {
		return errors.New("invalid data type for email validation")
	}

	return nil
}

// Phương thức Validate cho Validate_IsLength
func (v validateRepository) IsMaxLengthString(data string, maxLength int) error {
	length := len(data)

	if length > maxLength {
		return errors.New("the length of the string must be large " + fmt.Sprint(maxLength) + " characters")
	}

	return nil
}

// Phương thức Validate cho Validate_IsLength
func (v validateRepository) IsRangeLength(data string, minLength, maxLength int) error {
	length := len(data)

	if length < minLength || length > maxLength {
		return errors.New("the length of the string must be between " + fmt.Sprint(minLength) + " and " + fmt.Sprint(maxLength) + " characters")
	}

	return nil
}

// Phương thức Validate cho Validate_IsPhoneNumber
func (v validateRepository) IsPhoneNumber(data string) error {
	// Biểu thức chính quy cho số điện thoại Việt Nam với các độ dài khác nhau
	var re = regexp.MustCompile(`^(03|05|07|08|09)\d{8}$|^(01|02|04|06)\d{7,9}$`)

	if !re.MatchString(data) {
		return errors.New("the phone number invalid")
	}

	return nil
}
