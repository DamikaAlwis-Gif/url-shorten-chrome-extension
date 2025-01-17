package database

import(
	"time"
)

type ShortenURL struct{
	ShortCode string `bson:"short_code"`
	OriginalURL string `bson:"original_url"`
	CreatedAt time.Time `bson:"created_at"`
	Expiry time.Time `bson:"expiry"`
	IsCustom bool `bson:"is_custom"`
}

type Click struct{
	ShortURL string `bson:"short_url"`
	Timestamp time.Time `bson:"timestamp"`
	IPAddress string `bson:"ip_address"`
	UserAgent string `bson:"user_agent"`
}

