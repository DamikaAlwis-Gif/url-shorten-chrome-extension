package repository

import(
	"context"
)

type MessageBroker interface{
	Publish(ctx context.Context, streamName string, message interface{}) error
	// Subscribe(ctx context.Context, streamName string, handler func(message interface{}) error) error
}