package routes

import (
	"errors"
	"fmt"
	"github.com/DamikaAlwis-Gif/shorten-url-app/database"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"net/http"
)

func fetch_original_url(short_code string)(string, error){
	// get database client
	rdb := database.GetRedisClient()
	// genarae short key short:abc123
	short_key := fmt.Sprintf("short:%s", short_code) 
	// get the original url from redis
	originalURL, err := rdb.Get(database.Ctx, short_key).Result()
	if err != nil {
		return "", err
	}
	// increment the clicks counter
	clicks_key := fmt.Sprintf("clicks:%s",short_code)
	_, err = rdb.Incr(database.Ctx,clicks_key).Result()
	if err!= nil {
    return "", err
  }
	return originalURL, nil
}


func resolveURL(c *gin.Context){
	short_url := c.Param("url")

	original_url , err := fetch_original_url(short_url)

	if err != nil {
		if errors.Is(err, redis.Nil){
			// short url not found
			c.JSON(http.StatusNotFound, gin.H{"error": "Short url not found"})
			return
		}else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		}
	}
	// redirect to the original url
	c.Redirect(http.StatusMovedPermanently, original_url)
	
}