package controllers

import (
	"auth-service/constants"
	"auth-service/domains"
	"auth-service/utils"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type SignupController struct {
	SignUpUseCase domains.SignUpUsecase
}

func (sc *SignupController) SignUp(c *gin.Context) {
	// 1. convert từ body sang struct request.
	var request domains.SignUpRequest
	var response domains.SignUpResponse

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
	err = sc.SignUpUseCase.ValidateDataRequest(request)

	if err != nil {
		// set data trả về
		response.Code = constants.CODE_BAD_REQUEST
		response.Status = constants.STATUS_BAD_REQUEST
		response.Message = err.Error()

		c.JSON(http.StatusBadRequest, response)
		return
	}

	// 3. Kiểm tra tồn tại của email trong hệ thống
	_, err = sc.SignUpUseCase.GetUserByEmail(request.Email)
	if err == nil {
		// set data trả về
		response.Code = constants.CODE_INVALID_FIELD
		response.Status = constants.STATUS_INVALID_FIELD
		response.Message = "Email đăng ký đã tồn tại."

		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	// 4. Mã hoá dữ liệu (password, ...)
	passwordHash, err := utils.HashPasswordWithBcrypt(request.Password)
	if err != nil {
		// set data trả về
		response.Code = constants.CODE_SERVER_INTERNAL_ERROR
		response.Status = constants.STATUS_SERVER_INTERNAL_ERROR
		response.Message = "INTERNAL SERVER ERROR."

		c.JSON(http.StatusInternalServerError, response)
		return
	}

	fmt.Println(passwordHash)
	// 4. Lưu trữ thông tin bảng user.

	// 5. Lưu trữ thông tin bảng user_profile
	c.JSON(http.StatusOK, response)
}
