package utils

import (
	"crypto/sha1"
	"fmt"
)

func HashPassword(password string) string {
	hasher := sha1.New()
	hasher.Write([]byte(password))
	sha := fmt.Sprintf("%X", hasher.Sum(nil))
	return sha
}
