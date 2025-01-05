package database

import (
	"context"
	"log"
	"os"
	"github.com/go-redis/redis/v8"
)

var Ctx = context.Background()
var client *redis.Client

// create a new client if not already exists
func InitRedis() *redis.Client{
	if client == nil {
		client = redis.NewClient( &redis.Options{
		Addr:     os.Getenv("REDIS_ADDR"),
    Password: os.Getenv("REDIS_PASSWORD"),
    DB:       0,
		PoolSize: 10, // number of connections in the pool
		})
	

	if _, err := client.Ping(Ctx).Result() ; err != nil {
		log.Fatalf("Failed to connect to redis %v.", err.Error() )
	}

	log.Print("Connected to redis")
	
	}
	return client
}


func GetRedisClient() *redis.Client {
	return client
}


