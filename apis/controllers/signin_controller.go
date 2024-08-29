package controllers

import (
	"auth-service/constants"
	"auth-service/domains"
	"auth-service/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

type SignInController struct {
	SignInUsecase domains.SignInUsecase
}

func (sc *SignInController) SignIn(c *gin.Context) {
	// 1. convert từ body sang struct request.
	var request domains.SignInRequest
	var response domains.SignInResponse

	// 1. convert từ body sang struct request.
	err := c.ShouldBind(&request)
	if err != nil {
		// set data trả về
		response.Code = constants.CODE_BAD_REQUEST
		response.Status = constants.STATUS_BAD_REQUEST
		response.Message = "Invalid JSON format."

		c.JSON(http.StatusBadRequest, response)
		return
	}

	// 2. Kiểm tra dữ liệu nhất quán
	err = sc.SignInUsecase.ValidateDataRequest(request)

	if err != nil {
		// set data trả về
		response.Code = constants.CODE_BAD_REQUEST
		response.Status = constants.STATUS_BAD_REQUEST
		response.Message = err.Error()

		c.JSON(http.StatusBadRequest, response)
		return
	}

	// 2.1 get IP client và device gán vào data request
	// lấy thông tin device và IP đăng nhập
	request.IpAddress = utils.GetClientIP(c)
	request.Device = utils.GetDevice(c)

	// 3. Kiểm tra tồn tại của email trên database
	userData, passwordHash, err := sc.SignInUsecase.GetUserByEmail(request.Email)

	if err != nil {
		// set data trả về
		response.Code = constants.CODE_INVALID_FIELD
		response.Status = constants.STATUS_INVALID_FIELD
		response.Message = "Email hoặc mật khẩu không hợp lệ."

		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	// 4. Confirm password với dữ liệu trong hệ thống.
	err = sc.SignInUsecase.ComparePasswordUser(passwordHash, request.Password)

	if err != nil {
		// set data trả về
		response.Code = constants.CODE_INVALID_FIELD
		response.Status = constants.STATUS_INVALID_FIELD
		response.Message = "Email hoặc mật khẩu không hợp lệ."

		c.JSON(http.StatusInternalServerError, response)
		return
	}

	// 5. Kiểm tra xác thực OTP nếu chưa thì gửi OTP xác thực.
	err = sc.SignInUsecase.CheckAccountVerification(userData, request)

	if err != nil {
		// set data trả về
		response.Code = constants.CODE_INVALID_FIELD
		response.Status = constants.STATUS_INVALID_FIELD
		response.Message = "NETWORK ERROR."

		c.JSON(http.StatusInternalServerError, response)
		return
	}

	// 6. Tạo token gửi về cho người dùng +  Lưu thông tin cơ bản của người dùng lên Redis để caching data.
	dataResponse, err := sc.SignInUsecase.CreateTokenAndDataResponse(userData, request)

	if err != nil {
		// set data trả về
		response.Code = constants.CODE_BAD_REQUEST
		response.Status = constants.STATUS_BAD_REQUEST
		response.Message = "NETWORK ERROR."

		c.JSON(http.StatusInternalServerError, response)
		return
	}

	// set data trả về
	response.Code = constants.CODE_SUCCESS
	response.Status = constants.STATUS_SUCCESS
	response.Message = "Login success."
	response.Data = dataResponse

	c.JSON(http.StatusOK, response)
}
