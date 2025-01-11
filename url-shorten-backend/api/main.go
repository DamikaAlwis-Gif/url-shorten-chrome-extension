package main

import (
	"context"
	"log"
	"github.com/DamikaAlwis-Gif/shorten-url-app/config"
	"github.com/DamikaAlwis-Gif/shorten-url-app/database"
	"github.com/DamikaAlwis-Gif/shorten-url-app/service"
	"github.com/DamikaAlwis-Gif/shorten-url-app/routes"
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

	// Define CORS configuration
	corsConfig := cors.Config{
		// Allow all origins (can also restrict to specific domains)
		AllowOrigins: []string{"*"}, // "*" allows all origins, replace with specific domains for tighter security
		// Allow HTTP methods (GET, POST, PUT, DELETE, etc.)
		AllowMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		// Allow specific headers
		AllowHeaders: []string{"Origin", "Content-Type", "Authorization"},
		// Allow credentials (cookies, authorization headers, etc.)
		AllowCredentials: true,
	}

	// Use CORS middleware with the specified configuration
	r.Use(cors.New(corsConfig))

	// connect to the redis db
	redis := &database.Redis{}
	if err := redis.InitDB(ctx); err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}
	defer redis.CloseDBConnection()

	// connect to the MongoDB db
	mongoDB := &database.MongoDB{}
	if err := mongoDB.InitDB(ctx); err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	defer mongoDB.CloseDBConnection(ctx)

	srv := &service.Service{Redis: redis, MongoDB: mongoDB}

	routes.SetupRoutes(r, srv )
	r.Run(AppConfig.AppPort)

}