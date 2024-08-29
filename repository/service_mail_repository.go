package repository

import (
	"auth-service/constants"
	"auth-service/domains"
	"auth-service/utils"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
)

type serviceMailRepository struct {
}

func NewServiceMailRepository() domains.ServiceMailRepository {
	return &serviceMailRepository{}
}

func (smr *serviceMailRepository) SendOTPCodeToMail(params map[string]interface{}) error {
	// URL của API bạn muốn gọi
	url := os.Getenv(constants.URL_API_SERVICE) + constants.PATH_SEND_EMAIL_OTP

	urlSendOTP, err := utils.AddParamsToURL(url, params)

	if err != nil {
		return err
	}

	// Tạo một request mới
	req, err := http.NewRequest("GET", urlSendOTP, nil)

	if err != nil {
		return err
	}

	// Tạo một HTTP client
	client := &http.Client{}

	// Thực hiện request
	response, err := client.Do(req)

	if err != nil {
		fmt.Println("Error making GET request:", err)
		return err
	}

	defer response.Body.Close()

	// Đọc toàn bộ nội dung của response body
	body, err := io.ReadAll(response.Body)

	if err != nil {
		fmt.Println("Error reading response body:", err)
		return err
	}

	// Khai báo map để chứa kết quả
	var result map[string]interface{}

	// Chuyển đổi từ JSON sang map[string]interface{}
	err = json.Unmarshal([]byte(string(body)), &result)
	if err != nil {

		fmt.Println("Error:", err)
		return err
	}

	if fmt.Sprint(result["code"]) != fmt.Sprint(constants.CODE_SUCCESS) {
		return errors.New(fmt.Sprint(result["message"]))
	}

	return nil
}
