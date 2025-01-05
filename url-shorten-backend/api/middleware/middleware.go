package middleware

import (
	"errors"
	"log"
	"net/http"
	"strconv"
	"github.com/DamikaAlwis-Gif/shorten-url-app/config"
	"github.com/DamikaAlwis-Gif/shorten-url-app/database"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
)

func RateLimitMiddleware() gin.HandlerFunc{
	return func(c *gin.Context){
  // implement rate limiting
	rdb := database.GetRedisClient()
	// get the client's IP address
	ipAddress := c.ClientIP()
	// check the remaing quota
	val, err := rdb.Get(database.Ctx, ipAddress).Result()
	remaining_quota, _ := strconv.Atoi(val)

	if err != nil {
		// if the key does not exist
		if errors.Is(err,redis.Nil) {
			default_quota:= config.AppConfig.APIQuoata
			quota_reset_time_min := config.AppConfig.QuotaResetTime
			
			if err =rdb.Set(database.Ctx, ipAddress, default_quota,quota_reset_time_min).Err(); err != nil{
				log.Printf("Error initializing quota for IP %s: %v", ipAddress, err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
				c.Abort()
				return
			}

  	}else {
			// Other Redis errors
			log.Printf("Redis error for IP %s: %v", ipAddress, err)
			c.JSON(http.StatusInternalServerError, gin.H{"error" : "Internal Server Error"})
			c.Abort()
			return
		}
	}
	if remaining_quota <= 0 {
		// ttl -> time to live
		// Retrieve TTL to inform the user when the quota resets
		ttl, err := rdb.TTL(database.Ctx, ipAddress).Result()
		if err != nil {
    	log.Printf("Error retrieving TTL for IP %s: %v", ipAddress, err)
    	c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			c.Abort()
    	return
		}
		// respond with rate limit  exceeded message
    c.JSON(http.StatusTooManyRequests, gin.H{
			"error": "Rate limit exceeded",
			"ttl" : ttl.Seconds(),
			})
		c.Abort()		
     return
  	}
	// proceed to the next handler
	c.Next()
  }
	
}