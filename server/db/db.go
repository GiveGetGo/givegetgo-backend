package db

import (
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// ConnectPostgresDB initializes and returns a GORM database connection
func ConnectPostgresDB(dburl string) (*gorm.DB, error) {
	// Open the connection
	db, err := gorm.Open(postgres.Open(dburl), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	log.Println("Successfully connected to PostgreSQL")
	return db, nil
}
