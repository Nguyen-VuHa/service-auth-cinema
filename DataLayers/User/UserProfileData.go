package user_data_layer

import (
	models "service-auth/Models"

	"gorm.io/gorm"
)

// *Áp dụng Design Pattern Strategy cho method Execute
type UserProfileExecuteQuery interface {
	ExecuteUserProfile(DB *gorm.DB) (*models.UserProfile, error)
}

// create strategy user execute
func (data *UserDataLayer) UserProfileExecute(execute UserProfileExecuteQuery) (*models.UserProfile, error) {
	return execute.ExecuteUserProfile(data.DB)
}

type CreateUserProfileExecute struct {
	Data *models.UserProfile
}

// thực thi tạo user với hàm ExecuteUserProfile
func (create *CreateUserProfileExecute) ExecuteUserProfile(DB *gorm.DB) (*models.UserProfile, error) {
	errCreate := DB.Create(&create.Data).Error

	return create.Data, errCreate
}
