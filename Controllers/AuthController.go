package controllers

import (
	"net/http"
	constants "service-auth/Constants"
	"service-auth/DTO"
	initializers "service-auth/Initializers"
	auth_services "service-auth/Services/AuthServices"
	viewmodels "service-auth/ViewModels"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// Khai báo struct IntanceAuthController thông qua dependency injection (auth_services.AuthService)
type AuthController struct {
	authService *auth_services.AuthService
}

// khởi tạo intance NewAuthController định nghĩa struct AuthController
func NewAuthController(authService *auth_services.AuthService) *AuthController {
	return &AuthController{authService}
}

func (service *AuthController) SignUpController(c *gin.Context) {
	var signUpResponse viewmodels.SigUpViewModel //

	var bodyRequest DTO.SignUp_Request // khởi tạo bodyRequest

	if err := c.ShouldBindJSON(&bodyRequest); err != nil { // bind data từ request sang bodyRequest
		// write log
		initializers.Logger.Error("Function SignUpController()",
			zap.String("Error Bind JSON", err.Error()),
		)

		// set data ViewModel reponse to user
		signUpResponse.Code = constants.CODE_BAD_REQUEST
		signUpResponse.Status = constants.STATUS_BAD_REQUEST
		signUpResponse.Message = "Invalid JSON format."

		c.JSON(http.StatusBadRequest, signUpResponse)
		return
	}

	// 2.call service Auth execute function sign up
	_, errResponse, httpS := service.authService.SignUpAccount(bodyRequest)

	// Set Response trả về
	signUpResponse.BaseReponseDTO = errResponse

	c.JSON(httpS.HTTPStatus, signUpResponse)
}
