package controllers

import (
	"auth-service/constants"
	"auth-service/domains"
	"auth-service/utils"
	"net/http"

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

	c.JSON(http.StatusOK, data_repsonse)
}
