package repository

// import (
// 	"context"
// 	"time"
// 	"github.com/DamikaAlwis-Gif/shorten-url-app/database"
// 	// "github.com/DamikaAlwis-Gif/shorten-url-app/service"
// 	"github.com/DamikaAlwis-Gif/shorten-url-app/config"
// )

// func LogClick(ctx context.Context, srv *service.Service, shortURL, ipAddress, userAgent string) error {
// 	mdb := srv.MongoDB
// 	collection, err := mdb.GetCollection(config.AppConfig.MongoDBDatabaseName, config.AppConfig.MongoDBCOllectionNameClicks)
// 	if err != nil {
// 		return err
// 	}

// 	// Create a new click entry
// 	click := database.Click{
// 		ShortURL:  shortURL,
// 		Timestamp: time.Now(),
// 		IPAddress: ipAddress,
// 		UserAgent: userAgent,
// 	}

// 	// Insert the click entry into MongoDB
// 	_, err = collection.InsertOne(ctx, click)
// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }
