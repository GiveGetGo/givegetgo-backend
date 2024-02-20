package db

import (
	"log"
	"server/schema"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Database interface {
	AutoMigrate(models ...interface{}) error
	Create(value interface{}) *gorm.DB
	Where(query interface{}, args ...interface{}) *gorm.DB
	First(dest interface{}, conds ...interface{}) *gorm.DB
}

// Ensure that *gorm.DB satisfies the Database interface
var _ Database = (*gorm.DB)(nil)

// ConnectPostgresDB initializes and returns a GORM database connection
func ConnectPostgresDB(dburl string) *gorm.DB {
	// Open the connection
	db, err := gorm.Open(postgres.Open(dburl), &gorm.Config{})
	if err != nil {
		log.Fatalf("Error connecting to PostgreSQL: %v", err)
		return nil
	}

	log.Println("Successfully connected to PostgreSQL")
	return db
}

// AutoMigratePostgresDB migrates the database schema
func AutoMigratePostgresDB(db *gorm.DB) error {
	// Migrate the schema
	err := db.AutoMigrate(&schema.User{}, &schema.RegisterEmailVerification{})
	if err != nil {
		log.Fatalf("Error migrating PostgreSQL schema: %v", err)
		return err
	}

	log.Println("Successfully migrated PostgreSQL schema")
	return nil
}
