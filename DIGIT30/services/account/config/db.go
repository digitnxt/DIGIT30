package config

import (
	"account/model"
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"os"
	"strconv"
)

var database *gorm.DB

// DatabaseInit initializes the GORM DB connection using values from environment variables.
func DatabaseInit() {
	// Retrieve configuration from environment variables.
	host := os.Getenv("DB_HOST")
	if host == "" {
		host = "localhost"
	}

	user := os.Getenv("DB_USER")
	if user == "" {
		user = "postgres"
	}

	password := os.Getenv("DB_PASSWORD")
	if password == "" {
		password = "postgres"
	}

	dbName := os.Getenv("DB_NAME")
	if dbName == "" {
		dbName = "account"
	}

	portStr := os.Getenv("DB_PORT")
	port := 5432
	if portStr != "" {
		if p, err := strconv.Atoi(portStr); err == nil {
			port = p
		}
	}

	// Construct DSN (Data Source Name)
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable TimeZone=Asia/Jakarta",
		host, user, password, dbName, port)

	var err error
	database, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(fmt.Sprintf("failed to connect to database: %v", err))
	}

	// Create the UUID extension if it doesn't exist.
	if err := database.Exec(`CREATE EXTENSION IF NOT EXISTS "uuid-ossp";`).Error; err != nil {
		panic(fmt.Sprintf("failed to create uuid-ossp extension: %v", err))
	}

	// Automatically migrate your Account model (this will create/update your table schema)
	if err := database.AutoMigrate(&model.Account{}); err != nil {
		panic(fmt.Sprintf("failed to auto-migrate: %v", err))
	}
}

// DB returns the current GORM DB instance.
func DB() *gorm.DB {
	return database
}
