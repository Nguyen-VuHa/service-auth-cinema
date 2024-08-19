package controllers

import (
	"auth-service/constants"
	"auth-service/domains"
	"auth-service/models"
	"auth-service/utils"
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

	// 4. Lưu trữ thông tin bảng user.
	var userData models.User

	userData.Email = request.Email
	userData.Password = passwordHash
	userData.UserStatus = constants.USER_STATUS_PENDING
	userData.LoginMethodID = constants.LOGIN_NORMAL_ID

	err = sc.SignUpUseCase.CreateUser(&userData)

	if err != nil {
		// set data trả về
		response.Code = constants.CODE_SERVER_INTERNAL_ERROR
		response.Status = constants.STATUS_SERVER_INTERNAL_ERROR
		response.Message = "INTERNAL SERVER ERROR."

		c.JSON(http.StatusInternalServerError, response)
		return
	}

	// 5. Lưu trữ thông tin bảng user_profile
	var dataUserProfile = make(map[string]interface{})
	dataUserProfile[constants.USER_PROFILE_FULLNAME] = request.FullName
	dataUserProfile[constants.USER_PROFILE_BIRTHDAY] = request.BirthDay
	dataUserProfile[constants.USER_PROFILE_PHONENUMBER] = request.PhoneNumber

	var profileKeys = []string{constants.USER_PROFILE_FULLNAME, constants.USER_PROFILE_BIRTHDAY, constants.USER_PROFILE_PHONENUMBER}

	// insert thông tin vào user profile với các field còn lại
	for _, key := range profileKeys { // 3 số biến object cần lưu vào user profile (FullName, BirthDay, PhoneNumber)
		var userProfileData models.UserProfile // Khai báo biến để chứa thông tin detail user hợp lệ
		userProfileData.ProfileKey = key
		userProfileData.ProfileValue = dataUserProfile[key].(string)
		userProfileData.UserID = userData.UserID // Gán UserID khoá ngoại trong UserProfile

		err = sc.SignUpUseCase.CreateUserProfile(&userProfileData)

		if err != nil {
			break
		}
	}

	if err != nil {
		// set data trả về
		response.Code = constants.CODE_SERVER_INTERNAL_ERROR
		response.Status = constants.STATUS_SERVER_INTERNAL_ERROR
		response.Message = "INTERNAL SERVER ERROR."

		c.JSON(http.StatusInternalServerError, response)
	}

	response.Code = constants.CODE_SUCCESS
	response.Status = constants.STATUS_SUCCESS
	response.Message = "Tạo tài khoản thành công."

	c.JSON(http.StatusOK, response)
}
