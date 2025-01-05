package helpers

import (
	"math/big"
	"net/url"
	"os"
	"strings"
	"fmt"
	"github.com/google/uuid"
	
)
func IsDomainURL(url_string string) (bool, error) {
	parsed_url , err := url.Parse(url_string) ; if err != nil {
		return false, err
	}
	

	if parsed_url.Host == os.Getenv("HOST"){
		return true, nil
	}
	return false, nil
}


// for ascii characters len(string) ok
// utf8.RuneCountInString(string)

func customBaseEncode(num *big.Int, char_set string) string{
	var encoded strings.Builder
	base := big.NewInt(int64(len(char_set)))

	mod := new(big.Int)
	for num.Cmp(big.NewInt(0)) > 0 {
		num.DivMod(num, base, mod)
		encoded.WriteByte(char_set[mod.Int64()]) // Int64 is from big package
	}
	return encoded.String()
}

func GenarateShortCode(length int) (string , error) {
	const char_set = "23456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz"
	uuid := uuid.New()// genarate uuid
	num := new(big.Int).SetBytes(uuid[:]) // genarate a number from uuid
	short_code := customBaseEncode(num, char_set) // encode the number using a base
	if length > len(short_code) {
	  return "" , fmt.Errorf("length should be less than or equal to %d", len(short_code))
	}
	return short_code[:length], nil

}