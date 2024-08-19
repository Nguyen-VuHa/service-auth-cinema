package repository

import (
	"auth-service/domains"
	"auth-service/models"

	"gorm.io/gorm"
)

type userProfileRepository struct {
	db *gorm.DB
}

func NewUserProfileRepository(db *gorm.DB) domains.UserProfileRepository {
	return &userProfileRepository{
		db: db,
	}
}

func (u *userProfileRepository) Create(userProfile *models.UserProfile) error {
	createAction := u.db.Create(&userProfile)

	if createAction.Error != nil {
		return createAction.Error
	}

	return nil
}
