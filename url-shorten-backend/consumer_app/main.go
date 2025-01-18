package main

import (
    "context"
    "encoding/json"
    "fmt"
    "os"
    "time"
    "github.com/go-redis/redis/v8"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
    "github.com/joho/godotenv"
)

type ClickData struct {
    ShortURL  string    `json:"short_url" bson:"short_url"`
    Timestamp time.Time `json:"timestamp" bson:"timestamp"`
    IPAddress string    `json:"ip_address" bson:"ip_address"`
    UserAgent string    `json:"user_agent" bson:"user_agent"`
}

func main() {
    // Load environment variables from .env file
    err := godotenv.Load()
    if err != nil {
        fmt.Println("Error loading .env file")
    }

    // Connect to Redis
    rdb := redis.NewClient(&redis.Options{
        Addr:     os.Getenv("REDIS_DB_ADDR"),
        Password: os.Getenv("REDIS_DB_PASS"),
        DB:       0,
        Username: os.Getenv("REDIS_DB_USERNAME"),
    })
    if err := rdb.Ping(context.Background()).Err(); err != nil {
        panic(fmt.Sprintf("Failed to connect to Redis: %v", err))
    }

    // Connect to MongoDB
    mongoClient, err := mongo.Connect(context.Background(), options.Client().ApplyURI(os.Getenv("MONGODB_URI")))
    if err != nil {
        panic(fmt.Sprintf("Failed to connect to MongoDB: %v", err))
    }
    db := mongoClient.Database(os.Getenv("MONGODB_DATABASE_NAME"))
    collection := db.Collection(os.Getenv("MONGODB_COLLECTION_NAME_CLICKS"))

    // Create Consumer Group if it doesn't exist (with ">" to start from the latest messages)
    err = rdb.XGroupCreateMkStream(context.Background(), "click_logs", "my_consumer_group", "0").Err()
    if err != nil && err.Error() != "BUSYGROUP Consumer Group name already exists" {
        fmt.Printf("Error creating consumer group: %v\n", err)
        return
    }

    // Start consuming messages as part of the consumer group
    for {
        // XReadGroup will block and read new messages from the stream for this consumer
        msgs, err := rdb.XReadGroup(context.Background(), &redis.XReadGroupArgs{
            Group:    "my_consumer_group", // Group name
            Consumer: "consumer1",         // Consumer name (each consumer in the group should have a unique name)
            Streams:  []string{"click_logs", ">"},
            Count:    1,   // Number of messages to read per request
            Block:    0,    // Block indefinitely until new messages arrive
        }).Result()

        if err != nil {
            fmt.Printf("Error reading from group: %v\n", err)
            continue
        }

        fmt.Printf("Received new click data: %+v\n", msgs)

        // Process the messages
        for _, stream := range msgs {
            for _, msg := range stream.Messages {
                if data, ok := msg.Values["data"].(string); ok {
                    // Parse the message
                    var clickData ClickData
                    if err := json.Unmarshal([]byte(data), &clickData); err != nil {
                        fmt.Printf("Failed to parse message: %v\n", err)
                        continue
                    }

                    // Save to MongoDB
                    _, err := collection.InsertOne(context.Background(), clickData)
                    if err != nil {
                        fmt.Printf("Failed to save to MongoDB: %v\n", err)
                        continue
                    }

                    fmt.Printf("Saved click data to MongoDB: %+v\n", clickData)
                }
            }
        }
    }
}
