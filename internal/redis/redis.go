package redis

import (
	"context"
	"fmt"
	"log"

	"github.com/go-redis/redis/v8"
)

// Connect is used to connect to redis and return a pool and a connection
func Connect(host string, port int, username, password string) *redis.Client {
	log.Println("Connecting to redis...")

	// redisPool := &redis.Pool{
	// 	Dial: func() (redis.Conn, error) {
	// 		return redis.Dial("tcp", fmt.Sprintf("%s:%d", host, port), redis.DialUsername(username), redis.DialPassword(password))
	// 	},
	// }

	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", host, port),
		Username: username,
		Password: password,
	})

	_, err := rdb.Ping(context.Background()).Result()
	if err != nil {
		log.Println("Redis PING failed, exitting...")
		log.Fatalf("%s\n", err)
	}

	log.Println("Connected to redis")

	return rdb
}

// Close cleans up the pool and the connection
func Close(rdb *redis.Client) {
	rdb.Close()
}
