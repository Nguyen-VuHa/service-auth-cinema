package validation_service

import helpers "service-auth/Helpers"

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

// phương thức Auth_IsEmail -> xác thực email trong phần Auth
func (authValid *AuthValidation) Auth_IsEmail(data interface{}) string {

	email, ok := data.(string) // validation dữ liệu đúng định dạng

	if !ok {
		return "Dữ liệu không đúng định dạng"
	}

	// Tạo các đối tượng xác thực cụ thể
	isRequire := Validate_IsRequire{}
	isEmail := &Validate_IsEmail{}
	isLenght := &Validate_IsLength{Min: 0, Max: 80}

	// is require
	authValid.validationContext.SetStrategy(isRequire)
	errIsRequire := authValid.validationContext.Validate(email)

	if errIsRequire != nil {
		objectLog := map[string]interface{}{
			"Error Validation ": errIsRequire.Error(),
		}

		helpers.WriteLogApp("Function Auth_IsEmail() - ValidationService", objectLog, "ERROR")

		return "Trường Email bắt buộc"
	}

	// is email
	authValid.validationContext.SetStrategy(isEmail)
	errIsEmail := authValid.validationContext.Validate(email)

	if errIsEmail != nil {
		objectLog := map[string]interface{}{
			"Error Validation ": errIsEmail.Error(),
		}

		helpers.WriteLogApp("Function Auth_IsEmail() - ValidationService", objectLog, "ERROR")

		return "Email không đúng định dạng."
	}

	// is check length
	authValid.validationContext.SetStrategy(isLenght)
	errIsLength := authValid.validationContext.Validate(email)

	if errIsLength != nil {
		objectLog := map[string]interface{}{
			"Error Validation ": errIsLength.Error(),
		}

		helpers.WriteLogApp("Function Auth_IsEmail() - ValidationService", objectLog, "ERROR")

		return "Độ dài chuỗi không được vượt quá 50 ký tự"
	}

	return ""
}
