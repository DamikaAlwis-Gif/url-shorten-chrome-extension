package config
import (
	"log"
	"os"
	"strconv"
	"time"
)

type Config struct {
	RedisDbAddr string
	RedisDbPass string
	AppPort  string
	Host string
	APIQuota   int
	QuotaResetTime  time.Duration
	
}
var AppConfig *Config
func LoadConfig() *Config{

	AppConfig = &Config{
		RedisDbAddr:     getEnv("REDIS_DB_ADDR", "localhost:6379"),
    RedisDbPass:     getEnv("REDIS_DB_PASS", ""),
    AppPort:         getEnv("APP_PORT", "8080"),
    Host:           getEnv("HOST", "localhost"),
    APIQuota:       getEnvAsInt("API_QUOTA", 10),
    QuotaResetTime: getEnvAsDuration("QUOTA_RESET_TIME", 30), 
		
	}
	return AppConfig
}

// retrives a string env variable or returns a default value
func getEnv(key string , defualtValue string) string {
	value , exists := os.LookupEnv(key)
	if !exists {
		return defualtValue
	}
	return value
}

// retrives values as int
func getEnvAsInt(key string , defualtValue int) int {
	valueStr , exists := os.LookupEnv(key)
	if !exists {
		return defualtValue
	}
	value , err := strconv.Atoi(valueStr)
	if err!= nil {
    log.Printf("Error parsing %s as int: %v", key, err)
		return defualtValue
  }
	return value

}

// getEnvAsDuration retrieves a duration environment variable (in minutes) or returns a default value
func getEnvAsDuration(key string, defualtValue int) time.Duration{
	minutes := getEnvAsInt(key, defualtValue)
	return time.Duration(minutes) * time.Minute
}



