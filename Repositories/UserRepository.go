package repositories

import (
	constants "service-auth/Constants"
	"service-auth/DTO"
	user_data_layer "service-auth/DataLayers/User"
	models "service-auth/Models"
)

// khởi tạo interface repository cho phần UserRepository
type UserRepository interface {
	// Function Get Access Database
	GetUserByEmail(email string) (models.User, error)

	// Function Edit Access Database
	CreateNewUser(userDataRequest DTO.SignUp_Request) error
	CreateUserLoginFacebook(userDataRequest DTO.Callback_SignIn_Facebook) error

	// Function Logic
}

// Khai báo struct IntanceUserDataLayer thông qua dependency injection (InterfaceUserDataLayer)
type IntanceUserDataLayer struct {
	userData *user_data_layer.UserDataLayer
}

// khởi tạo intance NewIntanceUserDataLayer chưa struct IntanceUserDataLayer
func NewIntanceUserDataLayer(userData *user_data_layer.UserDataLayer) *IntanceUserDataLayer {
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
	userData.UserStatus = constants.USER_STATUS_PENDING
	userData.LoginMethodID = constants.LOGIN_NORMAL_ID

	var actionCreateUser user_data_layer.CreateUserExecute
	actionCreateUser.Data = &userData

	var userCreateNew *models.User // khai báo biến lưu trữ record khi tạo mới user.
	userCreateNew, err = intance.userData.UserExecute(&actionCreateUser)

	// khởi tạo mảng dữ liệu key value -> userProfile
	var dataUserProfile = make(map[string]interface{})

	dataUserProfile[constants.USER_PROFILE_FULLNAME] = userDataRequest.FullName
	dataUserProfile[constants.USER_PROFILE_BIRTHDAY] = userDataRequest.BirthDay
	dataUserProfile[constants.USER_PROFILE_PHONENUMBER] = userDataRequest.PhoneNumber

	var profileKeys = []string{constants.USER_PROFILE_FULLNAME, constants.USER_PROFILE_BIRTHDAY, constants.USER_PROFILE_PHONENUMBER}

	// insert thông tin vào user profile với các field còn lại
	for _, key := range profileKeys { // 3 số biến object cần lưu vào user profile (FullName, BirthDay, PhoneNumber)
		var userProfileData models.UserProfile // Khai báo biến để chứa thông tin detail user hợp lệ
		userProfileData.ProfileKey = key
		userProfileData.ProfileValue = dataUserProfile[key].(string)
		userProfileData.UserID = userCreateNew.UserID // Gán UserID khoá ngoại trong UserProfile

		// Khai báo struct CreateUserProfileExecute trong UserProfileExecute
		var actionCreateUserProfile user_data_layer.CreateUserProfileExecute
		actionCreateUserProfile.Data = &userProfileData // set Data với tham trị userProfileData

		_, err = intance.userData.UserProfileExecute(&actionCreateUserProfile)
		if err != nil {
			return err
		}
	}

	// trả về lỗi (nếu có)
	return err
}

func (intance *IntanceUserDataLayer) CreateUserLoginFacebook(userDataRequest DTO.Callback_SignIn_Facebook) error {
	var err error            // Khai báo biến để chứa lỗi trong quá trình thực thi
	var userData models.User // Khai báo biến để chứa thông tin người dùng ghi được

	// set userData từ userDataRequest
	userData.Email = userDataRequest.ID + "@facebook.com"
	// thêm một số trường với rule khi tạo mới tài khoản
	userData.UserStatus = constants.USER_STATUS_PENDING
	userData.LoginMethodID = constants.LOGIN_GOOGLE_ID

	var actionCreateUser user_data_layer.CreateUserExecute
	actionCreateUser.Data = &userData

	var userCreateNew *models.User // khai báo biến lưu trữ record khi tạo mới user.
	userCreateNew, err = intance.userData.UserExecute(&actionCreateUser)

	// khởi tạo mảng dữ liệu key value -> userProfile
	var dataUserProfile = make(map[string]interface{})

	dataUserProfile[constants.USER_PROFILE_FULLNAME] = userDataRequest.Name

	var profileKeys = []string{constants.USER_PROFILE_FULLNAME}

	// insert thông tin vào user profile với các field còn lại
	for _, key := range profileKeys { // 1 số biến object cần lưu vào user profile (FullName)
		var userProfileData models.UserProfile // Khai báo biến để chứa thông tin detail user hợp lệ
		userProfileData.ProfileKey = key
		userProfileData.ProfileValue = dataUserProfile[key].(string)
		userProfileData.UserID = userCreateNew.UserID // Gán UserID khoá ngoại trong UserProfile

		// Khai báo struct CreateUserProfileExecute trong UserProfileExecute
		var actionCreateUserProfile user_data_layer.CreateUserProfileExecute
		actionCreateUserProfile.Data = &userProfileData // set Data với tham trị userProfileData

		_, err = intance.userData.UserProfileExecute(&actionCreateUserProfile)
		if err != nil {
			return err
		}
	}

	// insert thông tin vào auth third party
	var authThirdParty models.AuthThirdParty

	authThirdParty.AccessToken = userDataRequest.AccessToken
	authThirdParty.ProviderID = userDataRequest.ID
	authThirdParty.Provider = "Facebook"
	authThirdParty.ExpiredTime = userDataRequest.Expiry
	authThirdParty.UserID = userCreateNew.UserID

	// Khai báo struct actionCreateAuthThirdParty trong CreateAuthThirdParty
	var actionCreateAuthThirdParty user_data_layer.CreateAuthThirdPartyExecute
	actionCreateAuthThirdParty.Data = &authThirdParty // set Data với tham trị authThirdParty

	_, err = intance.userData.AuthThirdPartyExecute(&actionCreateAuthThirdParty)
	if err != nil {
		return err
	}

	// trả về lỗi (nếu có)
	return err
}
