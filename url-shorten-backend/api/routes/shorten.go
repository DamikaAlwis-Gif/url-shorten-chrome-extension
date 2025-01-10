package routes

import (
	"net/http"
	"time"
	"log"
	"fmt"
	"github.com/DamikaAlwis-Gif/shorten-url-app/config"
	"github.com/DamikaAlwis-Gif/shorten-url-app/database"
	"github.com/DamikaAlwis-Gif/shorten-url-app/helpers"
	"github.com/DamikaAlwis-Gif/shorten-url-app/custom_errors"
	"github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"errors"
)


type request struct {
	URL          string         `json:"url"`
	CustomShort  string         `json:"short"`
	Expiry			 time.Duration  `json:"expiry"`
}


type response struct {
	URL                string          `json:"url"`
	CustomShort        string          `json:"short"`
	Expiry             time.Duration   `json:"expiry"`
	XRateRemaining     int             `json:"rate_limit"`
	XRateLimitReset    time.Duration   `json:"rate_limit_reset"` 
}

func get_short(custom_short string) (string, error){

	var id string
	var err error

	// there is no custom short given by user, genarate a new short
	if custom_short == "" {
		id, err = helpers.GenarateShortCode(6)
		if err != nil{
			log.Print(err.Error())
			return "", err
		}
	} else {
		// if a custom short given
		id = custom_short
	}
	return id, nil
	
}

func set_short_url(short_url string, original_url string, expiry time.Duration) (error) {
	rdb := database.GetRedisClient()
	url_collection := database.GetCollection("url_shortener","urls")
	var err error

	// if expiry is not provided, set it to 24 hours (default)
	if expiry == 0{
		expiry = time.Duration(24)
	}
	
	count , err := rdb.Exists(database.Ctx, short_url).Result()
	if err != nil {
		log.Print(err.Error())
    return err
	}  

	if count >0 {
		
	}else{
		// if short key doesn't exist -> set it
		if err == redis.Nil{
			err = rdb.Set(database.Ctx, short_key, original_url, expiry * time.Hour).Err()
  		if err!= nil {
    		log.Print(err.Error())
    		return  err
  		}
  		return nil
		}
		log.Print(err.Error())
		return err
	}

	}

	_, err = rdb.Get(database.Ctx,short_key).Result()
	// if there is an error -> short key already exists
	if err == nil{
		return custom_errors.ErrShortKeyExists
	}else {
		// if short key doesn't exist -> set it
		if err == redis.Nil{
			err = rdb.Set(database.Ctx, short_key, original_url, expiry * time.Hour).Err()
  		if err!= nil {
    		log.Print(err.Error())
    		return  err
  		}
  		return nil
		}
		log.Print(err.Error())
		return err
	}
 
}


func shortenURL(c *gin.Context){
	// validate request body
	var req request
	if err := c.ShouldBindJSON(&req); err!= nil {
    c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
    return
  }

	rdb := database.GetRedisClient()
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
	

	short_code , err := get_short(req.CustomShort)
	if err!= nil {
    log.Print(err.Error())
    c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
    return
  }
	err = set_short_url(short_code, req.URL, req.Expiry)
	if err!= nil {
		if errors.Is(err, custom_errors.ErrShortKeyExists) {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
		}

		log.Print(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
    
  }

	ipAddress := c.ClientIP()
	remaining_quota , err := rdb.Decr(database.Ctx, ipAddress).Result()
	if err!= nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
	}
	ttl , err := rdb.TTL(database.Ctx, ipAddress).Result() 
	if err!= nil {
    log.Printf("Error retrieving TTL for IP %s: %v", ipAddress, err)
    c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
    return
  }
	short_url := config.AppConfig.Host + "/"+ short_code

	response := response{URL: req.URL, CustomShort: short_url, Expiry: req.Expiry, XRateRemaining: int(remaining_quota), XRateLimitReset: time.Duration(ttl.Minutes())}

	c.JSON(http.StatusOK, response)

}


