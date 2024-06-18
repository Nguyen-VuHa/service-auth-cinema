package controllers

import (
	"net/http"
	constants "service-auth/Constants"
	"service-auth/DTO"
	helpers "service-auth/Helpers"
	auth_services "service-auth/Services/AuthServices"
	validation_service "service-auth/Services/ValidationService"
	viewmodels "service-auth/ViewModels"

	"github.com/gin-gonic/gin"
)

// Khai báo struct IntanceAuthController thông qua dependency injection (auth_services.AuthService)
type AuthController struct {
	authService       *auth_services.AuthService
	validationService *validation_service.AuthValidation
}

// khởi tạo intance NewAuthController định nghĩa struct AuthController
func NewAuthController(authService *auth_services.AuthService, validationService *validation_service.AuthValidation) *AuthController {
	return &AuthController{
		authService,
		validationService,
	}
}

func (service *AuthController) SignUpController(c *gin.Context) {
	var signUpResponse viewmodels.SigUpViewModel //

	var bodyRequest DTO.SignUp_Request // khởi tạo bodyRequest

	if err := c.ShouldBindJSON(&bodyRequest); err != nil { // bind data từ request sang bodyRequest
		// write log
		objectLog := map[string]interface{}{
			"Error Bind JSON": err.Error(),
		}
		helpers.WriteLogApp("Function SignUpController() - AuthController", objectLog, "ERROR")

		// set data ViewModel reponse to user
		signUpResponse.Code = constants.CODE_BAD_REQUEST
		signUpResponse.Status = constants.STATUS_BAD_REQUEST
		signUpResponse.Message = "Invalid JSON format."

		c.JSON(http.StatusBadRequest, signUpResponse)
		return
	}

	msgErr := service.validationService.Auth_IsEmail(bodyRequest.Email)
	if msgErr != "" {
		// set data ViewModel reponse to user
		signUpResponse.Code = constants.CODE_BAD_REQUEST
		signUpResponse.Status = constants.STATUS_BAD_REQUEST
		signUpResponse.Message = msgErr

		c.JSON(http.StatusBadRequest, signUpResponse)
		return
	}

	// Xác thực password

	// 2.call service Auth execute function sign up
	_, errResponse, httpS := service.authService.SignUpAccount(bodyRequest)

	// Set Response trả về
	signUpResponse.BaseReponseDTO = errResponse

	c.JSON(httpS.HTTPStatus, signUpResponse)
}
