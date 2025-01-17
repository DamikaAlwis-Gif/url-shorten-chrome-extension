package repository

import (
    "context"
    "fmt"
		"time"
		// "go.mongodb.org/mongo-driver/mongo/options"
    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/mongo"
    "github.com/DamikaAlwis-Gif/shorten-url-app/config"
		"github.com/DamikaAlwis-Gif/shorten-url-app/database"
		"github.com/DamikaAlwis-Gif/shorten-url-app/custom_errors"
)

type MongoRepository struct{
	client *mongo.Client
}

func NewMongoRepository (client *mongo.Client) *MongoRepository {
	return &MongoRepository{client: client}
}

// FindShortURL fetches the complete ShortenURL object from the database by its short code.
func (repo *MongoRepository) FindOriginalURLDetailsByShortCode(ctx context.Context, shortCode string) (*database.ShortenURL, error) {
	collection := repo.client.Database(config.AppConfig.MongoDBDatabaseName).Collection(config.AppConfig.MongoDBCollectionNameUrls)

	// Query filter for the short code
	filter := bson.M{"short_code": shortCode}

	// Create a variable to hold the fetched document
	var urlDoc database.ShortenURL

	// Fetch the document from MongoDB
	err := collection.FindOne(ctx, filter).Decode(&urlDoc)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, custom_errors.ErrShortURLNotFound
		}
		return nil, fmt.Errorf("error finding short URL in DB: %w", err)
	}

	return &urlDoc, nil
}

func (repo *MongoRepository) SaveOriginalURL(ctx context.Context, shortCode, originalURL string, isCustom bool, expiry time.Duration) error{
	collection := repo.client.Database(config.AppConfig.MongoDBDatabaseName).Collection(config.AppConfig.MongoDBCollectionNameUrls)
	
		 urlDoc := database.ShortenURL{
        ShortCode:    shortCode,
        OriginalURL: originalURL,
        CreatedAt:   time.Now(),
        Expiry:      time.Now().Add(expiry),
				IsCustom: isCustom,
    }

    _, err := collection.InsertOne(ctx, urlDoc)
    if err != nil {
				if mongo.IsDuplicateKeyError(err){
					return custom_errors.ErrShortKeyExists
				}
        return fmt.Errorf("error storing short URL in DB: %w", err)
    }
    return nil
	
 
}
