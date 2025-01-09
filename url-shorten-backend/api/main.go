package main

import (
	"github.com/gin-gonic/gin"
	"log"
	"github.com/DamikaAlwis-Gif/shorten-url-app/routes"
	"github.com/DamikaAlwis-Gif/shorten-url-app/config"
	"github.com/DamikaAlwis-Gif/shorten-url-app/database"
	"github.com/joho/godotenv"
	"github.com/gin-contrib/cors"
)


func main(){
	r := gin.Default()

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
	redisClient := database.InitRedis()
	defer redisClient.Close() // close connection when the app shuts down

	// connect to the MongoDB db
	database.InitMongoDB()
	defer database.CloseMongoDBConnection()

	routes.SetupRoutes(r)
	r.Run(AppConfig.AppPort)

}