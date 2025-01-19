package service

import "github.com/DamikaAlwis-Gif/shorten-url-app/repository"

type AppService struct {
	URLService       *URLService
	RateLimitService *RateLimitService
	ClickLogService  *ClickLogService
}

func NewAppService(cacheRepo repository.CacheRepository, dbRepo repository.DBRepository, broker repository.MessageBroker) *AppService {
	urlService := &URLService{cacheRepo : cacheRepo, dbRepo:  dbRepo }
	rateLimitService := &RateLimitService{repo: cacheRepo}
	clickLogService := &ClickLogService{broker : broker}

	return &AppService{URLService: urlService, RateLimitService: rateLimitService, ClickLogService: clickLogService}
}