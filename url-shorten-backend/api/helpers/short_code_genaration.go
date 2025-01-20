package helpers
import (
	"crypto/md5"
	"math/big"
	"math/rand"
	"time"
)

const (
	charSet = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"
	randomLength = 6 
	codeLength = 8

)

// Generate a random string of the given length
func randomString(length int) string {
	rand.NewSource(time.Now().UnixNano())
	bytes := make([]byte, length)
	for i := range bytes {
		bytes[i] = charSet[rand.Intn(len(charSet))]
	}
	return string(bytes)
}

// Base62 encoding function
func base62Encode(num *big.Int) string {
	base := big.NewInt(int64(len(charSet)))
	zero := big.NewInt(0)
	mod := new(big.Int)

	var encoded string
	for num.Cmp(zero) > 0 {
		num.DivMod(num, base, mod) // Divide num by base, get the remainder in mod
		encoded = string(charSet[mod.Int64()]) + encoded
	}
	return encoded
}

// Generate a short code using Base62
func GenerateShortURL(longURL string) (string, error) {
	// Combine the long URL with a random string
	randomPart := randomString(randomLength)
	input := longURL + randomPart

	// Compute the MD5 hash
	hash := md5.Sum([]byte(input))

	// Convert the hash bytes into a big.Int
	num := new(big.Int).SetBytes(hash[:])

	// Encode the number in Base62
	encoded := base62Encode(num)

	// Ensure the short code is the desired length
	if len(encoded) < codeLength {
		// Pad with leading characters if too short
		encoded = randomString(codeLength-len(encoded)) + encoded
	}

	// Truncate to the desired length
	return encoded[:codeLength], nil
}



// import (
// 	"math/rand"
// 	"strings"
// 	"time"
// )
// const(
// 	// Base62 character set
//  base62Chars string = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"
// )

// // encodeBase62 converts a given integer (timestamp) into a Base62 string
// func encodeBase62(num int64) string {
// 	base := int64(len(base62Chars))
// 	var sb strings.Builder

// 	for num > 0 {
// 		remainder := num % base
// 		sb.WriteByte(base62Chars[remainder])
// 		num /= base
// 	}
// 	// Reverse the string to get the correct Base62 representation
// 	return sb.String()
// }


// // GenerateShortURL generates a short URL with better randomness
// func GenerateShortURL(truncateLength int) string {
// 	// Get the current timestamp in milliseconds
// 	timestamp := time.Now().UnixNano()
// 	// Generate a random 16-bit value
// 	randomPart := int64(rand.Intn(1 << 16))
// 	// Combine timestamp and random part
// 	combined := (timestamp << 16) | randomPart
// 	// Encode the combined value to Base62
// 	encoded := encodeBase62(combined)
// 	// Truncate the result
// 	if len(encoded) > truncateLength {
// 		encoded = encoded[:truncateLength]
// 	}

// 	return encoded
// }

