package models

import (
	"auth-service/constants"
	"time"

	"github.com/google/uuid"
)

type UserStatus string

var (
	Pending UserStatus = constants.USER_STATUS_PENDING
	Active  UserStatus = constants.USER_STATUS_ACTIVE
	Hidden  UserStatus = constants.USER_STATUS_HIDDEN
	Blocked UserStatus = constants.USER_STATUS_BLOCKED
)

type User struct {
	UserID         uuid.UUID `gorm:"primaryKey;type:uuid;default:uuid_generate_v4()"`
	Email          string    `gorm:"uniqueIndex;type:varchar(100)"`
	Password       string    `gorm:"type:varchar(250)"`
	UserStatus     UserStatus
	CreatedAt      time.Time `gorm:"default:current_timestamp"`
	UpdatedAt      time.Time
	LoginMethodID  uint
	LoginMethod    LoginMethod
	Profiles       []UserProfile
	AuthThirdParty []AuthThirdParty
	Sessions       []UserSession
}

type LoginMethod struct {
	LoginMethodID uint   `gorm:"primaryKey;autoIncrement"`
	LoginMethod   string `gorm:"type:varchar(30)"`
	Users         []User
}

type UserProfile struct {
	ProfileID    uint      `gorm:"primaryKey;autoIncrement"`
	UserID       uuid.UUID `gorm:"type:uuid"`
	ProfileKey   string    `gorm:"type:text"`
	ProfileValue string    `gorm:"type:text"`
}

type AuthThirdParty struct {
	AuthID       uint      `gorm:"primaryKey;autoIncrement"`
	UserID       uuid.UUID `gorm:"type:uuid"`
	Provider     string    `gorm:"type:varchar(100)"`
	ProviderID   string    `gorm:"uniqueIndex;type:char(30)"`
	AccessToken  string    `gorm:"type:text"`
	RefreshToken string    `gorm:"type:text"`
	ExpiredTime  time.Time
	AuthDetail   string `gorm:"type:text"`
}

type UserSession struct {
	UserSessionID uint      `gorm:"primaryKey;autoIncrement"`
	UserID        uuid.UUID `gorm:"type:uuid"`
	DeviceInfo    string    `gorm:"type:varchar(100)"`
	IPAddress     string    `gorm:"type:varchar(20)"`
	LoginTime     time.Time
	LogoutTime    time.Time
	RefreshToken  string `gorm:"type:varchar(255)"`
}
