package custom_errors

import "errors"

var ErrShortKeyExists = errors.New("short key already exists")
var ErrShortURLNotFound = errors.New("short url not found")