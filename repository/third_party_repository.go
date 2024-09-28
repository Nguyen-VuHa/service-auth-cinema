package repository

import (
	"auth-service/domains"
	"auth-service/models"

	"gorm.io/gorm"
)

type thirdPartyRepository struct {
	db *gorm.DB
}

func NewThirdPartyRepository(db *gorm.DB) domains.ThirdPartyRepository {
	return &thirdPartyRepository{
		db: db,
	}
}

func (u *thirdPartyRepository) Create(user *models.AuthThirdParty) error {
	createAction := u.db.Create(&user)

	if createAction.Error != nil {
		return createAction.Error
	}

	return nil
}

func (u *thirdPartyRepository) GetByUserID(user_id string) (models.AuthThirdParty, error) {
	var third_party_token models.AuthThirdParty

	err := u.db.Model(&models.AuthThirdParty{}).Where("user_id = ?", user_id).First(&third_party_token).Error

	return third_party_token, err
}
