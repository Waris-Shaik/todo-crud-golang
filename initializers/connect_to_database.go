package initializers

import (
	"fmt"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectToDB() {
	var err error

	// Get DB_URL from .env
	dsn := os.Getenv("DB_URL")

	// DB_URL error
	if dsn == "" {
		log.Fatal("DB_URL is not defined in .env file.")
	}

	// database connection
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})

	// DB connection success
	fmt.Println("Database Connected Successfully..🔥🔥🔥")

	// DB connection failure
	if err != nil {
		log.Fatal("Failed to connect to database..⛔⛔⛔")
	}

}
