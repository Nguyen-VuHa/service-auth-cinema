package validation_service

import (
	"fmt"
	helpers "service-auth/Helpers"
	"time"
)

// AuthValidation struct
type AuthValidation struct {
	validationContext *ValidationContext
}

// NewAuthValidation constructor
func NewAuthValidation() *AuthValidation {
	return &AuthValidation{
		validationContext: &ValidationContext{},
	}
}

// phương thức Auth_IsEmail -> kiểm tra email trong phần Auth
func (authValid *AuthValidation) Auth_IsEmail(data interface{}) string {
	email, ok := data.(string) // validation dữ liệu đúng định dạng

	if !ok {
		objectLog := map[string]interface{}{
			"Error Validation Email": ok,
		}

		helpers.WriteLogApp("Function Auth_IsEmail() - ValidationService", objectLog, "ERROR")

		return "Dữ liệu không đúng định dạng"
	}

	// Tạo các đối tượng kiểm tra cụ thể
	isRequire := Validate_IsRequire{}
	isEmail := &Validate_IsEmail{}
	isLenght := &Validate_IsLength{Min: 0, Max: 100}

	// is require
	authValid.validationContext.SetStrategy(isRequire)
	errIsRequire := authValid.validationContext.Validate(email)

	if errIsRequire != nil {
		objectLog := map[string]interface{}{
			"Error Validation Email": errIsRequire.Error(),
		}

		helpers.WriteLogApp("Function Auth_IsEmail() - ValidationService", objectLog, "ERROR")

		return "Trường `Email` bắt buộc"
	}

	// is email
	authValid.validationContext.SetStrategy(isEmail)
	errIsEmail := authValid.validationContext.Validate(email)

	if errIsEmail != nil {
		objectLog := map[string]interface{}{
			"Error Validation Email": errIsEmail.Error(),
		}

		helpers.WriteLogApp("Function Auth_IsEmail() - ValidationService", objectLog, "ERROR")

		return "Email không đúng định dạng."
	}

	// is check length
	authValid.validationContext.SetStrategy(isLenght)
	errIsLength := authValid.validationContext.Validate(email)

	if errIsLength != nil {
		objectLog := map[string]interface{}{
			"Error Validation Email": errIsLength.Error(),
		}

		helpers.WriteLogApp("Function Auth_IsEmail() - ValidationService", objectLog, "ERROR")

		return `Độ dài chuỗi không được vượt quá ` + fmt.Sprint(isLenght.Max) + ` ký tự`
	}

	return ""
}

// phương thức Auth_IsPassword -> kiểm tra password trong phần Auth
func (authValid *AuthValidation) Auth_IsPassword(data interface{}) string {
	password, ok := data.(string) // validation dữ liệu đúng định dạng

	if !ok {
		objectLog := map[string]interface{}{
			"Error Validation Password": ok,
		}

		helpers.WriteLogApp("Function Auth_IsEmail() - ValidationService", objectLog, "ERROR")

		return "Dữ liệu không đúng định dạng"
	}

	// Tạo các đối tượng kiểm tra cụ thể
	isRequire := Validate_IsRequire{}
	isLenght := &Validate_IsLength{Min: 8, Max: 50}

	// is require
	authValid.validationContext.SetStrategy(isRequire)
	errIsRequire := authValid.validationContext.Validate(password)

	if errIsRequire != nil {
		objectLog := map[string]interface{}{
			"Error Validation Password": errIsRequire.Error(),
		}

		helpers.WriteLogApp("Function Auth_IsEmail() - ValidationService", objectLog, "ERROR")

		return "Trường `Mật khẩu` bắt buộc"
	}

	// is check length
	authValid.validationContext.SetStrategy(isLenght)
	errIsLength := authValid.validationContext.Validate(password)

	if errIsLength != nil {
		objectLog := map[string]interface{}{
			"Error Validation Password": errIsLength.Error(),
		}

		helpers.WriteLogApp("Function Auth_IsEmail() - ValidationService", objectLog, "ERROR")

		return `Độ dài chuỗi tối thiểu ` + fmt.Sprint(isLenght.Min) + ` ký tự và không được vượt quá ` + fmt.Sprint(isLenght.Max) + ` ký tự`
	}

	return ""
}

// phương thức Auth_IsFullName -> kiểm tra tên người dùng trong phần Auth
func (authValid *AuthValidation) Auth_IsFullName(data interface{}) string {
	fullName, ok := data.(string) // validation dữ liệu đúng định dạng

	if !ok {
		objectLog := map[string]interface{}{
			"Error Validation FullName": ok,
		}

		helpers.WriteLogApp("Function Auth_IsEmail() - ValidationService", objectLog, "ERROR")

		return "Dữ liệu không đúng định dạng"
	}

	// Tạo các đối tượng kiểm tra cụ thể
	isRequire := Validate_IsRequire{}
	isLenght := &Validate_IsLength{Min: 0, Max: 100}

	// is require
	authValid.validationContext.SetStrategy(isRequire)
	errIsRequire := authValid.validationContext.Validate(fullName)

	if errIsRequire != nil {
		objectLog := map[string]interface{}{
			"Error Validation FullName": errIsRequire.Error(),
		}

		helpers.WriteLogApp("Function Auth_IsEmail() - ValidationService", objectLog, "ERROR")

		return "Trường `Họ & Tên` bắt buộc"
	}

	// is check length
	authValid.validationContext.SetStrategy(isLenght)
	errIsLength := authValid.validationContext.Validate(fullName)

	if errIsLength != nil {
		objectLog := map[string]interface{}{
			"Error Validation FullName": errIsLength.Error(),
		}

		helpers.WriteLogApp("Function Auth_IsEmail() - ValidationService", objectLog, "ERROR")

		return `Độ dài chuỗi không được vượt quá ` + fmt.Sprint(isLenght.Max) + ` ký tự`
	}

	return ""
}

// phương thức Auth_IsBirthDay -> kiểm tra ngày sinh hợp lệ dùng trong phần Auth
func (authValid *AuthValidation) Auth_IsBirthDay(data interface{}) string {
	birthDay, ok := data.(time.Time)

	if !ok {
		objectLog := map[string]interface{}{
			"Error Validation BirthDay": ok,
		}

		helpers.WriteLogApp("Function Auth_IsEmail() - ValidationService", objectLog, "ERROR")
		return "Dữ liệu không đúng định dạng"
	}

	// Check if birthday is not in the future
	if birthDay.After(time.Now()) {
		objectLog := map[string]interface{}{
			"Error Validation BirthDay": "Birthday is invalid",
		}

		helpers.WriteLogApp("Function Auth_IsEmail() - ValidationService", objectLog, "ERROR")

		return "Ngày sinh không thể ở tương lai"
	}

	// Check if year of birthday is >= 1900
	if birthDay.Year() < 1900 {
		objectLog := map[string]interface{}{
			"Error Validation BirthDay": "Birthday is invalid",
		}

		helpers.WriteLogApp("Function Auth_IsEmail() - ValidationService", objectLog, "ERROR")

		return "Ngày sinh nhật phải lớn năm 1900"
	}

	return ""
}

// phương thức Auth_IsPhoneNumber -> kiểm tra định dạng số điện thoại trong phần Auth
func (authValid *AuthValidation) Auth_IsPhoneNumber(data interface{}) string {
	phoneNumber, ok := data.(string)

	if !ok {
		objectLog := map[string]interface{}{
			"Error Validation PhoneNumber": ok,
		}

		helpers.WriteLogApp("Function Auth_IsEmail() - ValidationService", objectLog, "ERROR")

		return "Dữ liệu không đúng định dạng"
	}

	// Tạo các đối tượng kiểm tra cụ thể
	isRequire := Validate_IsRequire{}
	isPhoneNumber := &Validate_IsPhoneNumber{}

	// is require
	authValid.validationContext.SetStrategy(isRequire)
	errIsRequire := authValid.validationContext.Validate(phoneNumber)

	if errIsRequire != nil {
		objectLog := map[string]interface{}{
			"Error Validation PhoneNumber": errIsRequire.Error(),
		}

		helpers.WriteLogApp("Function Auth_IsEmail() - ValidationService", objectLog, "ERROR")

		return "Trường `Số điện thoại` bắt buộc"
	}

	// is email
	authValid.validationContext.SetStrategy(isPhoneNumber)
	errIsEmail := authValid.validationContext.Validate(phoneNumber)

	if errIsEmail != nil {
		objectLog := map[string]interface{}{
			"Error Validation PhoneNumber": errIsEmail.Error(),
		}

		helpers.WriteLogApp("Function Auth_IsEmail() - ValidationService", objectLog, "ERROR")

		return "Số điện thoại không đúng định dạng."
	}

	return ""
}
