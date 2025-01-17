package repository

import (
    "context"
    "github.com/go-redis/redis/v8"
    "fmt"
    "github.com/DamikaAlwis-Gif/shorten-url-app/custom_errors"
		"time"
    "strconv"
)

type RedisRepository struct{
	client *redis.Client
}


func NewRedisRepository(client *redis.Client) *RedisRepository {
  return &RedisRepository{client: client}
}

func (repo *RedisRepository) FindOrignalURLByShortCode(ctx context.Context , shortCode string) (string, error){
	key := fmt.Sprintf("short:%s", shortCode)
  originalURL, err := repo.client.Get(ctx, key).Result()
	if err != nil {
      if err == redis.Nil {
        return "", custom_errors.ErrShortURLNotFound
      }
      return "", fmt.Errorf("error fetching original URL from cache: %w", err)
    }
  return originalURL, nil

}

func (repo *RedisRepository) SaveOriginalURL(ctx context.Context, shortCode string, originalURL string, expiry time.Duration) error{
	key := fmt.Sprintf("short:%s", shortCode)
  err := repo.client.Set(ctx, key, originalURL, expiry).Err()
  if err != nil {
        return fmt.Errorf("error caching URL in Redis: %w", err)
  }
  return nil
}

func (repo *RedisRepository) SetQuota(ctx context.Context, key string, quota int, expiry time.Duration) error {
    err := repo.client.Set(ctx, key, quota, expiry).Err()
    if err != nil {
        return fmt.Errorf("error setting quota in Redis: %w", err)
    }
    return nil
}

func (repo *RedisRepository) GetQuota(ctx context.Context, key string) (int, error) {
    val, err := repo.client.Get(ctx, key).Result()
    if err != nil {
        if err == redis.Nil {
            return 0, custom_errors.ErrKeyNotFound
        }
        return 0, fmt.Errorf("error fetching quota from Redis: %w", err)
    }

    quota, err := strconv.Atoi(val)
    if err != nil {
        return 0, fmt.Errorf("error converting quota value to integer: %w", err)
    }

    return quota, nil
}

func (repo *RedisRepository) DecrementQuota(ctx context.Context, key string) (int, error) {
    val, err := repo.client.Decr(ctx, key).Result()
    if err != nil {
        return 0, fmt.Errorf("error decrementing quota in Redis: %w", err)
    }
    return int(val), nil
}

func (repo *RedisRepository) GetTTL(ctx context.Context, key string) (time.Duration, error) {
    ttl, err := repo.client.TTL(ctx, key).Result()
    if err != nil {
        if err == redis.Nil {
            return 0, custom_errors.ErrKeyNotFound
        }
        return 0, fmt.Errorf("error retrieving TTL for key %s: %w", key, err)
    }
    return ttl, nil
}
