package domains

import (
	"auth-service/models"
	"time"
)

type UserRepository interface {
	Create(user *models.User) error
	Update(user *models.User) error

	GetByEmail(email string) (models.User, error)
	GetByID(id string) (models.User, error)

	GetByEmailPreload(email string, preloads ...interface{}) (models.User, error)
	GetByIDPreload(user_id string, preloads ...interface{}) (models.User, error)
}

type UserProfileRepository interface {
	Create(userProfile *models.UserProfile) error
}

type UserDTO struct {
	UserID        string    `json:"user_id"`
	Email         string    `json:"email"`
	UserStatus    string    `json:"user_status"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
	LoginMethodID uint      `json:"login_method_id"`
	LoginMethod   string    `json:"login_method"`
	PhoneNumber   string    `json:"phone_number"`
	FullName      string    `json:"full_name"`
	BirthDay      string    `json:"birth_day"`
}
