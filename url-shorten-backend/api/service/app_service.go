package service

import "github.com/DamikaAlwis-Gif/shorten-url-app/repository"

type AppService struct {
	URLService       *URLService
	RateLimitService *RateLimitService
}

func NewAppService(cacheRepo repository.CacheRepository, dbRepo repository.DBRepository) *AppService {
	urlService := &URLService{cacheRepo : cacheRepo, dbRepo:  dbRepo }
	rateLimitService := &RateLimitService{repo: cacheRepo}

	return &AppService{URLService: urlService, RateLimitService: rateLimitService}
}