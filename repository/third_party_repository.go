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
