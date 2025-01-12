package auth

import (
	"crypto/rand"
	"encoding/hex"
)

// Generates a secure random 16 byte API key
//
// Returns the key or an error if one occured.
func GenerateApiKey() (string, error) {
	keyBytes := make([]byte, 16)
	if _, err := rand.Read(keyBytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(keyBytes), nil
}
