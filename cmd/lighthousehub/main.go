package main

import (
	"log"

	"github.com/lighthouse-p2p/hub/internal/config"
	"github.com/lighthouse-p2p/hub/internal/database"
	"github.com/lighthouse-p2p/hub/internal/web"
)

func main() {
	cfg := &config.Config{}
	cfg.LoadConfig()

	db, err := database.Connect(cfg.PostgresConfig.GormDSN)
	if err != nil {
		log.Fatalln("Unable to connect to the database")
		log.Fatalf("%s\n", err)
	}
	cfg.Database = db

	web.InitHTTP(cfg)
}
