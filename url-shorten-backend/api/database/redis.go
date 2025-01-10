package database

import (
	"context"
	"fmt"
	"log"
	"os"
	"sync"
	"time"
	"github.com/go-redis/redis/v8"
)

// Redis struct to hold the Redis client and sync.Once
type Redis struct {
	Client *redis.Client
	once   sync.Once
}

// Initialize and connect to Redis if not already done
func (r *Redis) InitDB(parentCtx context.Context) error {
	var err error
	r.once.Do(func() {

		ctx, cancel := context.WithTimeout(parentCtx, 10*time.Second)
		defer cancel()

		if os.Getenv("REDIS_ADDR") == "" {
      err = fmt.Errorf("missing REDIS_ADDR in environment variable")
      return
    }
		r.Client = redis.NewClient(&redis.Options{
			Addr:     os.Getenv("REDIS_ADDR"),  // Redis address (e.g., "localhost:6379")
			Password: os.Getenv("REDIS_PASSWORD"), // Redis password if any
			DB:       0,                          // Default DB
			PoolSize: 10,                         // Pool size for Redis connections
		})

		// Test Redis connection by pinging the server
		if err = r.Client.Ping(ctx).Err(); err != nil {
			err = fmt.Errorf("failed to connect to Redis: %w", err)
			return
		}

		log.Print("connected to Redis")
	})

	return err
}

// GetRedisClient returns the Redis client instance
func (r *Redis) GetDBClient(parentCtx context.Context) (*redis.Client, error) {
	if r.Client == nil{
	 log.Println("Redis client not initialized. Initializing now...")
	 if err := r.InitDB(parentCtx); err != nil {
		return nil, err
	 }
	}
	return r.Client, nil
}

// CloseRedisConnection closes the Redis connection
func (r *Redis) CloseDBConnection() error {
	if r.Client != nil {
		// Close the Redis connection
		if err := r.Client.Close(); err != nil {
			log.Printf("Failed to close Redis connection: %v", err)
			return err
		}
		log.Println("Redis connection closed")
	}
	return nil
}
