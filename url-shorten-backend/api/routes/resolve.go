package routes

import (
	"time"
	"context"
	"fmt"
	"errors"
	"net/http"
	"github.com/gin-gonic/gin"
	"github.com/DamikaAlwis-Gif/shorten-url-app/repository"
	"github.com/DamikaAlwis-Gif/shorten-url-app/custom_errors"
	
)

func resolveURL(c *gin.Context){
	ctx := c.Request.Context()
	shortURL := c.Param("url")
	
	// Resolve the URL
  originalURL, err := repository.ResolveShortURL(ctx, shortURL)
	
  if err != nil {
    if errors.Is(err, custom_errors.ErrShortURLNotFound) {
        c.JSON(http.StatusNotFound, gin.H{"error": "Short URL not found"})
    } else {
      c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
    }
    return
  }

	go func(){

		logCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()
		err := repository.LogClick(logCtx, shortURL, c.ClientIP(), c.Request.UserAgent())
		if err!= nil {
			fmt.Printf("Error logging click for short URL %s: %v\n", shortURL, err)
		}
	}()

	// redirect to the original url
	c.Redirect(http.StatusMovedPermanently, originalURL)
	
}