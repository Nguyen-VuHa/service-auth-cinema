package domains

import "time"

type GetDetailUserUsecase interface {
	GetDetailUserOnRedis(user_id string) (DetailUserData, error)
	GetDetailUserOnDatabase(user_id string) (DetailUserData, error)
}

type GetDetailUserResponse struct {
	ResponseBasic
	Data DetailUserData `json:"data"`
}

type DetailUserData struct {
	UserID      string    `json:"user_id"`
	Email       string    `json:"email"`
	CreatedAt   time.Time `json:"created_at"`
	PhoneNumber string    `json:"phone_number"`
	FullName    string    `json:"full_name"`
	BirthDay    string    `json:"birth_day"`
}
