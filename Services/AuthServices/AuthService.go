package auth_services

import (
	"net/http"
	constants "service-auth/Constants"
	"service-auth/DTO"
	helpers "service-auth/Helpers"
	repositories "service-auth/Repositories"
)

// Khai báo struct AuthService thông qua dependency injection (repositories.UserRepository)
type AuthService struct {
	userRepository repositories.UserRepository
}

// khởi tạo intance NewAuthService định nghĩa struct AuthService
func NewAuthService(userRepository repositories.UserRepository) *AuthService {
	return &AuthService{userRepository}
}

func (repo *AuthService) SignUpAccount(dataRequest DTO.SignUp_Request) (DTO.AuthService_SignUp_Response, DTO.BaseReponseDTO, DTO.HTTPStatusDTO) {
	var dataResponse DTO.AuthService_SignUp_Response // khởi tạo biến lưu giá trị trả về với stuct AuthService_SignUp_Response
	var err error                                    // khai báo biến trả về lỗi khi thực thi function này
	var errResponse DTO.BaseReponseDTO               // khai báo đối tượng trả về thông báo cho client khi thực thi function này
	var httpStatus DTO.HTTPStatusDTO                 // khai báo đối tượng trả về mã lỗi http cho request

	// Logic đăng ký
	// 1. Kiểm tra tồn tại của email
	// Gọi function GetUserByEmail từ UserRepository
	_, err = repo.userRepository.GetUserByEmail(dataRequest.Email)

	if err == nil { // email tồn tại -> thông báo mã lỗi và trả về kết quả failed
		// write log
		errJSON, _ := helpers.JSON_Stringify(err)
		objectLog := map[string]interface{}{
			"Error Find User By Email": errJSON,
		}
		helpers.WriteLogApp("Function SignUpAccount() - AuthService", objectLog, "ERROR")

		// set dữ liệu cho errRespone
		errResponse.Code = constants.CODE_INVALID_FIELD
		errResponse.Status = constants.STATUS_INVALID_FIELD
		errResponse.Message = "Email đăng ký đã tồn tại."

		// set trạng thái trả lỗi HTTPStatus
		httpStatus.HTTPStatus = http.StatusUnprocessableEntity
		return dataResponse, errResponse, httpStatus
	}

	// 2. hash password với thư viện bcrypt:
	var passwordHash = "" // khai báo biến lưu trữ kết quả tra về từ fucntion HashPasswordWithBcrypt
	passwordHash, err = helpers.HashPasswordWithBcrypt(dataRequest.Password)

	// set passwordHash vào lại object DTO.SignUp_Request
	dataRequest.Password = passwordHash

	if err != nil { // lỗi trong quá trình hash password ở function HashPasswordWithBcrypt
		// write log
		errJSON, _ := helpers.JSON_Stringify(err)
		objectLog := map[string]interface{}{
			"Hash Password Failed": errJSON,
		}
		helpers.WriteLogApp("Function SignUpAccount() - AuthService", objectLog, "ERROR")

		// set dữ liệu cho errRespone
		errResponse.Code = constants.CODE_SERVER_INTERNAL_ERROR
		errResponse.Status = constants.STATUS_SERVER_INTERNAL_ERROR
		errResponse.Message = "INTERNAL SERVER ERROR."

		// set trạng thái trả lỗi HTTPStatus
		httpStatus.HTTPStatus = http.StatusInternalServerError
		return dataResponse, errResponse, httpStatus
	}

	// 3. insert thông tin vào database tương ứng
	err = repo.userRepository.CreateNewUser(dataRequest)

	if err != nil {
		// write log
		errJSON, _ := helpers.JSON_Stringify(err)
		objectLog := map[string]interface{}{
			"Storage failed ": errJSON,
		}
		helpers.WriteLogApp("Function SignUpAccount() - AuthService", objectLog, "ERROR")

		// set dữ liệu cho errRespone
		errResponse.Code = constants.CODE_SERVER_INTERNAL_ERROR
		errResponse.Status = constants.STATUS_SERVER_INTERNAL_ERROR
		errResponse.Message = "INTERNAL SERVER ERROR."

		// set trạng thái trả lỗi HTTPStatus
		httpStatus.HTTPStatus = http.StatusInternalServerError
		return dataResponse, errResponse, httpStatus
	}

	// 4. trả về kết quả
	errResponse.Code = constants.CODE_SUCCESS
	errResponse.Status = constants.STATUS_SUCCESS
	errResponse.Message = "Đăng ký thành công."

	return dataResponse, errResponse, httpStatus
}
