package repository

import (
	"context"
	"errors"
	"fmt"
	"time"
	"github.com/DamikaAlwis-Gif/shorten-url-app/custom_errors"
	"github.com/DamikaAlwis-Gif/shorten-url-app/database"
	"github.com/go-redis/redis/v8"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func FetchOriginalURLPersistantDB(ctx context.Context, shortURL string) (string, error) {
	// get mongo client
  mdb:= &database.MongoDB{}
	collection, err := mdb.GetCollection("url_shortner","urls")
	if err!= nil {
    return "", err
  }

	// filter retrive short url
	filter := bson.M{
		"short_url": shortURL,
	}
	// define the projection to return only the original_url field
	projection := bson.M{
		"original_url": 1,
		"_id"         : 0,
	}
	// Define a result structure to fetch the OriginalURL field
  var result struct {
        OriginalURL string `bson:"original_url"`
  }
  err = collection.FindOne(ctx, filter, options.FindOne().SetProjection(projection)).Decode(&result)
	if err!= nil {
		if errors.Is(err, mongo.ErrNoDocuments){
			return "", custom_errors.ErrShortURLNotFound
		}
		return "", fmt.Errorf("error finding original URL: %w", err)
	}
  
  return result.OriginalURL, nil
}

func FetchOriginalURLCache(ctx context.Context,shortURL string)(string, error){
	rdb := &database.Redis{}
  redisClient, err := rdb.GetDBClient(ctx)
  if err!= nil {
    return "", err
  }

  // generate short key to query redis eg : short:abc123
  short_key := fmt.Sprintf("short:%s", shortURL)

  // get the original url from redis
  originalURL, err := redisClient.Get(ctx, short_key).Result()
  if  err != nil {
		// if redis doesn't have the short key
    if errors.Is(err, redis.Nil){
      return "", custom_errors.ErrShortURLNotFound
    }
		// for other errors
    return "", fmt.Errorf("error fetching original URL from cache: %w", err)
    
  }
	return originalURL, nil

}


// CacheOriginalURL caches the original URL in Redis.
func CacheOriginalURL(ctx context.Context, shortURL string, originalURL string) error {
	rdb := &database.Redis{}
	redisClient, err := rdb.GetDBClient(ctx)
	if err != nil {
		return err
	}

	// Generate the Redis key
	shortKey := fmt.Sprintf("short:%s", shortURL)

	// Cache the URL with a TTL (e.g., 24 hours)
	err = redisClient.Set(ctx, shortKey, originalURL, 24 * time.Hour).Err() // 24 hours
	if err != nil {
		return fmt.Errorf("error caching URL in Redis: %w", err)
	}

	return nil
}


func ResolveShortURL(ctx context.Context, shortURL string) (string, error) {
    originalURL, err := FetchOriginalURLCache(ctx, shortURL)
    if err == nil {
        return originalURL, nil
    }
    if !errors.Is(err, custom_errors.ErrShortURLNotFound) {
        return "", err
    }
    originalURL, err = FetchOriginalURLPersistantDB(ctx, shortURL)
    if err != nil {
        return "", err
    }
    _ = CacheOriginalURL(ctx, shortURL, originalURL) // Ignore caching errors
    return originalURL, nil
}
