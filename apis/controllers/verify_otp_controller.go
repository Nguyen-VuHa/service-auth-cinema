package controllers

import (
	"auth-service/constants"
	"auth-service/domains"
	"auth-service/utils"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type VerifyOTPController struct {
	VerifyOTPUsecase domains.VerifyOTPUsecase
}

func (votp *VerifyOTPController) VerifyOTP(c *gin.Context) {
	var data_response domains.VerifyOTPResponse

	params_email := c.Query("_email")
	params_otp := c.Query("_otp")

	// kiểm tra params truyền lên
	if params_email == "" || params_otp == "" {
		data_response.Code = constants.CODE_BAD_REQUEST
		data_response.Status = constants.STATUS_BAD_REQUEST
		data_response.Message = "Params invalid"

		c.JSON(http.StatusBadRequest, data_response)
		return
	}

	var verify_otp_request domains.VerifyOTPRequest

	verify_otp_request.Email = params_email
	verify_otp_request.OTP = params_otp

	// 2.1 get IP client và device gán vào data request
	// lấy thông tin device và IP đăng nhập
	verify_otp_request.IpAddress = utils.GetClientIP(c)
	verify_otp_request.Device = utils.GetDevice(c)

	user_data, is_email_valid, err := votp.VerifyOTPUsecase.ValidateEmail(params_email)

	if err != nil || !is_email_valid {
		fmt.Println(err)

		data_response.Code = constants.CODE_BAD_REQUEST
		data_response.Status = constants.STATUS_BAD_REQUEST
		data_response.Message = err.Error()

		c.JSON(http.StatusBadRequest, data_response)
		return
	}

	is_valid_otp, err := votp.VerifyOTPUsecase.CheckOTPValid(verify_otp_request, user_data)

	if err != nil {
		fmt.Println(err)

		data_response.Code = constants.CODE_BAD_REQUEST
		data_response.Status = constants.STATUS_BAD_REQUEST
		data_response.Message = "NETWORK ERROR."

		c.JSON(http.StatusBadRequest, data_response)
		return
	}

	if !is_valid_otp {
		data_response.Code = constants.CODE_BAD_REQUEST
		data_response.Status = constants.STATUS_BAD_REQUEST
		data_response.Message = "Xác thực OTP thất bại."

		c.JSON(http.StatusBadRequest, data_response)
		return
	}

	// update status and create token for user
	data, err := votp.VerifyOTPUsecase.UpdateUser(verify_otp_request, user_data)

	if err != nil {
		fmt.Println(err)

		data_response.Code = constants.CODE_BAD_REQUEST
		data_response.Status = constants.STATUS_BAD_REQUEST
		data_response.Message = "NETWORK ERROR."

		c.JSON(http.StatusBadRequest, data_response)
		return
	}

	data_response.Data = data

	data_response.Code = constants.CODE_SUCCESS
	data_response.Status = constants.STATUS_SUCCESS
	data_response.Message = "Verify OTP success."

	c.JSON(http.StatusOK, data_response)
}
