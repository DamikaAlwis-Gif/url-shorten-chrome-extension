package service

import (
    "github.com/DamikaAlwis-Gif/shorten-url-app/database"
)

// Service holds the dependencies for the application
type Service struct {
    Redis   *database.Redis
    MongoDB *database.MongoDB
}
