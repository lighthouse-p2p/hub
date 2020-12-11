package main

import (
	"log"

	"github.com/lighthouse-p2p/hub/internal/config"
	"github.com/lighthouse-p2p/hub/internal/database"
	"github.com/lighthouse-p2p/hub/internal/redis"
	"github.com/lighthouse-p2p/hub/internal/web"
)

func main() {
	cfg := &config.Config{}
	cfg.LoadConfig()

	db, err := database.Connect(cfg.PostgresConfig.GormDSN)
	if err != nil {
		log.Println("Unable to connect to the database")
		log.Fatalf("%s\n", err)
	}
	cfg.Database = db

	redisPool, redisConn := redis.Connect(
		cfg.RedisConfig.Host,
		cfg.RedisConfig.Port,
		cfg.RedisConfig.User,
		cfg.RedisConfig.Password,
	)
	cfg.Redis.Pool = redisPool
	cfg.Redis.Conn = redisConn

	defer redis.Close(redisPool, redisConn)

	web.InitHTTP(cfg)
}
