package data_layers

import (
	"fmt"
	models "service-auth/Models"

	"gorm.io/gorm"
)

// Khai báo struct UserDataLayer thông qua dependency injection (DB)
type UserDataLayer struct {
	DB *gorm.DB
}

// khởi tạo một instance cho UserDataLayer
func NewUserDataLayer(DB *gorm.DB) *UserDataLayer {
	return &UserDataLayer{DB}
}

// GetUserByConditions truy xuất người dùng dựa trên nhiều điều kiện động.
// Các điều kiện được cung cấp dưới dạng map, trong đó key là tên cột và value là giá trị tương ứng cần khớp.
func (data *UserDataLayer) GetUserByConditions(conditions map[string]interface{}) (models.User, error) {
	var userData models.User // Khai báo biến để chứa thông tin người dùng truy xuất được
	var err error            // Khai báo biến để chứa lỗi trong quá trình thực thi

	query := data.DB.Model(&models.User{}) // Bắt đầu với truy vấn cơ bản từ kết nối cơ sở dữ liệu

	// Duyệt qua map các điều kiện
	for key, value := range conditions {
		// Xây dựng truy vấn động với các điều kiện
		// fmt.Sprintf("%s = ?", key) tạo một chuỗi điều kiện như "column_name = ?"
		query = query.Where(fmt.Sprintf("%s = ?", key), value)
	}

	// Thực hiện truy vấn và lưu kết quả vào biến userData
	// .Error trả về lỗi, nếu có, xảy ra trong quá trình thực hiện truy vấn
	err = query.First(&userData).Error

	// Trả về người dùng truy xuất được và lỗi nếu có
	return userData, err
}
