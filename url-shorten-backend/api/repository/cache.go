package repository

import (
	"context"
	"time"
)

type CacheRepository interface{
	FindOrignalURLByShortCode(ctx context.Context , shortCode string) (string, error)
	SaveOriginalURL(ctx context.Context, shortCode string, originalURL string, expiry time.Duration) error
	SetQuota(ctx context.Context, key string, quota int, expiry time.Duration) error
	GetQuota(ctx context.Context, key string) (int, error)
	DecrementQuota(ctx context.Context, key string) (int, error)
	GetTTL(ctx context.Context, key string) (time.Duration, error)
}
