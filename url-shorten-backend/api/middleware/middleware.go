package middleware

import (
	// "errors"
	// "log"
	"net/http"
	"github.com/DamikaAlwis-Gif/shorten-url-app/service"
	// "strconv"
	// "github.com/DamikaAlwis-Gif/shorten-url-app/config"
	"github.com/DamikaAlwis-Gif/shorten-url-app/custom_errors"
	// "github.com/DamikaAlwis-Gif/shorten-url-app/repository"
	// "github.com/DamikaAlwis-Gif/shorten-url-app/service"
	"github.com/gin-gonic/gin"
	// "github.com/go-redis/redis/v8"
)


// RateLimitMiddleware checks the rate limit for the current IP address
func RateLimitMiddleware(rls *service.RateLimitService) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		ipAddress := c.ClientIP()

		// Handle rate limiting
		if err := rls.HandleRateLimit(ctx, ipAddress); err != nil {
			if rateLimitErr, ok := err.(*custom_errors.ErrRateLimitExceeded); ok {
				// Respond with rate limit exceeded message
				c.JSON(http.StatusTooManyRequests, gin.H{
					"error": "Rate limit exceeded",
					"ttl":   rateLimitErr.TTL,
				})
				c.Abort() // Prevent further processing
				return
			}

			// Handle internal errors
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
			c.Abort()
			return
		}

		c.Next()
	}
}

// func RateLimitMiddleware(repo repository.CacheRepository) gin.HandlerFunc {
//     return func(c *gin.Context) {
//         ctx := c.Request.Context()
//         ipAddress := c.ClientIP()

//         // Get the current quota
//         quota, err := repo.GetQuota(ctx, ipAddress)
//         if err != nil {
//             if errors.Is(err, custom_errors.ErrKeyNotFound) {
//                 // Initialize quota if not found
//                 defaultQuota := config.AppConfig.APIQuota
//                 quotaResetTime := config.AppConfig.QuotaResetTime

//                 if err := repo.SetQuota(ctx, ipAddress, defaultQuota, quotaResetTime); err != nil {
//                     log.Printf("Error initializing quota for IP %s: %v", ipAddress, err)
//                     c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
//                     c.Abort()
//                     return
//                 }
//             } else {
//                 log.Printf("Error retrieving quota for IP %s: %v", ipAddress, err)
//                 c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
//                 c.Abort()
//                 return
//             }
//         } else {
//             // If the quota is less than or equal to 0, block the request
//             if quota <= 0 {
//                 // Use the GetTTL method to retrieve the remaining TTL
//                 ttl, err := repo.GetTTL(ctx, ipAddress)
//                 if err != nil {
//                     log.Printf("Error retrieving TTL for IP %s: %v", ipAddress, err)
//                     c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
//                     c.Abort()
//                     return
//                 }

//                 c.JSON(http.StatusTooManyRequests, gin.H{
//                     "error": "Rate limit exceeded",
//                     "ttl":   ttl.Seconds(),
//                 })
//                 c.Abort()
//                 return
//             }

//         }

//         c.Next()
//     }
// }



// RateLimitMiddleware implements rate limiting
// func RateLimitMiddleware(srv *service.Service,) gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		// Create a context tied to the current HTTP request
// 		ctx := c.Request.Context()

// 		// Get Redis client
// 		rdb := srv.Redis.GetDBClient()
// 		// Get the client's IP address
// 		ipAddress := c.ClientIP()

// 		// Check the remaining quota for the IP
// 		val, err := rdb.Get(ctx, ipAddress).Result()

// 		if err != nil {
// 			// If the key does not exist, initialize the quota
// 			if errors.Is(err, redis.Nil) {
// 				defaultQuota := config.AppConfig.APIQuota
// 				quotaResetTime := config.AppConfig.QuotaResetTime

// 				if err := rdb.Set(ctx, ipAddress, defaultQuota, quotaResetTime).Err(); err != nil {
// 					log.Printf("Error initializing quota for IP %s: %v", ipAddress, err)
// 					c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
// 					c.Abort()
// 					return
// 				}
				

// 			} else {
// 				// Handle other Redis errors
// 				log.Printf("Redis error for IP %s: %v", ipAddress, err)
// 				c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
// 				c.Abort()
// 				return
// 			}

// 		} else {
// 			// Convert the remaining quota value to an integer
// 			remainingQuota, err := strconv.Atoi(val)
// 			if err != nil {
// 				log.Printf("Error converting remaining quota value: %v", err)
// 				c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
// 				c.Abort()
// 				return
// 			}

// 			// If the remaining quota is less than or equal to 0, return rate limit exceeded
// 			if remainingQuota <= 0 {
// 				// Retrieve the TTL (time-to-live) to inform the client when the quota will reset
// 				ttl, err := rdb.TTL(ctx, ipAddress).Result()
// 				if err != nil {
// 					log.Printf("Error retrieving TTL for IP %s: %v", ipAddress, err)
// 					c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
// 					c.Abort()
// 					return
// 				}

// 				// Respond with rate limit exceeded message, along with the TTL (when the rate limit resets)
// 				c.JSON(http.StatusTooManyRequests, gin.H{
// 					"error": "Rate limit exceeded",
// 					"ttl":   ttl.Seconds(),
// 				})
// 				c.Abort() // Abort the request processing
// 				return
// 			}
// 		}

		
// 		// Proceed to the next handler if the rate limit is not exceeded
// 		c.Next()
// 	}
// }
