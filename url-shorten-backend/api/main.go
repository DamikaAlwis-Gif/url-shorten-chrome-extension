package main

import (
	"context"
	"log"

	"github.com/DamikaAlwis-Gif/shorten-url-app/config"
	"github.com/DamikaAlwis-Gif/shorten-url-app/database"
	"github.com/DamikaAlwis-Gif/shorten-url-app/repository"
	"github.com/DamikaAlwis-Gif/shorten-url-app/routes"
	"github.com/DamikaAlwis-Gif/shorten-url-app/service"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main(){
	r := gin.Default()
	// set up concelabel context
	ctx, cancel  := context.WithCancel(context.Background())
	defer cancel()

	// load environment variables from .env file
	err := godotenv.Load()
	if err!= nil {
		log.Fatalf("Error loading .env file: %v", err.Error())
	}
	AppConfig := config.LoadConfig()
	log.Printf("Loaded .env file: %v", AppConfig)

	// Set up CORS configuration
	setupCORS(r)

	// Initialize databases
	redis, mongoDB := initializeDatabases(ctx)
	defer redis.CloseDBConnection()
	defer mongoDB.CloseDBConnection(ctx)

	
	dbRepo := repository.NewMongoRepository(mongoDB.GetDBClient())
	cacheRepo := repository.NewRedisRepository(redis.GetDBClient())
	srv := service.NewAppService(cacheRepo, dbRepo)

	routes.SetupRoutes(r, srv)
	r.Run(AppConfig.AppPort)

}



// setupCORS configures the CORS middleware for the Gin router.
func setupCORS(r *gin.Engine) {
	corsConfig := cors.Config{
		AllowOrigins:     []string{"*"}, // Replace "*" with specific domains for security
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
	}
	r.Use(cors.New(corsConfig))
}


// initializeDatabases sets up connections to Redis and MongoDB.
func initializeDatabases(ctx context.Context) (*database.Redis, *database.MongoDB) {
	// Initialize Redis
	redis := &database.Redis{}
	if err := redis.InitDB(ctx); err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}

	// Initialize MongoDB
	mongoDB := &database.MongoDB{}
	if err := mongoDB.InitDB(ctx); err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}

	return redis, mongoDB
}

