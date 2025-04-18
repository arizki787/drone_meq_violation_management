package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func SetUpDatabaseConnection() *gorm.DB {
	err := godotenv.Load(".env")
	if err != nil {
		log.Println("No .env file found, continuing with environment variables")
	}

	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASS")
	dbHost := os.Getenv("DB_HOST")
	dbName := os.Getenv("DB_NAME")
	dbPort := os.Getenv("DB_PORT")

	// Ensure that all required environment variables are set
	if dbUser == "" || dbPass == "" || dbHost == "" || dbName == "" || dbPort == "" {
		log.Fatal("Database configuration is missing one or more required environment variables")
	}

	// Build the DSN (Data Source Name) for PostgreSQL
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Jakarta",
		dbHost, dbUser, dbPass, dbName, dbPort,
	)

	// Open a connection to the PostgreSQL database
	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN:                  dsn,
		PreferSimpleProtocol: true,
	}), &gorm.Config{})
	if err != nil {
		log.Fatalf("Error opening database connection: %v", err)
	}

	log.Println("Database connection established successfully")
	return db
}

func CloseDatabaseConnection(db *gorm.DB) {
	dbSQL, err := db.DB()
	if err != nil {
		log.Fatalf("Error getting database connection: %v", err)
	}
	err = dbSQL.Close()
	if err != nil {
		log.Fatalf("Error closing database connection: %v", err)
	}
	log.Println("Database connection closed successfully")
}