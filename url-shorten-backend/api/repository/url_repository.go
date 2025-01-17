package repository

// import (
// 	"context"
// 	"errors"
// 	"fmt"
// 	"time"
// 	"github.com/DamikaAlwis-Gif/shorten-url-app/custom_errors"
// 	"github.com/DamikaAlwis-Gif/shorten-url-app/database"
// 	"github.com/DamikaAlwis-Gif/shorten-url-app/helpers"
// 	"github.com/DamikaAlwis-Gif/shorten-url-app/service"
// 	"github.com/DamikaAlwis-Gif/shorten-url-app/config"
// 	"github.com/go-redis/redis/v8"
// 	"go.mongodb.org/mongo-driver/bson"
// 	"go.mongodb.org/mongo-driver/mongo"
// 	"go.mongodb.org/mongo-driver/mongo/options"
// )

// func fetchOriginalURLPersistantDB(ctx context.Context, srv *service.Service, shortURL string) (string, error) {
// 	// get mongo client
// 	mdb := srv.MongoDB
// 	collection, err := mdb.GetCollection(config.AppConfig.MongoDBDatabaseName, config.AppConfig.MongoDBCollectionNameUrls)
// 	if err!= nil {
//     return "", err
//   }

// 	// filter retrive short url
// 	filter := bson.M{
// 		"short_url": shortURL,
// 	}
// 	// define the projection to return only the original_url field
// 	projection := bson.M{
// 		"original_url": 1,
// 		"_id"         : 0,
// 	}
// 	// Define a result structure to fetch the OriginalURL field
//   var result struct {
//         OriginalURL string `bson:"original_url"`
//   }
//   err = collection.FindOne(ctx, filter, options.FindOne().SetProjection(projection)).Decode(&result)
// 	if err!= nil {
// 		if errors.Is(err, mongo.ErrNoDocuments){
// 			return "", custom_errors.ErrShortURLNotFound
// 		}
// 		return "", fmt.Errorf("error finding original URL: %w", err)
// 	}
  
//   return result.OriginalURL, nil
// }

// func fetchOriginalURLCache(ctx context.Context, srv *service.Service, shortURL string)(string, error){
//   rdb := srv.Redis.GetDBClient()
//   // generate short key to query redis eg : short:abc123
//   short_key := fmt.Sprintf("short:%s", shortURL)

//   // get the original url from redis
//   originalURL, err := rdb.Get(ctx, short_key).Result()
//   if  err != nil {
// 		// if redis doesn't have the short key
//     if errors.Is(err, redis.Nil){
//       return "", custom_errors.ErrShortURLNotFound
//     }
// 		// for other errors
//     return "", fmt.Errorf("error fetching original URL from cache: %w", err)
    
//   }
// 	return originalURL, nil

// }


// // cacheOriginalURL caches the original URL in Redis.
// func cacheOriginalURL(ctx context.Context,srv *service.Service, shortURL string, originalURL string) error {
// 	rdb := srv.Redis.GetDBClient()
	
// 	// Generate the Redis key
// 	shortKey := fmt.Sprintf("short:%s", shortURL)

// 	// Cache the URL with a TTL (e.g., 24 hours)
// 	err := rdb.Set(ctx, shortKey, originalURL, 24 * time.Hour).Err() // 24 hours
// 	if err != nil {
// 		return fmt.Errorf("error caching URL in Redis: %w", err)
// 	}

// 	return nil
// }


// func ResolveShortURL(ctx context.Context, srv *service.Service, shortURL string) (string, error) {
//     originalURL, err := fetchOriginalURLCache(ctx,srv, shortURL)
//     if err == nil {
//         return originalURL, nil
//     }
//     if !errors.Is(err, custom_errors.ErrShortURLNotFound) {
//         return "", err
//     }
//     originalURL, err = fetchOriginalURLPersistantDB(ctx, srv, shortURL)
//     if err != nil {
//         return "", err
//     }
//     _ = cacheOriginalURL(ctx, srv, shortURL, originalURL) // Ignore caching errors
//     return originalURL, nil
// }


// func setShortURLCache(ctx context.Context, srv *service.Service, shortURL string, originalURL string, expiry time.Duration) error{
// 	rdb := srv.Redis.GetDBClient()
	
// 	if expiry <= 0 {
// 		return fmt.Errorf("expiry time should be greater than current time")
// 	}
	
// 	err := rdb.Set(ctx, shortURL, originalURL, expiry).Err()
// 	if err!= nil {
//     return fmt.Errorf("failed to cache short URL: %w", err)
//   }
// 	return nil

// }

// func setShortURLPersistantDB(ctx context.Context, srv *service.Service, shortURL string, originalURL string, isCustomURL bool, expiry time.Duration) error {
// 	// Get the MongoDB collection
// 	mdb := srv.MongoDB
// 	collection, err := mdb.GetCollection(config.AppConfig.MongoDBDatabaseName, config.AppConfig.MongoDBCollectionNameUrls)
// 	if err != nil {
// 		return fmt.Errorf("failed to get MongoDB collection: %w", err)
// 	}
// 	// check that expiry time is positive and greater than the current time
// 	if expiry <= 0 {
//     return fmt.Errorf("expiry time should be greater than current time")
//   }

// 	// if it's a custom short URL, check if it already exists
// 	if isCustomURL{

// 		// check if the custom short URL already exists in the database
// 		filter := bson.M{"short_url": shortURL}
// 		var existingDoc database.ShortenURL
// 		err := collection.FindOne(ctx, filter).Decode(&existingDoc)
// 		if err == nil {
// 			// If we find an existing custom short URL, return an error
//       return custom_errors.ErrShortKeyExists

//     } else if !errors.Is(err, mongo.ErrNoDocuments) {
// 			// If there's any other error while querying
//       return fmt.Errorf("failed to check if custom short URL already exists: %w", err)
//     }

// 	}

// 	// create a new document with the provided shortURL and originalURL
// 	urlDoc := database.ShortenURL{
// 		ShortURL:      shortURL,
//     OriginalURL:   originalURL,
// 		CreatedAt:  time.Now(),
// 		Expiry: time.Now().Add(expiry),
// 	}

// 	// insert the new document
// 	_, err = collection.InsertOne(ctx, urlDoc)
//   if err!= nil {
//     return fmt.Errorf("failed to insert short URL to MongoDB: %w", err)
//   }
//   return nil

// }


// func SetShortURL(ctx context.Context, srv *service.Service, customShort, originalURL string, expiry time.Duration) (shortURL string, err error) {
// 	// if expiry not specified, set it to 24 hours
// 	if expiry == 0 {
// 		expiry = 24 * time.Hour
// 	}
// 	// custom short provided or not provided
// 	isCustomShort := customShort != ""
// 	// get short code
// 	shortCode, err := helpers.GetShort(customShort)
// 	if err!= nil {
//     return "", err
//   }
// 	// set short URL in database and cache
// 	err = setShortURLPersistantDB(ctx,srv, shortCode, originalURL, isCustomShort, expiry)
// 	if err != nil {
// 		return "", fmt.Errorf("error setting short persistent DB: %w", err)
// 	}
// 	// create short code
// 	redisKey := fmt.Sprintf("short:%s", shortCode)
// 	err = setShortURLCache(ctx, srv, redisKey, originalURL, expiry)
// 	if err != nil {
// 		return "", fmt.Errorf("error caching short URL: %w", err)
// 	}
// 	return shortCode, nil
// }


// type RateLimitResponse struct {
// 	RemainingQuota int
// 	ResetAfter      time.Duration
// }

// func GetRateLimit(ctx context.Context, srv *service.Service, ipAddress string) (*RateLimitResponse, error) {
// 	rdb := srv.Redis.GetDBClient()
  
// 	// decrement the remaining quota and fetch it
// 	remainingQuota, err := rdb.Decr(ctx, ipAddress).Result()
// 	if err != nil{
// 		return nil , fmt.Errorf("failed  to fetch remaining quota: %w", err)
// 	}
// 	// get the quota reset time 
// 	resetAfter, err := rdb.TTL(ctx, ipAddress).Result()
// 	if err!= nil {
//     return nil, fmt.Errorf("failed to get reset after time: %w", err)
//   }

// 	return &RateLimitResponse{RemainingQuota: int(remainingQuota), ResetAfter: resetAfter}, nil

  
// }

