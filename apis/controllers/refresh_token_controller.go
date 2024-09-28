package controllers

import (
	"auth-service/constants"
	"auth-service/domains"
	"auth-service/utils"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type RefreshTokenController struct {
	RefreshTokenUsecase domains.RefreshTokenUsecase
}

func (rfc *RefreshTokenController) RefreshToken(c *gin.Context) {
	var data_repsonse domains.RefreshTokenResponse

	refresh_token := c.Query("_token")
	user_id := c.Query("_user_id")

	if refresh_token == "" || user_id == "" {
		data_repsonse.Code = constants.CODE_BAD_REQUEST
		data_repsonse.Status = constants.STATUS_BAD_REQUEST
		data_repsonse.Message = "Params invalid."

		c.JSON(http.StatusBadRequest, data_repsonse)
		return
	}

	token_arr := strings.Split(refresh_token, ".")

	// kiểm tra chỉ method thuộc các phương thức đặt biệt mới xác thực.
	arr_method := []string{fmt.Sprint(constants.LOGIN_GOOGLE_ID)}

	mothod_is_valid := utils.ItemIsArrayString(token_arr[0], arr_method)

	if mothod_is_valid { // Nếu refresh token là id của phương thức thì refresh cho method đó thường là Google
		is_token_valid, err := rfc.RefreshTokenUsecase.ValidateTokenGoogle(user_id, strings.Join(token_arr[1:], "."))

		if err != nil {
			data_repsonse.Code = constants.CODE_BAD_REQUEST
			data_repsonse.Status = constants.STATUS_BAD_REQUEST
			data_repsonse.Message = "Refresh token failed."

			c.JSON(http.StatusBadRequest, data_repsonse)
			return
		}

		if !is_token_valid {
			data_repsonse.Code = constants.CODE_BAD_REQUEST
			data_repsonse.Status = constants.STATUS_BAD_REQUEST
			data_repsonse.Message = "Refresh token failed."

			c.JSON(http.StatusBadRequest, data_repsonse)
			return
		}

		access_token, err := rfc.RefreshTokenUsecase.CreateRefreshTokenGoogle(user_id)

		if err != nil {
			data_repsonse.Code = constants.CODE_BAD_REQUEST
			data_repsonse.Status = constants.STATUS_BAD_REQUEST
			data_repsonse.Message = "Refresh token failed."

			c.JSON(http.StatusBadRequest, data_repsonse)
			return
		}

		data_repsonse.Data = access_token

		data_repsonse.Code = constants.CODE_SUCCESS
		data_repsonse.Status = constants.STATUS_SUCCESS
		data_repsonse.Message = "Refresh token success"
	} else {
		var data_request domains.RefreshTokenRequest

		data_request.UserID = user_id
		data_request.RefreshToken = refresh_token
		data_request.Device = utils.GetDevice(c)

		err := rfc.RefreshTokenUsecase.ValidateRefreshToken(data_request)

		if err != nil {
			data_repsonse.Code = constants.CODE_BAD_REQUEST
			data_repsonse.Status = constants.STATUS_BAD_REQUEST
			data_repsonse.Message = "Refresh token failed."

			c.JSON(http.StatusBadRequest, data_repsonse)
			return
		}

		access_token, err := rfc.RefreshTokenUsecase.CreateRefreshToken(data_request)

		if err != nil {
			data_repsonse.Code = constants.CODE_BAD_REQUEST
			data_repsonse.Status = constants.STATUS_BAD_REQUEST
			data_repsonse.Message = "Refresh token failed."

			c.JSON(http.StatusBadRequest, data_repsonse)
			return
		}

		data_repsonse.Data = access_token

		data_repsonse.Code = constants.CODE_SUCCESS
		data_repsonse.Status = constants.STATUS_SUCCESS
		data_repsonse.Message = "Refresh token success"
	}

	c.JSON(http.StatusOK, data_repsonse)
}
