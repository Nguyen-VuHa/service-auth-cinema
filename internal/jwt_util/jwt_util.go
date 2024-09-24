package jwt_util

import (
	"auth-service/domains"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

// JwtNewWithClaim tạo một JWT mới với các claim được truyền vào.
// Nó sử dụng phương pháp ký HS256 và trả về một con trỏ đến đối tượng jwt.Token.
func JwtNewWithClaim(dataClaim domains.JWTToken) *jwt.Token {
	// Tạo một token mới với phương pháp ký là HS256 (HMAC với SHA-256)
	// và các claim được truyền vào dưới dạng jwt.MapClaims.
	return jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		// "user_id" chứa ID người dùng từ dataClaim.
		"user_id": dataClaim.UserID,
		// "exp" chứa thời gian hết hạn của token,
		// được tính bằng cách cộng thời gian hiện tại với thời gian hết hạn trong dataClaim.
		// Hàm Unix() trả về thời gian dưới dạng số giây kể từ kỷ nguyên Unix (1/1/1970).
		"exp": time.Now().Add(dataClaim.Exp).Unix(),
	})
}

func CreateJWTToken(data domains.JWTToken, signKey string) (string, error) {
	newToken := JwtNewWithClaim(data) // Gọi hàm JwtNewWithClaim để tạo một token mới với các claim từ data.

	// Ký token mới với khóa bí mật.
	// Hàm SignedString sẽ trả về token đã ký dưới dạng chuỗi.
	signToken, err := newToken.SignedString([]byte(signKey))

	// Kiểm tra nếu có lỗi xảy ra trong quá trình ký token.
	if err != nil {
		return "", err // Nếu có lỗi, trả về chuỗi rỗng và lỗi đó.
	}

	return signToken, nil // Trả về token đã ký và nil cho lỗi (không có lỗi).
}

func VerifyJWTToken(tokenString string, signKey string) (string, error) {
	// Parse token và xác thực chữ ký với phương pháp ký HS256.
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Kiểm tra xem phương pháp ký có đúng là HS256 không.
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		// Trả về khóa bí mật để xác thực chữ ký.
		return []byte(signKey), nil
	})

	// Kiểm tra nếu có lỗi trong quá trình giải mã token.
	if err != nil {
		return "", err
	}

	// Kiểm tra tính hợp lệ của token (có hợp lệ về chữ ký và chưa hết hạn hay không).
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		float_unix := claims["exp"].(float64)

		// Chuyển đổi float64 sang int64 (Unix time theo giây)
		seconds := int64(float_unix)

		// Lấy thời gian hiện tại dưới dạng Unix timestamp
		now := time.Now().Unix()

		if now > seconds {
			return "", errors.New("token expired")
		}

		return claims["user_id"].(string), nil
	} else {
		return "", fmt.Errorf("invalid token")
	}
}

// CreateAccessToken tạo một token truy cập (access token) với thời gian hết hạn 30 phút.
// Nó sử dụng hàm CreateJWTToken để tạo và ký token.
func CreateAccessToken(dataToken domains.JWTToken, signKey string) (string, error) {
	// Thiết lập thời gian hết hạn cho token là 30 phút.
	// Bạn có thể thay đổi giá trị này nếu cần.
	dataToken.Exp = time.Minute * 30

	// Gọi hàm CreateJWTToken để tạo và ký token truy cập với khóa bí mật
	accessToken, errAccessToken := CreateJWTToken(dataToken, signKey)

	// Kiểm tra nếu có lỗi xảy ra trong quá trình tạo và ký token truy cập.
	if errAccessToken != nil {
		// Trả về một chuỗi rỗng và một thông báo lỗi "NETWORK ERROR".
		return "", errAccessToken
	}

	// Trả về token truy cập đã tạo và nil cho lỗi (không có lỗi).
	return accessToken, nil
}

// CreateRefreshToken tạo một refresh token với thời gian hết hạn là 30 ngày.
// Nó sử dụng hàm CreateJWTToken để tạo và ký token.
func CreateRefreshToken(dataToken domains.JWTToken, signKey string) (string, error) {
	// Thiết lập thời gian hết hạn cho refresh token là 30 ngày.
	dataToken.Exp = time.Hour * 24 * 30

	// Gọi hàm CreateJWTToken để tạo và ký refresh token với khóa bí mật
	refreshToken, errRefreshToken := CreateJWTToken(dataToken, signKey)

	// Kiểm tra nếu có lỗi xảy ra trong quá trình tạo và ký refresh token.
	if errRefreshToken != nil {
		return "", errRefreshToken // Trả về một chuỗi rỗng và một thông báo lỗi "NETWORK ERROR".
	}

	// Trả về refresh token đã tạo và nil cho lỗi (không có lỗi).
	return refreshToken, nil
}
