package routes
import (
	"github.com/gin-gonic/gin"
	"github.com/DamikaAlwis-Gif/shorten-url-app/middleware"
	"github.com/DamikaAlwis-Gif/shorten-url-app/service"
)



func SetupRoutes(r *gin.Engine, srv *service.AppService){
	// define routes
	r.POST("/shorten",middleware.RateLimitMiddleware(srv.RateLimitService), func( c *gin.Context){
		shortenURL(c, srv.URLService, srv.RateLimitService)
	})
	r.GET("/:short_code", func(c *gin.Context){
		resolveURL(c, srv.URLService)
	})
	r.GET("/", func(c *gin.Context){
		c.JSON(200, gin.H{"message": "Welcome to the URL shortening service!"})
	})

}