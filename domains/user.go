package domains

import (
	"auth-service/models"
)

type UserRepository interface {
	Create(user *models.User) error
	GetByEmail(email string) (models.User, error)
	GetByID(id string) (models.User, error)
}
