package database

import (
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var dbConnection *gorm.DB

// GetDBConnection is used to globally access the gorm connection
func GetDBConnection() *gorm.DB {
	return dbConnection
}

// Connect is used to connect to the gorm db
func Connect(dsn string) {
	log.Println("Connecting to the database...")

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalln("Unable to connect to the database")
		log.Fatalf("%s\n", err)
	}

	dbConnection = db

	log.Println("Connected to the database!")
}
