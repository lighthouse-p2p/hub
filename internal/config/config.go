package config

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/gomodule/redigo/redis"
	"github.com/joho/godotenv"
	"gorm.io/gorm"
)

// Config holds the structure for the application config
type Config struct {
	PostgresConfig struct {
		DB       string
		User     string
		Password string
		Host     string
		Port     int
		TLS      bool
		Timezone string
		GormDSN  string
	}

	RedisConfig struct {
		Host     string
		Port     int
		User     string
		Password string
	}

	HTTPConfig struct {
		HTTPAddr      string
		SignalingAddr string
	}

	// ChainConfig is hardncoded for now
	ChainConfig struct {
		// Worth is the value of coins per mb
		Worth float64

		// InitialCoins is the number of coins initially given to the user while registration
		InitialCoins float64
	}

	// Database holds a gorm.DB instance for access with the handler
	Database *gorm.DB

	// Redis holds the global redis pool and connection
	Redis struct {
		Pool *redis.Pool
		Conn redis.Conn
	}
}

// LoadConfig loads the config from the environment file
func (c *Config) LoadConfig() {
	log.Println("Loading the config from .env")

	err := godotenv.Load()
	if err != nil {
		log.Fatalln("Error loading .env file")
	}

	pgDB, pgDBPresent := os.LookupEnv("POSTGRES_DB")
	pgUser, pgUserPresent := os.LookupEnv("POSTGRES_USER")
	pgPass, pgPassPresent := os.LookupEnv("POSTGRES_PASSWORD")
	pgHost, pgHostPresent := os.LookupEnv("POSTGRES_HOST")
	pgPort, pgPortPresent := os.LookupEnv("POSTGRES_PORT")
	pgTLS, pgTLSPresent := os.LookupEnv("POSTGRES_TLS")
	pgTimezone, pgTimezonePresent := os.LookupEnv("POSTGRES_TIMEZONE")

	redisHost, redisHostPresent := os.LookupEnv("REDIS_HOST")
	redisPort, redisPortPresent := os.LookupEnv("REDIS_PORT")
	redisUser, redisUserPresent := os.LookupEnv("REDIS_USER")
	redisPassword, redisPasswordPresent := os.LookupEnv("REDIS_PASSWORD")

	httpAddr, httpAddrPresent := os.LookupEnv("HTTP_ADDR")
	signalingAddr, signalingAddrPresent := os.LookupEnv("SIGNALING_ADDR")

	if !pgDBPresent ||
		!pgUserPresent ||
		!pgPassPresent ||
		!pgHostPresent ||
		!pgPortPresent ||
		!pgTLSPresent ||
		!pgTimezonePresent ||
		!httpAddrPresent ||
		!signalingAddrPresent ||
		!redisHostPresent ||
		!redisPortPresent ||
		!redisUserPresent ||
		!redisPasswordPresent {
		log.Fatalln(".env is not complete :)")
	}

	c.PostgresConfig.DB = pgDB
	c.PostgresConfig.User = pgUser
	c.PostgresConfig.Password = pgPass
	c.PostgresConfig.Host = pgHost
	c.PostgresConfig.Port, err = strconv.Atoi(pgPort)
	if err != nil {
		log.Fatalln("POSTGRES_PORT is not an integer")
	}
	c.PostgresConfig.TLS, err = strconv.ParseBool(pgTLS)
	if err != nil {
		log.Fatalln("POSTGRES_TLS is not a bool")
	}
	c.PostgresConfig.Timezone = pgTimezone

	sslMode := "disable"
	if c.PostgresConfig.TLS {
		sslMode = "verify-ca"
	}

	c.PostgresConfig.GormDSN = fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%d sslmode=%s TimeZone=%s",
		pgHost,
		pgUser,
		pgPass,
		pgDB,
		c.PostgresConfig.Port,
		sslMode,
		pgTimezone,
	)

	c.RedisConfig.Host = redisHost
	c.RedisConfig.Port, err = strconv.Atoi(redisPort)
	if err != nil {
		log.Fatalln("REDIS_PORT is not an integer")
	}
	c.RedisConfig.User = redisUser
	c.RedisConfig.Password = redisPassword

	c.HTTPConfig.HTTPAddr = httpAddr
	c.HTTPConfig.SignalingAddr = signalingAddr

	c.ChainConfig.InitialCoins = 100
	c.ChainConfig.Worth = 1

	log.Println("Config loaded successfully")
}
