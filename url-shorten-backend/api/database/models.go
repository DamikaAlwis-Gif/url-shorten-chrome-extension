package database

import(
	"time"
)

type ShortenURL struct{
	ShortURL string `bson:"short_url"`
	OriginalURL string `bson:"original_url"`
	CreatedAt time.Time `bson:"created_at"`
	Expiry time.Time `bson:"expiry"`
}

type Click struct{
	ShortURL string `bson:"short_url"`
	Timestamp time.Time `bson:"timestamp"`
	IPAddress string `bson:"ip_address"`
	UserAgent string `bson:"user_agent"`
}

