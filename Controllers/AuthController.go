package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	constants "service-auth/Constants"
	"service-auth/DTO"
	helpers "service-auth/Helpers"
	initializers "service-auth/Initializers"
	auth_services "service-auth/Services/AuthServices"
	validation_service "service-auth/Services/ValidationService"
	validations "service-auth/Validations"
	viewmodels "service-auth/ViewModels"

	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
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
	var signUpResponse viewmodels.SigUpViewModel // Khởi tạo data response SigUpViewModel

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

	// gọi hàm Valid_Auth_SignUp() để kiểm tra dữ liệu trong Body Request
	messageError := validations.Valid_Auth_SignUp(bodyRequest, service.validationService)

	if messageError != "" {
		// set data ViewModel reponse to user
		signUpResponse.Code = constants.CODE_BAD_REQUEST
		signUpResponse.Status = constants.STATUS_BAD_REQUEST
		signUpResponse.Message = messageError

		c.JSON(http.StatusBadRequest, signUpResponse)
		return
	}

	// 2.call service Auth execute function sign up
	_, errResponse, httpS := service.authService.SignUpAccount(bodyRequest)

	// Set Response trả về
	signUpResponse.BaseReponseDTO = errResponse

	c.JSON(httpS.HTTPStatus, signUpResponse)
}

func (service *AuthController) SignInCotroller(c *gin.Context) {
	var signInResponse viewmodels.SignInViewModel // Khởi tạo data response SigInViewModel

	var bodyRequest DTO.SignIn_Request // khởi tạo bodyRequest

	if err := c.ShouldBindJSON(&bodyRequest); err != nil { // bind data từ request sang bodyRequest
		// write log
		objectLog := map[string]interface{}{
			"Error Bind JSON": err.Error(),
		}

		helpers.WriteLogApp("Function SignInCotroller() - AuthController", objectLog, "ERROR")

		// set data ViewModel reponse to user
		signInResponse.Code = constants.CODE_BAD_REQUEST
		signInResponse.Status = constants.STATUS_BAD_REQUEST
		signInResponse.Message = "Invalid JSON format."

		c.JSON(http.StatusBadRequest, signInResponse)
		return
	}

	// gọi hàm Valid_Auth_SignIn() để kiểm tra dữ liệu trong Body Request
	messageError := validations.Valid_Auth_SignIn(bodyRequest, service.validationService)

	if messageError != "" {
		// set data ViewModel reponse to user
		signInResponse.Code = constants.CODE_BAD_REQUEST
		signInResponse.Status = constants.STATUS_BAD_REQUEST
		signInResponse.Message = messageError

		c.JSON(http.StatusBadRequest, signInResponse)
		return
	}

	// lấy thông tin device và IP đăng nhập
	ipRequest := helpers.GetClientIP(c)
	deviceRequest := helpers.GetDevice(c)

	// set thông tin ipaddr vs device vào bodyRequest
	bodyRequest.IPAddress = ipRequest
	bodyRequest.Device = deviceRequest

	data, baseResponse, httpS := service.authService.SignInAccount(bodyRequest)

	// Set Response trả về
	signInResponse.BaseReponseDTO = baseResponse
	signInResponse.Data = data

	c.JSON(httpS.HTTPStatus, signInResponse)
}

func (service *AuthController) SignInWithFacebookCotroller(c *gin.Context) {
	var facebookResponse viewmodels.SignInFacebookViewModel // Khởi tạo data response SigInViewModel

	var bodyRequest DTO.SignInFacebook_Request // khởi tạo bodyRequest

	if err := c.ShouldBindJSON(&bodyRequest); err != nil { // bind data từ request sang bodyRequest
		// write log
		objectLog := map[string]interface{}{
			"Error Bind JSON": err.Error(),
		}

		helpers.WriteLogApp("Function SignInCotroller() - AuthController", objectLog, "ERROR")

		// set data ViewModel reponse to user
		facebookResponse.Code = constants.CODE_BAD_REQUEST
		facebookResponse.Status = constants.STATUS_BAD_REQUEST
		facebookResponse.Message = "Invalid JSON format."

		c.JSON(http.StatusBadRequest, facebookResponse)
		return
	}

	data, baseResponse, _ := service.authService.SignUpWithFacebook(bodyRequest)

	// Set Response trả về
	facebookResponse.BaseReponseDTO = baseResponse
	facebookResponse.Data = data

	// c.JSON(httpS.HTTPStatus, facebookResponse.Data.URL)
	c.JSON(http.StatusOK, facebookResponse.Data.URL)
}

func (service *AuthController) CallBackWithFacebookCotroller(c *gin.Context) {
	code := c.Query("code")

	fmt.Println(code)

	token, err := initializers.FacebookConfig.Exchange(oauth2.NoContext, code)
	if err != nil {
		// http.Error(w, "Failed to exchange token: "+err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Println(token)

	client := initializers.FacebookConfig.Client(oauth2.NoContext, token)
	response, err := client.Get("https://graph.facebook.com/me?fields=id,name,email")
	if err != nil {
		// http.Error(w, "Failed to get user info: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer response.Body.Close()

	var user struct {
		ID    string `json:"id"`
		Name  string `json:"name"`
		Email string `json:"email"`
	}

	if err := json.NewDecoder(response.Body).Decode(&user); err != nil {
		// http.Error(w, "Failed to decode user info: "+err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Printf("User Info: %s\n", user)
}
