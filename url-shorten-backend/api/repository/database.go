package repository

import (
	"context"
	"time"

	"github.com/DamikaAlwis-Gif/shorten-url-app/database"
)

type DBRepository interface{
	
	FindOriginalURLDetailsByShortCode(ctx context.Context, shortCode string) (*database.ShortenURL, error)
	SaveOriginalURL(ctx context.Context, shortCode, originalURL string, isCustom bool, expiry time.Duration) error
	
}