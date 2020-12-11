package database

import (
	"log"

	"github.com/lighthouse-p2p/hub/internal/models"
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

	db.AutoMigrate(&models.DHT{}, &models.CoinChain{})
	log.Println("Connected to the database!")

	return db, nil
}
