package database

import (
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Connect is used to connect to the gorm db
func Connect(dsn string) (*gorm.DB, error) {
	log.Println("Connecting to the database...")

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	log.Println("Connected to the database!")

	return db, nil
}
