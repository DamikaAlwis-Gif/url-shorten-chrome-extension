package database

import(
	"context"
	"log"
	"os"
	"sync"
	"time"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

)

var mongo_client *mongo.Client
var once sync.Once

func InitMongoDB(){
	once.Do( func(){
		mongoURI := os.Getenv("MONGODB_URI")
		if mongoURI == "" {
			log.Fatal("Missing MONGODB_URI environment variable")
		}
		// Set client options
		clientOptions := options.Client().ApplyURI(mongoURI)

		// Set up a timeout context for the connection
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		// connect to MongoDB
		var err error 
		mongo_client, err = mongo.Connect(ctx, clientOptions)
		if err!= nil{
			log.Fatalf("Failed to connect to MongoDB: %v", err)
		}

		// Test the connection
		err = mongo_client.Ping(ctx, nil)
		if err!= nil{
      log.Fatalf("Failed to ping MongoDB: %v", err)
    }
		log.Println("Connected to MongoDB!")

	})
	
}

// get mongodb client instance
func GetMongoClient() *mongo.Client {
	if mongo_client == nil{
		log.Println("MongoDB client not initialized. Initializing now...")
		InitMongoDB()
	}
	return mongo_client  
}
// get a specific MongoDB collection
func GetCollection(databaseName string, collectionName string) *mongo.Collection {
	return mongo_client.Database(databaseName).Collection(collectionName)
}

//claose the MongoDB connection

func CloseMongoDBConnection(){
	if mongo_client != nil{
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		err := mongo_client.Disconnect(ctx)
		if err!= nil{
			log.Printf("Failed to disconnect MongoDB: %v", err)
		}else {
			log.Println("MongoDB connection closed")
		}
	}
}
