package service

import (
    "errors"
	"context"
	"log"
	"time"
    "fmt"
	"github.com/DamikaAlwis-Gif/shorten-url-app/custom_errors"
	"github.com/DamikaAlwis-Gif/shorten-url-app/database"
	"github.com/DamikaAlwis-Gif/shorten-url-app/helpers"
	"github.com/DamikaAlwis-Gif/shorten-url-app/repository"
)

// Service holds the dependencies for the application
type Service struct {
    Redis   *database.Redis
    MongoDB *database.MongoDB
}

type URLService struct{
    cacheRepo repository.CacheRepository
    dbRepo repository.DBRepository
}

func NewURLService(cacheRepo repository.CacheRepository, dbRepo repository.DBRepository) *URLService{
    return &URLService{cacheRepo, dbRepo}
}

func (s *URLService) ResolveShortURL(ctx context.Context, shortCode string) (string, error) {
    // Check cache first
    originalURL, err := s.cacheRepo.FindOrignalURLByShortCode(ctx, shortCode)
    if err == nil {
        return originalURL, nil
    }

    // If not found in cache, check database
    urlDoc, err := s.dbRepo.FindOriginalURLDetailsByShortCode(ctx, shortCode)
    if err != nil {
        return "", err
    }

    // Check for URL expiration
    if time.Now().After(urlDoc.Expiry) {
        return "", custom_errors.ErrURLExpired
    }

    // Calculate the remaining TTL
    expiry := time.Until(urlDoc.Expiry)
    if expiry > 24*time.Hour {
        expiry = 24 * time.Hour // Cap the cache TTL to 24 hours
    }

    // Cache the result for future requests
    if err := s.cacheRepo.SaveOriginalURL(ctx, shortCode, urlDoc.OriginalURL, expiry); err != nil {
        // Log the error but continue with the database value
        log.Printf("Error caching URL for shortCode %s: %v\n", shortCode, err)
    }

    // Return the original URL from the database
    return urlDoc.OriginalURL, nil
}


func (s *URLService) CreateShortURL(ctx context.Context, shortCode string, isCustom bool, originalURL string, expiry time.Duration) (string, error){
    
    if (!isCustom) {
        var err error
        retryCount := 3
        for i := 0; i < retryCount ; i++ {
            
            shortCode, err = helpers.GenerateShortURL(originalURL)
            if err!= nil {
                return "", fmt.Errorf("error generating short URL: %w", err)
            }
            err = s.dbRepo.SaveOriginalURL(ctx, shortCode, originalURL, isCustom, expiry)
            if err == nil {
                break
            }
            if !errors.Is(err, custom_errors.ErrShortKeyExists){
                return "", err
            }
            log.Printf("Short code %s already exists, retrying... (%d/%d)", shortCode, i+1, retryCount)
            if i == retryCount - 1 {
                return "", fmt.Errorf("failed to generate unique short code after %d retries", retryCount)
            }
        }
       
    }else {
        // save the original url in DB
        err := s.dbRepo.SaveOriginalURL(ctx, shortCode, originalURL, isCustom, expiry)
        if err!= nil {
        return "", fmt.Errorf("error saving original url to db: %w", err)
        }

    }

    // save the original  url in cache 
    err := s.cacheRepo.SaveOriginalURL(ctx, shortCode, originalURL, expiry)
    if err!= nil {
        log.Printf("Error saving to cache: %v", err)
    }

    return shortCode, nil
    
}