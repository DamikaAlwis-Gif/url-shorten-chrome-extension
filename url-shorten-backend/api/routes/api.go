package routes
import (
	"github.com/gin-gonic/gin"
	"github.com/DamikaAlwis-Gif/shorten-url-app/middleware"
)



func SetupRoutes(r *gin.Engine){
	// define routes
	r.POST("/shorten",middleware.RateLimitMiddleware(),shortenURL)
	r.GET("/:url",resolveURL)
	r.GET("/", func(c *gin.Context){
		c.JSON(200, gin.H{"message": "Welcome to the URL shortening service!"})
	})

}