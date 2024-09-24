package controllers

import (
	"auth-service/constants"
	"auth-service/domains"
	"net/http"

	"github.com/gin-gonic/gin"
)

type GetDetailUserController struct {
	GetDetailUserUsecase domains.GetDetailUserUsecase
}

func (gduc *GetDetailUserController) GetDetailByID(c *gin.Context) {
	var data_response domains.GetDetailUserResponse

	user_id := c.Query("_user_id")

	if user_id == "" {
		data_response.Code = constants.CODE_BAD_REQUEST
		data_response.Status = constants.STATUS_BAD_REQUEST
		data_response.Message = "Params invalid."

		c.JSON(http.StatusBadRequest, data_response)
		return
	}

	data, err := gduc.GetDetailUserUsecase.GetDetailUserOnRedis(user_id)

	if err != nil {
		data_response.Code = constants.CODE_BAD_REQUEST
		data_response.Status = constants.STATUS_BAD_REQUEST
		data_response.Message = "Get detail user failed."

		c.JSON(http.StatusBadRequest, data_response)
		return
	}

	// trường hợp user không tồn tại trn Redis phải kiểm tra trong database.
	if data.Email == "" {
		data, err = gduc.GetDetailUserUsecase.GetDetailUserOnDatabase(user_id)

		if err != nil {
			data_response.Code = constants.CODE_BAD_REQUEST
			data_response.Status = constants.STATUS_BAD_REQUEST
			data_response.Message = "Get detail user failed."

			c.JSON(http.StatusBadRequest, data_response)
			return
		}
	}

	data_response.Code = constants.CODE_SUCCESS
	data_response.Status = constants.STATUS_SUCCESS
	data_response.Message = "Get detail user success."
	data_response.Data = data

	c.JSON(http.StatusOK, data_response)
}
