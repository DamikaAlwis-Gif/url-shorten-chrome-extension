package service

import (
  "context"
  "time"
	"github.com/DamikaAlwis-Gif/shorten-url-app/config"
	"log"
	"errors"
  "github.com/DamikaAlwis-Gif/shorten-url-app/custom_errors"
  "github.com/DamikaAlwis-Gif/shorten-url-app/repository"
)

type RateLimitService struct {
	repo repository.CacheRepository
}

func NewRateLimitService(repo repository.CacheRepository) *RateLimitService {
	return &RateLimitService{repo: repo}
}

// HandleRateLimit checks and enforces the rate limit based on IP address
func (rls *RateLimitService) HandleRateLimit(ctx context.Context, ipAddress string) error {
	// Retrieve current quota
	quota, err := rls.repo.GetQuota(ctx, ipAddress)
	if err != nil {
		if errors.Is(err, custom_errors.ErrKeyNotFound) {
			// If no quota is found, initialize with default quota
			return rls.initializeQuota(ctx, ipAddress)
		}
		log.Printf("Error retrieving quota for IP %s: %v", ipAddress, err)
		return err
	}

	// If quota is <= 0, return rate limit exceeded error
	if quota <= 0 {
		return rls.handleRateLimitExceeded(ctx, ipAddress)
	}
	

	return nil
}

// initializeQuota sets the initial quota and TTL for a new IP address
func (rls *RateLimitService) initializeQuota(ctx context.Context, ipAddress string) error {
	defaultQuota := config.AppConfig.APIQuota
	quotaResetTime := config.AppConfig.QuotaResetTime

	// Set the quota in the cache with the specified TTL
	err := rls.repo.SetQuota(ctx, ipAddress, defaultQuota, quotaResetTime)
	if err != nil {
		log.Printf("Error initializing quota for IP %s: %v", ipAddress, err)
		return err
	}

	return nil
}

// handleRateLimitExceeded is called when the quota is exceeded
func (rls *RateLimitService) handleRateLimitExceeded(ctx context.Context, ipAddress string) error {
	// Get the remaining TTL (time-to-live) for the current quota
	ttl, err := rls.repo.GetTTL(ctx, ipAddress)
	if err != nil {
		log.Printf("Error retrieving TTL for IP %s: %v", ipAddress, err)
		return err
	}

	// Return a rate limit exceeded error with the TTL
	return &custom_errors.ErrRateLimitExceeded{TTL : ttl}
}

func (rls *RateLimitService) DecrementQuota(ctx context.Context, ipAddress string)(int , error){
	ramainingQuota , err := rls.repo.DecrementQuota(ctx, ipAddress)
	if err!= nil {
    log.Printf("Error decrementing quota for IP %s: %v", ipAddress, err)
    return -1, err
  }
	return ramainingQuota, nil
}

func (rls *RateLimitService) GetQuotaResetTime(ctx context.Context, ipAddress string) (time.Duration, error) {
	quotaResetTime, err := rls.repo.GetTTL(ctx, ipAddress)
  if err!= nil {
    log.Printf("Error retrieving quota reset time for IP %s: %v", ipAddress, err)
    return 0, err
  }
  return quotaResetTime, nil
}

