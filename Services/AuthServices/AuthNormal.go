package auth_services

import (
	"context"
	"fmt"
	"net/http"
	"os"
	constants "service-auth/Constants"
	"service-auth/DTO"
	helpers "service-auth/Helpers"
	initializers "service-auth/Initializers"
	repositories "service-auth/Repositories"
	"time"
)

// khởi tạo context background để run Redis
var ctx = context.Background()

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

func (repo *AuthService) SignInAccount(dataRequest DTO.SignIn_Request) (DTO.AuthService_SignIn_Response, DTO.BaseReponseDTO, DTO.HTTPStatusDTO) {
	var dataResponse DTO.AuthService_SignIn_Response // khởi tạo biến lưu giá trị trả về với stuct AuthService_SignIn_Response
	var err error                                    // khai báo biến trả về lỗi khi thực thi function này
	var errResponse DTO.BaseReponseDTO               // khai báo đối tượng trả về thông báo cho client khi thực thi function này
	var httpStatus DTO.HTTPStatusDTO                 // khai báo đối tượng trả về mã lỗi http cho request

	// Logic đăng nhập
	// 1. Kiểm tra tồn tại của email
	// Gọi function GetUserByEmail từ UserRepository
	userDetail, err := repo.userRepository.GetUserByEmail(dataRequest.Email)

	if err != nil { // email không tồn tại tồn tại -> thông báo mã lỗi và trả về kết quả failed
		// write log
		errJSON, _ := helpers.JSON_Stringify(err)
		objectLog := map[string]interface{}{
			"Error Find User By Email": errJSON,
		}
		helpers.WriteLogApp("Function SignInAccount() - AuthService", objectLog, "ERROR")

		// set dữ liệu cho errRespone
		errResponse.Code = constants.CODE_INVALID_FIELD
		errResponse.Status = constants.STATUS_INVALID_FIELD
		errResponse.Message = "Email hoặc mật khẩu không hợp lệ."

		// set trạng thái trả lỗi HTTPStatus
		httpStatus.HTTPStatus = http.StatusUnprocessableEntity
		return dataResponse, errResponse, httpStatus
	}

	// 2. Confirm password input với pash hash trong lưu trữ.
	errComparePassword := helpers.ComparePasswordByBcrypt(userDetail.Password, dataRequest.Password)

	if errComparePassword != nil { // xác thực mật khẩu không hợp lệ. password wrong
		// write log
		errJSON, _ := helpers.JSON_Stringify(err)
		objectLog := map[string]interface{}{
			"Error Compare Password ": errJSON,
			"User Compare ":           userDetail.Email,
		}
		helpers.WriteLogApp("Function SignInAccount() - AuthService", objectLog, "ERROR")

		// set dữ liệu cho errRespone
		errResponse.Code = constants.CODE_INVALID_FIELD
		errResponse.Status = constants.STATUS_INVALID_FIELD
		errResponse.Message = "Email hoặc mật khẩu không hợp lệ."
		return dataResponse, errResponse, httpStatus
	}

	// Kiểm tra tài khoản nếu chưa xác thực gửi OTP xác thực
	// Call đến API đến service Mail

	// 3. tạo token và thông tin user trả về cho người dùng
	var tokenJWT DTO.JWTToken

	tokenJWT.UserID = userDetail.UserID

	// lấy secret key trong .env
	accessKey := os.Getenv(constants.JWT_ACCESS_SECRET)
	refreshKey := os.Getenv(constants.JWT_REFRESH_SECRET)

	// divice + hash secret
	accesSignKey := accessKey + dataRequest.Device
	refreshSignKey := refreshKey + dataRequest.Device

	var errGenerate error
	accessToken, errGenerate := helpers.CreateAccessToken(tokenJWT, accesSignKey)

	if errGenerate != nil { // create Token failed.
		// write log
		errJSON, _ := helpers.JSON_Stringify(errGenerate)
		objectLog := map[string]interface{}{
			"Create Access Token Failed ": errJSON,
		}

		helpers.WriteLogApp("Function SignInAccount() - AuthService", objectLog, "ERROR")

		// set dữ liệu cho errRespone
		errResponse.Code = constants.CODE_BAD_REQUEST
		errResponse.Status = constants.STATUS_BAD_REQUEST
		errResponse.Message = "NETWORK ERROR."
		return dataResponse, errResponse, httpStatus
	}

	refreshToken, errGenerate := helpers.CreateRefreshToken(tokenJWT, refreshSignKey)

	if errGenerate != nil { // create Token failed.
		// write log
		errJSON, _ := helpers.JSON_Stringify(errGenerate)
		objectLog := map[string]interface{}{
			"Create Access Token Failed ": errJSON,
		}

		helpers.WriteLogApp("Function SignInAccount() - AuthService", objectLog, "ERROR")

		// set dữ liệu cho errRespone
		errResponse.Code = constants.CODE_BAD_REQUEST
		errResponse.Status = constants.STATUS_BAD_REQUEST
		errResponse.Message = "NETWORK ERROR."
		return dataResponse, errResponse, httpStatus
	}

	// lưu thông tin user lên Redis - lưu cache tránh gọi vào DB khi xác thực token
	userKey := fmt.Sprint(userDetail.UserID) // Tạo key với định danh người dùng
	userFields := map[string]interface{}{
		"user_id":     fmt.Sprint(userDetail.UserID),
		"email":       userDetail.Email,
		"user_status": fmt.Sprint(userDetail.UserStatus),
		"user_type":   fmt.Sprint(userDetail.UserType),
	}

	errSaveRedis := initializers.RedisUser.HMSet(ctx, userKey, userFields).Err()

	if errSaveRedis != nil {
		// write log
		errJSON, _ := helpers.JSON_Stringify(errSaveRedis)
		objectLog := map[string]interface{}{
			"Save data user on Redis failed": errJSON,
		}

		helpers.WriteLogApp("Function SignInAccount() - AuthService", objectLog, "ERROR")
	}

	ttl := time.Hour // Thời gian cache là 1 tiếng
	err = initializers.RedisUser.Expire(ctx, userKey, ttl).Err()

	if err != nil {
		// write log
		errJSON, _ := helpers.JSON_Stringify(err)
		objectLog := map[string]interface{}{
			"Could not set TTL": errJSON,
		}

		helpers.WriteLogApp("Function SignInAccount() - AuthService", objectLog, "ERROR")
	}

	// Lấy khóa bí mật từ biến môi trường
	hashSecretKey := os.Getenv(constants.DEVICE_SECRET_KEY)

	// Tạo chuỗi để băm bằng cách kết hợp UserID, Device và khóa bí mật
	passHASH := fmt.Sprint(userDetail.UserID) + dataRequest.Device + hashSecretKey

	// Tạo HMAC hash từ chuỗi đã tạo và UserID
	// Sử dụng hàm GenerateHMAC để tạo HMAC
	authKey := helpers.GenerateHMAC(passHASH, fmt.Sprint(userDetail.UserID))

	deviceFields := map[string]interface{}{
		"device":        dataRequest.Device,
		"ip_address":    dataRequest.IPAddress,
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	}

	// lưu trữ token lên Redis sử dụng xác thực request của ng dùng trong middleware
	err = initializers.RedisAuth.HMSet(ctx, authKey, deviceFields).Err()

	if err != nil {
		// write log
		errJSON, _ := helpers.JSON_Stringify(err)
		objectLog := map[string]interface{}{
			"Save data token on Redis failed": errJSON,
		}

		helpers.WriteLogApp("Function SignInAccount() - AuthService", objectLog, "ERROR")
	}

	ttlToken := time.Hour * 24 * 30 // Thời gian cache là 30 ngày = thời gian hết hạn refresh token
	err = initializers.RedisAuth.Expire(ctx, authKey, ttlToken).Err()

	if err != nil {
		// write log
		errJSON, _ := helpers.JSON_Stringify(err)
		objectLog := map[string]interface{}{
			"Could not set TTL": errJSON,
		}

		helpers.WriteLogApp("Function SignInAccount() - AuthService", objectLog, "ERROR")
	}

	// set trạng thái baseResponse
	errResponse.Code = constants.CODE_SUCCESS
	errResponse.Status = constants.STATUS_SUCCESS
	errResponse.Message = "Xác thực tài khoản thành công."

	// set thông tin cần thiết cho xác thực FE
	dataResponse.UserID = userDetail.UserID
	dataResponse.AccessToken = accessToken
	dataResponse.RefreshToken = refreshToken

	return dataResponse, errResponse, httpStatus
}
