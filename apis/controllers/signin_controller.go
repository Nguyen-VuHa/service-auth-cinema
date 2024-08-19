package controllers

import "github.com/gin-gonic/gin"

type SignInController struct {
}

func (sc *SignInController) SignIn(c *gin.Context) {
	// 1. convert từ body sang struct request.

	// 2. Kiểm tra dữ liệu nhất quán

	// 3. Kiểm tra tồn tại của email trong trên bộ nhớ cache Redis không thì kiểm tra trên database

	// 4. Confirm password với dữ liệu trong hệ thống.

	// 5. Kiểm tra xác thực OTP nếu chưa thì gửi OTP xác thực.

	// 6. Tạo token gửi về cho người dùng.

	// 7. Lưu thông tin cơ bản của người dùng lên Redis để caching data.

	//
	c.JSON(200, gin.H{"message": "hello"})
}
