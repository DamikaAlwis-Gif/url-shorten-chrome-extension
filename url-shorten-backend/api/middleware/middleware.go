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

// RateLimitMiddleware implements rate limiting
func RateLimitMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Create a context tied to the current HTTP request
		ctx := c.Request.Context()

		// Get Redis client
		rdb := &database.Redis{}
		redisClient, err := rdb.GetDBClient(ctx)
		if err != nil {
			log.Printf("Error getting Redis client: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
			c.Abort()
			return
		}

		// Get the client's IP address
		ipAddress := c.ClientIP()

		// Check the remaining quota for the IP
		val, err := redisClient.Get(ctx, ipAddress).Result()
		if err != nil {
			// If the key does not exist, initialize the quota
			if errors.Is(err, redis.Nil) {
				defaultQuota := config.AppConfig.APIQuota
				quotaResetTime := config.AppConfig.QuotaResetTime

				if err := redisClient.Set(ctx, ipAddress, defaultQuota, quotaResetTime).Err(); err != nil {
					log.Printf("Error initializing quota for IP %s: %v", ipAddress, err)
					c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
					c.Abort()
					return
				}
			} else {
				// Handle other Redis errors
				log.Printf("Redis error for IP %s: %v", ipAddress, err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
				c.Abort()
				return
			}
		}

		// Convert the remaining quota value to an integer
		remainingQuota, err := strconv.Atoi(val)
		if err != nil {
			log.Printf("Error converting remaining quota value: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
			c.Abort()
			return
		}

		// If the remaining quota is less than or equal to 0, return rate limit exceeded
		if remainingQuota <= 0 {
			// Retrieve the TTL (time-to-live) to inform the client when the quota will reset
			ttl, err := redisClient.TTL(ctx, ipAddress).Result()
			if err != nil {
				log.Printf("Error retrieving TTL for IP %s: %v", ipAddress, err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
				c.Abort()
				return
			}

			// Respond with rate limit exceeded message, along with the TTL (when the rate limit resets)
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error": "Rate limit exceeded",
				"ttl":   ttl.Seconds(),
			})
			c.Abort() // Abort the request processing
			return
		}

		// Proceed to the next handler if the rate limit is not exceeded
		c.Next()
	}
}
