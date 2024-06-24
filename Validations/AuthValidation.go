package validations

import (
	"service-auth/DTO"
	validation_service "service-auth/Services/ValidationService"
)

func Valid_Auth_SignUp(bodyRequest DTO.SignUp_Request, validationService *validation_service.AuthValidation) string {
	var messageError = ""

	messageError = validationService.Auth_IsEmail(bodyRequest.Email) // Kiểm tra trường Email

	if messageError != "" { // nếu lỗi trả về message lỗi
		return messageError
	}

	messageError = validationService.Auth_IsPassword(bodyRequest.Password) // Kiểm tra trường Password

	if messageError != "" { // nếu lỗi trả về message lỗi
		return messageError
	}

	messageError = validationService.Auth_IsFullName(bodyRequest.FullName) // Kiểm tra trường FullName

	if messageError != "" { // nếu lỗi trả về message lỗi
		return messageError
	}

	messageError = validationService.Auth_IsBirthDay(bodyRequest.BirthDay) // Kiểm tra trường BirthDay

	if messageError != "" { // nếu lỗi trả về message lỗi
		return messageError
	}

	messageError = validationService.Auth_IsPhoneNumber(bodyRequest.PhoneNumber) // Kiểm tra trường PhoneNumber

	if messageError != "" { // nếu lỗi trả về message lỗi
		return messageError
	}

	return ""
}

func Valid_Auth_SignIn(bodyRequest DTO.SignIn_Request, validationService *validation_service.AuthValidation) string {
	var messageError = ""

	messageError = validationService.Auth_IsEmail(bodyRequest.Email) // Kiểm tra trường Email

	if messageError != "" { // nếu lỗi trả về message lỗi
		return messageError
	}

	messageError = validationService.Auth_IsPassword(bodyRequest.Password) // Kiểm tra trường Password

	if messageError != "" { // nếu lỗi trả về message lỗi
		return messageError
	}

	return ""
}
