package validation_service

import (
	"errors"
	"fmt"
	"regexp"
)

// Interface Validator định nghĩa phương thức Validate
type Validator interface {
	Validate(data string) error
}

// Cấu trúc ValidationContext
type ValidationContext struct {
	strategy Validator
}

// Phương thức SetStrategy để thiết lập chiến lược xác thực
func (v *ValidationContext) SetStrategy(strategy Validator) {
	v.strategy = strategy
}

// Phương thức Validate để thực hiện xác thực
func (v *ValidationContext) Validate(data string) error {
	return v.strategy.Validate(data)
}

// Validate_IsRequire  struct
type Validate_IsRequire struct{}

// Phương thức Validate cho Validate_IsRequire
func (v Validate_IsRequire) Validate(data string) error {
	// kiểm tra dữ liệu rỗng
	if data == "" {
		return errors.New("data is empty")
	}

	return nil
}

// Validate_IsRequire  struct
type Validate_IsEmail struct{}

// Phương thức Validate cho Validate_IsEmail
func (v Validate_IsEmail) Validate(data string) error {
	// kiểm tra dữ liệu rỗng
	re := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

	if !re.MatchString(data) {
		return errors.New("invalid data type for email validation")
	}

	return nil
}

// Cấu trúc Validate_IsLength
type Validate_IsLength struct {
	Min int
	Max int
}

// Phương thức Validate cho Validate_IsLength
func (v Validate_IsLength) Validate(data string) error {
	length := len(data)

	if length < v.Min || length > v.Max {
		return errors.New("the length of the string must be between " + fmt.Sprint(v.Min) + " and " + fmt.Sprint(v.Max) + " characters")
	}

	return nil
}

// Cấu trúc Validate_IsPhoneNumber (Only used VietNam)
type Validate_IsPhoneNumber struct{}

// Phương thức Validate cho Validate_IsPhoneNumber
func (v Validate_IsPhoneNumber) Validate(data string) error {
	// Biểu thức chính quy cho số điện thoại Việt Nam với các độ dài khác nhau
	var re = regexp.MustCompile(`^(03|05|07|08|09)\d{8}$|^(01|02|04|06)\d{7,9}$`)

	if !re.MatchString(data) {
		return errors.New("the phone number invalid")
	}

	return nil
}
