package repository
import (
    "context"
    "encoding/json"
    "github.com/go-redis/redis/v8"

)

type RedisBroker struct{
	client *redis.Client
}

func NewRedisBroker(client *redis.Client) *RedisBroker {
  return &RedisBroker{client: client}
}

// publish a message to a redis stream
// interface to handle any type of message
func (rb *RedisBroker) Publish(ctx context.Context, streamName string, message interface{}) error {
	  // serialize the message to json
    msgBytes, err := json.Marshal(message)
    if err != nil {
        return err
    }
		// use xadd to publish the message to the Redis stream
		  _, err = rb.client.XAdd(ctx, &redis.XAddArgs{
        Stream: streamName,                  // Name of the stream
        Values: map[string]interface{}{     // Fields of the message
            "data": string(msgBytes),       // Store the JSON-encoded message
        },
    }).Result()

    return err
}



