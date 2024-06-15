package controllers

import (
	"fmt"
	"net/http"
	constants "service-auth/Constants"
	"service-auth/DTO"
	auth_services "service-auth/Services/AuthServices"
	viewmodels "service-auth/ViewModels"

	"github.com/gin-gonic/gin"
)

type AuthController struct {
	authService *auth_services.AuthService
}

// Hàm khởi tạo UserController
func NewAuthController(authService *auth_services.AuthService) *AuthController {
	return &AuthController{authService}
}

func (service *AuthController) SignUpController(c *gin.Context) {
	var signUpResponse viewmodels.SigUpViewModel

	var bodyRequest DTO.SignUpRequest // khởi tạo bodyRequest

	if err := c.ShouldBindJSON(&bodyRequest); err != nil { // bind data từ request sang bodyRequest
		// write log
		fmt.Println(err.Error())

		// set data ViewModel reponse to user
		signUpResponse.Code = constants.CODE_BAD_REQUEST
		signUpResponse.Status = constants.STATUS_BAD_REQUEST
		signUpResponse.Message = "Invalid JSON format."

		c.JSON(http.StatusBadRequest, signUpResponse)

		return
	}

	fmt.Println(bodyRequest)
	// 2.call service Auth execute function sign up

	service.authService.SignInAccount(bodyRequest)
	// 3.1. return with function errors

	// 3.2. remain return success
	signUpResponse.Code = constants.CODE_SUCCESS
	signUpResponse.Status = constants.STATUS_SUCCESS
	signUpResponse.Message = "Đăng ký thành công"

	c.JSON(http.StatusOK, signUpResponse)
}
