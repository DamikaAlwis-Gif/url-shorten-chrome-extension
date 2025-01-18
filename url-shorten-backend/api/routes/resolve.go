package routes

import (
	"context"
	"errors"
	"log"
	// "fmt"
	"net/http"
	"time"
	"github.com/DamikaAlwis-Gif/shorten-url-app/custom_errors"
	// "github.com/DamikaAlwis-Gif/shorten-url-app/repository"
	"github.com/DamikaAlwis-Gif/shorten-url-app/service"
	"github.com/gin-gonic/gin"
)

func resolveURL(c *gin.Context, UrlSrv *service.URLService, clSrv *service.ClickLogService){
	ctx := c.Request.Context()
	shortCode := c.Param("short_code")

	// Resolve the URL
	originalURL , err := UrlSrv.ResolveShortURL(ctx, shortCode)
	
  if err != nil {
    if errors.Is(err, custom_errors.ErrShortURLNotFound) {
      c.JSON(http.StatusNotFound, gin.H{"error": "Short URL not found"})
			return
				
    }else if errors.Is(err, custom_errors.ErrURLExpired){
			c.JSON(http.StatusGone, gin.H{"error": "URL has expired"})
			return
		}else {
      c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			return
    }

  }else{

		// go func() {
    // bgCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    // defer cancel() // Ensure that the context is canceled when the task finishes

    // err := clSrv.LogClick(bgCtx, shortCode, c.ClientIP(), c.Request.UserAgent())
    // if err != nil {
    //     log.Printf("Error logging click for short URL %s: %v", shortCode, err)
    // } else {
    //     log.Printf("Published click for short URL %s from IP %s", shortCode, c.ClientIP())
    // }
		// }()
		// Extract data required for logging
	clientIP := c.ClientIP()
	userAgent := c.Request.UserAgent()

	// Start logging in a background goroutine
	go func(shortCode, clientIP, userAgent string) {
		bgCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second) // Independent context
		defer cancel()

		err := clSrv.LogClick(bgCtx, shortCode, clientIP, userAgent)
		if err != nil {
			log.Printf("Error logging click for short URL %s: %v", shortCode, err)
		} else {
			log.Printf("Published click for short URL %s from IP %s", shortCode, clientIP)
		}
	}(shortCode, clientIP, userAgent) // Pass the extracted data as parameters to the goroutine


		
		c.Redirect(http.StatusMovedPermanently, originalURL)
	}
	

	
}

