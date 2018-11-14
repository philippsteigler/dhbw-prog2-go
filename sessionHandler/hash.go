package sessionHandler

import (
	"crypto/sha256"
	"encoding/base64"
	"math/rand"
	"strings"
	"time"
)

const letterBytes = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func init() {
	rand.Seed(time.Now().UnixNano())
}

// How to generate a random string as salt:
// https://stackoverflow.com/questions/22892120/how-to-generate-a-random-string-of-a-fixed-length-in-go
func generateSalt() string {
	salt := make([]byte, 16)

	for i := range salt {
		salt[i] = letterBytes[rand.Intn(len(letterBytes))]
	}

	return string(salt)
}

// HashString() returns the passwords's SHA256 hash sum and the salt that has first been added
// This function is used to hash the password for a new user
func HashString(text string) (string, string) {
	salt := generateSalt()

	tmp := []string{text, salt}
	saltedString := strings.Join(tmp, "")

	bytes := sha256.Sum256([]byte(saltedString))
	hashedString := base64.URLEncoding.EncodeToString(bytes[:])

	return hashedString, salt
}

// GetHash() returns the passwords's SHA256 hash sum
// This function is used to calculate the hash sum of a password using a specific salt value
func GetHash(text, salt string) string {
	tmp := []string{text, salt}
	saltedString := strings.Join(tmp, "")

	bytes := sha256.Sum256([]byte(saltedString))
	hashedString := base64.URLEncoding.EncodeToString(bytes[:])

	return hashedString
}
