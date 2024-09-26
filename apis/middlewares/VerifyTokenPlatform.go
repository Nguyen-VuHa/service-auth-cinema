package middlewares

import (
	"auth-service/constants"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

func VerifyTokenFacebook(userAccessToken, appAccessToken string) (bool, error) {
	// Xây dựng URL để kiểm tra token
	url := fmt.Sprintf("https://graph.facebook.com/debug_token?input_token=%s&access_token=%s", userAccessToken, appAccessToken)

	// Gửi yêu cầu GET đến Facebook
	resp, err := http.Get(url)
	if err != nil {
		return false, fmt.Errorf("error sending request: %v", err)
	}
	defer resp.Body.Close()

	// Kiểm tra trạng thái phản hồi
	if resp.StatusCode != http.StatusOK {
		return false, fmt.Errorf("invalid response from Facebook: %d", resp.StatusCode)
	}

	// Phân tích phản hồi từ Facebook
	var result map[string]interface{}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return false, fmt.Errorf("invalid response from Facebook: %d", err)
	}

	// Kiểm tra xem "data" có tồn tại trong map hay không
	data, ok := result["data"].(map[string]interface{})
	if !ok {
		return false, fmt.Errorf("key 'data' not found or invalid type")
	}

	// Kiểm tra xem "is_valid" có tồn tại trong "data" map hay không
	isValid, ok := data["is_valid"].(bool)
	if !ok {
		return false, fmt.Errorf("key 'is_valid' not found or invalid type")
	}

	return isValid, nil

}

// TokenInfo chứa thông tin phản hồi từ Google
type GoogleTokenInfo struct {
	Iss   string `json:"iss"`
	Sub   string `json:"sub"`
	Aud   string `json:"aud"`
	Email string `json:"email"`
	Exp   string `json:"exp"`
	// Thêm các trường cần thiết khác
}

func VerifyTokenGoogle(token string) (bool, error) {
	url := fmt.Sprintf("https://oauth2.googleapis.com/tokeninfo?id_token=%s", token)

	resp, err := http.Get(url)

	if err != nil {
		fmt.Println(err)
		return false, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return false, fmt.Errorf("failed to verify token, status code: %d", resp.StatusCode)
	}

	var tokenInfo GoogleTokenInfo
	err = json.NewDecoder(resp.Body).Decode(&tokenInfo)

	if err != nil {
		fmt.Println(err)
		return false, err
	}

	client_ID := os.Getenv(constants.GOOGLE_CLIENT_ID)
	// Kiểm tra xem trường `aud` có khớp với `Client ID` hay không
	if tokenInfo.Aud != client_ID {
		return false, fmt.Errorf("failed to verify token")
	}

	// Kiểm tra aud và các điều kiện khác ở đây
	fmt.Println(tokenInfo)
	return true, nil
}
