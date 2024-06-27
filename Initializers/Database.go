package initializers

import (
	"os"
	models "service-auth/Models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDatabase() {
	dsn := os.Getenv("DNS_DATABASE")
	var err error

	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		SkipDefaultTransaction: true,
	})

	if err != nil {
		panic("Failed to connect to db")
	}
}

func MigrateDatabase() {
	// migrate user
	DB.AutoMigrate(&models.User{}, &models.LoginMethod{}, &models.UserProfile{}, &models.UserSession{}, &models.AuthThirdParty{})
}
