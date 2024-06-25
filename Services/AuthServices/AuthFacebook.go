package auth_services

import (
	"net/http"
	"os"
	constants "service-auth/Constants"
	"service-auth/DTO"
	initializers "service-auth/Initializers"
)

func (repo *AuthService) SignUpWithFacebook(dataRequest DTO.SignInFacebook_Request) (DTO.AuthService_SignInFacebook_Response, DTO.BaseReponseDTO, DTO.HTTPStatusDTO) {
	var dataResponse DTO.AuthService_SignInFacebook_Response // khởi tạo biến lưu giá trị trả về với stuct AuthService_SignInFacebook_Response
	// var err error                                            // khai báo biến trả về lỗi khi thực thi function này
	var errResponse DTO.BaseReponseDTO // khai báo đối tượng trả về thông báo cho client khi thực thi function này
	var httpStatus DTO.HTTPStatusDTO   // khai báo đối tượng trả về mã lỗi http cho request

	signKeyFacebook := os.Getenv(constants.FACEBOOK_SIGN_KEY)
	url := initializers.FacebookConfig.AuthCodeURL(signKeyFacebook)

	// 4. trả về kết quả
	errResponse.Code = constants.CODE_TEMPORARY_REDIRECT
	errResponse.Status = constants.STATUS_TEMPORARY_REDIRECT
	errResponse.Message = "Xác thực thành công."
	dataResponse.URL = url

	httpStatus.HTTPStatus = http.StatusTemporaryRedirect

	return dataResponse, errResponse, httpStatus
}
