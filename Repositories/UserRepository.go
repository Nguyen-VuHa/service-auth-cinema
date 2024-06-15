package repositories

import (
	"service-auth/DTO"
	data_layers "service-auth/DataLayers"
	models "service-auth/Models"
)

// khởi tạo interface repository cho phần UserRepository
type UserRepository interface {
	// Function Get Access Database
	GetUserByEmail(email string) (models.User, error)

	// Function Edit Access Database
	CreateNewUser(userDataRequest DTO.SignUp_Request) error

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

	// GetUserByConditions thuộc lớp UserDataLayer
	// Thực thi funtion GetUserByConditions với object condition
	userData, err = intance.userData.GetUserByConditions(condition)

	// trả về kết quả thực thi và lỗi (nếu có)
	return userData, err
}

// CreateNewUser xử lý tạo mới 1 user trong hệ thống với dữ liệu userData được truyền vào.
func (intance *IntanceUserDataLayer) CreateNewUser(userDataRequest DTO.SignUp_Request) error {
	var err error            // Khai báo biến để chứa lỗi trong quá trình thực thi
	var userData models.User // Khai báo biến để chứa thông tin người dùng ghi được

	// set userData từ userDataRequest
	userData.Email = userDataRequest.Email
	userData.Password = userDataRequest.Password
	// thêm một số trường với rule khi tạo mới tài khoản

	var actionCreateUser = data_layers.CreateUserExecute{
		Data: &userData,
	}

	err = intance.userData.UserExecute(&actionCreateUser)

	// trả về lỗi (nếu có)
	return err
}
