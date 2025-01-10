package repository

import (
	"context"
	"time"
	"github.com/DamikaAlwis-Gif/shorten-url-app/database"
)

func LogClick(ctx context.Context, shortURL, ipAddress, userAgent string) error {
	mdb := &database.MongoDB{}
	collection, err := mdb.GetCollection("url_shortener", "clicks")
	if err != nil {
		return err
	}

	// Create a new click entry
	click := database.Click{
		ShortURL:  shortURL,
		Timestamp: time.Now(),
		IPAddress: ipAddress,
		UserAgent: userAgent,
	}

	// Insert the click entry into MongoDB
	_, err = collection.InsertOne(ctx, click)
	if err != nil {
		return err
	}

	return nil
}
