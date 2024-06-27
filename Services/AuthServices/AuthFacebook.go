package auth_services

import (
	"net/http"
	"os"
	constants "service-auth/Constants"
	"service-auth/DTO"
	helpers "service-auth/Helpers"
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

func (repo *AuthService) CallbackSignUpWithFacebook(dataCallback DTO.Callback_SignIn_Facebook) (DTO.AuthService_Callback_Facebook_Response, DTO.BaseReponseDTO, DTO.HTTPStatusDTO) {
	var dataResponse DTO.AuthService_Callback_Facebook_Response // khởi tạo biến lưu giá trị trả về với stuct AuthService_Callback_Facebook_Response
	var err error                                               // khai báo biến trả về lỗi khi thực thi function này
	var baseResponse DTO.BaseReponseDTO                         // khai báo đối tượng trả về thông báo cho client khi thực thi function này
	var httpStatus DTO.HTTPStatusDTO                            // khai báo đối tượng trả về mã lỗi http cho request

	// Logic đăng nhập với facebook
	// 1. Kiểm tra tồn tại của email
	// Gọi function GetUserByEmail từ UserRepository
	_, err = repo.userRepository.GetUserByEmail(dataCallback.ID + "@facebook.com")

	// tạm thời chưa xử lý lưu lại tính tiếp
	if err == nil { // email tồn tại -> kiểm tra token còn hạng không nếu không bắt đăng nhập lại.
		// Xử lý xác thực hoặc trả về lỗi.
		// 	// write log
		// 	errJSON, _ := helpers.JSON_Stringify(err)
		// 	objectLog := map[string]interface{}{
		// 		"Error Find User By Email": errJSON,
		// 	}
		// 	helpers.WriteLogApp("Function SignUpAccount() - AuthService", objectLog, "ERROR")

		// 	// set dữ liệu cho errRespone
		// 	baseResponse.Code = constants.CODE_INVALID_FIELD
		// 	baseResponse.Status = constants.STATUS_INVALID_FIELD
		// 	baseResponse.Message = "Email đăng ký đã tồn tại."

		// 	// set trạng thái trả lỗi HTTPStatus
		// 	httpStatus.HTTPStatus = http.StatusUnprocessableEntity

		// 	return dataResponse, baseResponse, httpStatus
	} else {
		// nếu chưa tồn tại -> insert vào database

		// 3. insert thông tin vào database tương ứng
		err = repo.userRepository.CreateUserLoginFacebook(dataCallback)

		if err != nil {
			// write log
			errJSON, _ := helpers.JSON_Stringify(err)
			objectLog := map[string]interface{}{
				"Storage failed ": errJSON,
			}
			helpers.WriteLogApp("Function CallbackSignUpWithFacebook() - AuthService", objectLog, "ERROR")

			// set dữ liệu cho errRespone
			baseResponse.Code = constants.CODE_SERVER_INTERNAL_ERROR
			baseResponse.Status = constants.STATUS_SERVER_INTERNAL_ERROR
			baseResponse.Message = "INTERNAL SERVER ERROR."

			// set trạng thái trả lỗi HTTPStatus
			httpStatus.HTTPStatus = http.StatusInternalServerError
			return dataResponse, baseResponse, httpStatus
		}

	}

	return dataResponse, baseResponse, httpStatus
}
