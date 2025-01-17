package custom_errors

import (
	"errors"
	"fmt"
	"time"
)

var ErrShortKeyExists = errors.New("short key already exists")
var ErrShortURLNotFound = errors.New("short url not found")
var ErrURLExpired = errors.New("url expired")
var ErrKeyNotFound = errors.New("key not found")
// var ErrRateLimitExceeded = errors.New("rate limit exceeded")

type ErrRateLimitExceeded struct {
	TTL time.Duration // Time-to-live in seconds
}

func (e *ErrRateLimitExceeded) Error() string {
	return fmt.Sprintf("Rate limit exceeded, try again in %v minutes", e.TTL.Minutes())
}

