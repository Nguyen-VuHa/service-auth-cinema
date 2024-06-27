package user_data_layer

import (
	models "service-auth/Models"

	"gorm.io/gorm"
)

// *Áp dụng Design Pattern Strategy cho method Execute
type AuthThirdPartyExecuteQuery interface {
	ExecuteAuthThirdParty(DB *gorm.DB) (*models.AuthThirdParty, error)
}

// create strategy user execute
func (data *UserDataLayer) AuthThirdPartyExecute(execute AuthThirdPartyExecuteQuery) (*models.AuthThirdParty, error) {
	return execute.ExecuteAuthThirdParty(data.DB)
}

// Khởi tạo struct CreateAuthThirdPartyExecute
type CreateAuthThirdPartyExecute struct {
	Data *models.AuthThirdParty
}

// thực thi tạo user với hàm ExecuteUserProfile
func (create *CreateAuthThirdPartyExecute) ExecuteAuthThirdParty(DB *gorm.DB) (*models.AuthThirdParty, error) {
	errCreate := DB.Create(&create.Data).Error

	return create.Data, errCreate
}
