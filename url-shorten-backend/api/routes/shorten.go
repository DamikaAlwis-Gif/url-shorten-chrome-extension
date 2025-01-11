package routes

import (
	"errors"
	"log"
	"net/http"
	"time"
	"github.com/DamikaAlwis-Gif/shorten-url-app/config"
	"github.com/DamikaAlwis-Gif/shorten-url-app/custom_errors"
	"github.com/DamikaAlwis-Gif/shorten-url-app/helpers"
	"github.com/DamikaAlwis-Gif/shorten-url-app/repository"
	"github.com/DamikaAlwis-Gif/shorten-url-app/service"
	"github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"
)


type request struct {
	URL          string         `json:"url"`
	CustomShort  string         `json:"short"`
	Expiry			 time.Duration  `json:"expiry"`
}

type response struct {
	URL                string          `json:"url"`
	ShortURL           string          `json:"short"`
	Expiry             time.Duration   `json:"expiry"`
	XRateRemaining     int             `json:"rate_limit"`
	XRateLimitReset    time.Duration   `json:"rate_limit_reset"` 
}




func shortenURL(c *gin.Context, srv *service.Service){
	// validate request body
	var req request
	// get http request context
	ctx := c.Request.Context()

	// validate request
	if err := c.ShouldBindJSON(&req); err!= nil {
    c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
    return
  }

	// check if the input is an actual url
	if !govalidator.IsURL(req.URL){
		c.JSON(http.StatusBadRequest, gin.H{"error":"Invalid URL"})
		return
	}	
	
	// check if the entered url is hostname
	if isDomainURL , err := helpers.IsDomainURL(req.URL); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error validating URL"})
		return
	}else if isDomainURL {
		c.JSON(http.StatusBadRequest, gin.H{"error" : "Can't use this URL"})
	}  
	
	// set the short url in persistent storage and cache 
	shortCode, err := repository.SetShortURL(ctx, srv, req.CustomShort, req.URL, req.Expiry)
	if err != nil {
		if errors.Is(err, custom_errors.ErrShortKeyExists) {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		log.Print(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return 
  }
	// get remaing rate and rate reset time
	ipAddress := c.ClientIP()
	rateLimitResonse, err := repository.GetRateLimit(ctx,srv, ipAddress)
	if err!= nil {
    log.Printf("Error getting rate limit for IP %s: %v", ipAddress, err)
    c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
    return
  }
	remainingQuota := rateLimitResonse.RemainingQuota
	resetAfter := rateLimitResonse.ResetAfter
	
	shortURL := config.AppConfig.Host + "/"+ shortCode

	response := response{URL: req.URL, ShortURL : shortURL, Expiry: req.Expiry, XRateRemaining: remainingQuota, XRateLimitReset: time.Duration(resetAfter.Minutes())}

	c.JSON(http.StatusOK, response)

}


