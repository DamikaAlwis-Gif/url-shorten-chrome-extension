package helpers

import (
	"math/rand"
	"strings"
	"time"
)
const(
	// Base62 character set
 base62Chars string = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"
)

// encodeBase62 converts a given integer (timestamp) into a Base62 string
func encodeBase62(num int64) string {
	base := int64(len(base62Chars))
	var sb strings.Builder

	for num > 0 {
		remainder := num % base
		sb.WriteByte(base62Chars[remainder])
		num /= base
	}
	// Reverse the string to get the correct Base62 representation
	return sb.String()
}


// GenerateShortURL generates a short URL with better randomness
func GenerateShortURL(truncateLength int) string {
	// Get the current timestamp in milliseconds
	timestamp := time.Now().UnixNano()
	// Generate a random 16-bit value
	randomPart := int64(rand.Intn(1 << 16))
	// Combine timestamp and random part
	combined := (timestamp << 16) | randomPart
	// Encode the combined value to Base62
	encoded := encodeBase62(combined)
	// Truncate the result
	if len(encoded) > truncateLength {
		encoded = encoded[:truncateLength]
	}

	return encoded
}

