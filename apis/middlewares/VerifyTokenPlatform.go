package middlewares

import (
	"encoding/json"
	"fmt"
	"net/http"
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
