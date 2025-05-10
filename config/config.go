package config

import (
	"fmt"
	"os"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/joho/godotenv"
)

var DB *gorm.DB

// InitPostgres initializes the database connection with GORM and pool configurations
func InitPostgres() error {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		return fmt.Errorf("error loading .env file: %w", err)
	}

	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		return fmt.Errorf("DATABASE_URL is not set")
	}

	var err error
	DB, err = gorm.Open("postgres", dsn)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	// Manually configure the pool (optional, to have full control)
	sqlDB := DB.DB()
	sqlDB.SetMaxIdleConns(5)              // Set the max idle connections
	sqlDB.SetMaxOpenConns(50)             // Set the max open connections
	sqlDB.SetConnMaxLifetime(time.Hour)   // Set connection max lifetime
	sqlDB.SetConnMaxIdleTime(time.Minute) // Set max idle connection time

	return nil
}

// Close gracefully closes the DB connection
func Close() {
	if DB != nil {
		err := DB.Close()
		if err != nil {
			fmt.Printf("Error closing the database: %v\n", err)
		}
	}
}
