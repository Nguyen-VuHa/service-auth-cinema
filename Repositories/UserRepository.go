package repositories

import (
	data_layers "service-auth/DataLayers"
	models "service-auth/Models"
)

// khởi tạo interface repository cho phần UserRepository
type UserRepository interface {
	// Function Get Access Database
	GetUserByEmail(email string) (models.User, error)

	// Function Edit Access Database

	// Function Logic
}

// Khai báo struct IntanceUserDataLayer thông qua dependency injection (InterfaceUserDataLayer)
type IntanceUserDataLayer struct {
	userData *data_layers.UserDataLayer
}

// khởi tạo intance NewIntanceUserDataLayer chưa struct IntanceUserDataLayer
func NewIntanceUserDataLayer(userData *data_layers.UserDataLayer) *IntanceUserDataLayer {
	return &IntanceUserDataLayer{userData}
}

// GetUserByEmail truy xuất thông tin user với Email truyền vào
func (intance *IntanceUserDataLayer) GetUserByEmail(email string) (models.User, error) {
	var userData models.User // Khai báo biến để chứa thông tin người dùng truy xuất được
	var err error            // Khai báo biến để chứa lỗi trong quá trình thực thi

	// khởi tạo object confition cần thiết
	var condition = map[string]interface{}{
		"email": email,
	}

	userData, err = intance.userData.GetUserByConditions(condition)

	return userData, err
}
