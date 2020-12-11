package redis

import (
	"fmt"
	"log"

	"github.com/gomodule/redigo/redis"
)

// Connect is used to connect to redis and return a pool and a connection
func Connect(host string, port int, username, password string) (*redis.Pool, redis.Conn) {
	log.Println("Connecting to redis...")

	redisPool := &redis.Pool{
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", fmt.Sprintf("%s:%d", host, port), redis.DialUsername(username), redis.DialPassword(password))
		},
	}

	conn := redisPool.Get()

	_, err := conn.Do("PING")
	if err != nil {
		log.Println("Redis PING failed, exitting...")
		log.Fatalf("%s\n", err)
	}

	log.Println("Connected to redis")

	return redisPool, conn
}

// Close cleans up the pool and the connection
func Close(pool *redis.Pool, conn redis.Conn) {
	conn.Close()
	pool.Close()
}
