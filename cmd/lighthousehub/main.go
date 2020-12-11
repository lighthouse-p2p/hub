package main

import (
	"github.com/lighthouse-p2p/hub/internal/config"
	"github.com/lighthouse-p2p/hub/internal/database"
)

func main() {
	cfg := &config.Config{}
	cfg.LoadConfig()
	config.SetConfig(cfg)

	database.Connect(cfg.PostgresConfig.GormDSN)
}
