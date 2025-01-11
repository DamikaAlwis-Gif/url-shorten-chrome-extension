package routes
import (
	"github.com/gin-gonic/gin"
	"github.com/DamikaAlwis-Gif/shorten-url-app/middleware"
	"github.com/DamikaAlwis-Gif/shorten-url-app/service"
)



func SetupRoutes(r *gin.Engine, srv *service.Service){
	// define routes
	r.POST("/shorten",middleware.RateLimitMiddleware(srv), func( c *gin.Context){
		shortenURL(c, srv)
	})
	r.GET("/:url", func(c *gin.Context){
		resolveURL(c, srv)
	})
	r.GET("/", func(c *gin.Context){
		c.JSON(200, gin.H{"message": "Welcome to the URL shortening service!"})
	})

}