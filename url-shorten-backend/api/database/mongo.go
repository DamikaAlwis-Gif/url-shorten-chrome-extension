package database

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"github.com/DamikaAlwis-Gif/shorten-url-app/config"
)

type MongoDB struct{
	Client *mongo.Client
	once sync.Once

}
func (db *MongoDB)InitDB(parentCtx context.Context) error{

	var initErr error
	db.once.Do( func(){

		// Set up a timeout context for the connection
		ctx, cancel := context.WithTimeout(parentCtx, 10*time.Second)
		defer cancel()

		mongoURI := config.AppConfig.MongoDBAddr
		if mongoURI == "" {
			initErr = fmt.Errorf("missing MONGODB_URI environment variable")
			return
		}
		// Set client options
		clientOptions := options.Client().ApplyURI(mongoURI)

		// connect to MongoDB
		var err error
		db.Client, err = mongo.Connect(ctx, clientOptions)
		if err!= nil{
			initErr = fmt.Errorf("failed to connect to MongoDB at URI %s: %w", mongoURI, err)
			return
		}

		// Test the connection
		err = db.Client.Ping(ctx, nil)
		if err!= nil{
			initErr = fmt.Errorf("failed to ping MongoDB: %w", err)
			return
    }
		log.Println("connected to MongoDB!")

	})
	return initErr
	
}

// get mongodb client instance
func (db *MongoDB)GetDBClient() *mongo.Client{
	
	return db.Client
}

// get a specific MongoDB collection
func (db *MongoDB)GetCollection(databaseName string, collectionName string) (*mongo.Collection, error) {
	if db.Client == nil{
		return nil, fmt.Errorf("MongoDB client is not initialized")
	}
	return db.Client.Database(databaseName).Collection(collectionName), nil
}

//claose the MongoDB connection
func (db *MongoDB) CloseDBConnection(parentCtx context.Context){
	if db.Client != nil{
		ctx, cancel := context.WithTimeout(parentCtx, 10*time.Second)
		defer cancel()

		err := db.Client.Disconnect(ctx)
		if err!= nil{
			log.Printf("Failed to disconnect MongoDB: %v", err)
		}else {
			log.Println("MongoDB connection closed")
		}
	}
}
