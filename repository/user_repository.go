package repository

import (
	"auth-service/domains"
	"auth-service/models"

	"gorm.io/gorm"
)

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) domains.UserRepository {
	return &userRepository{
		db: db,
	}
}

func (u *userRepository) Create(user *models.User) error {
	createAction := u.db.Create(&user)

	if createAction.Error != nil {
		return createAction.Error
	}

	return nil
}

func (u *userRepository) Update(user *models.User) error {
	updateAction := u.db.Save(&user)

	if updateAction.Error != nil {
		return updateAction.Error
	}

	return nil
}

func (u *userRepository) GetByEmail(email string) (models.User, error) {
	var user models.User
	err := u.db.Model(&models.User{}).Where("email = ?", email).First(&user).Error

	return user, err
}

func (u *userRepository) GetByID(id string) (models.User, error) {
	var user models.User
	err := u.db.Model(&models.User{}).Where("user_id = ?", id).First(&user).Error

	return user, err
}

func (u *userRepository) GetByEmailPreload(email string, preloads ...interface{}) (models.User, error) {
	var user models.User

	query := u.db.Model(&models.User{}).Where("email = ?", email)

	// Lặp qua các preloads và áp dụng chúng vào query
	for _, preload := range preloads {
		query = query.Preload(preload.(string))
	}

	err := query.First(&user).Error

	return user, err
}
