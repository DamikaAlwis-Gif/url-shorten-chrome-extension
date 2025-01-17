package routes

import (
	// "context"
	"errors"
	// "fmt"
	"net/http"
	// "time"
	"github.com/DamikaAlwis-Gif/shorten-url-app/custom_errors"
	// "github.com/DamikaAlwis-Gif/shorten-url-app/repository"
	"github.com/DamikaAlwis-Gif/shorten-url-app/service"
	"github.com/gin-gonic/gin"
)

func resolveURL(c *gin.Context, srv *service.URLService){
	ctx := c.Request.Context()
	shortCode := c.Param("short_code")

	// Resolve the URL
	originalURL , err := srv.ResolveShortURL(ctx, shortCode)
	
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
		c.Redirect(http.StatusMovedPermanently, originalURL)
	}
	

	
}

// go func(){

	// 	logCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
  //   defer cancel()
	// 	err := repository.LogClick(logCtx,srv, shortURL, c.ClientIP(), c.Request.UserAgent())
	// 	if err!= nil {
	// 		fmt.Printf("Error logging click for short URL %s: %v\n", shortURL, err)
	// 	}
	// }()

	// redirect to the original url