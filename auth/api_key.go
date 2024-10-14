package auth

import (
	"crypto/rand"
	"eigen_db/constants"
	"encoding/hex"
	"os"
)

// Stores the passed in `key` in environment variables and on disk at `apiKeyFilePath`.
//
// Returns an error if one occured.
func storeApiKey(key string, apiKeyFilePath string) error {
	if err := os.Setenv(constants.ENV_VAR_API_KEY_NAME, key); err != nil {
		return err
	}
	if err := os.WriteFile(apiKeyFilePath, []byte(key), constants.API_KEY_FILE_CHMOD); err != nil {
		return err
	}
	return nil
}

// Generates a secure random 16 byte API key
//
// Returns the key or an error if one occured.
func generateApiKey() (string, error) {
	keyBytes := make([]byte, 16)
	if _, err := rand.Read(keyBytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(keyBytes), nil
}

// Generates/fetches API key and loads it into the environment variables
//
// Loads API key stored in `apiKeyFilePath` into environment variables for quick access.
//
// A new key is generated if regen = true, or the file `apiKeyFilePath` does not exist, or a custom key is specified.
//
// Returns the API key, or an error if one occured.
func SetupAPIKey(customApiKey string, regen bool, apiKeyFilePath string) (string, error) {
	if _, err := os.ReadFile(apiKeyFilePath); (err != nil && os.IsNotExist(err)) || regen || customApiKey != "" { // if api key file doesnt exist, or regen = true, or a custom key is passed
		var key string
		if customApiKey == "" {
			key, err = generateApiKey()
			if err != nil {
				return "", err
			}
		} else {
			key = customApiKey
		}
		if err := storeApiKey(key, apiKeyFilePath); err != nil {
			return "", err
		}
		return key, nil
	} else if err != nil {
		return "", err
	}

	// re-use the key on disk
	keyBytes, err := os.ReadFile(apiKeyFilePath)
	if err != nil {
		return "", err
	}
	key := string(keyBytes)
	if err := os.Setenv("EIGENDB_API_KEY", key); err != nil { // load key into env vars for quick access
		return "", err
	}

	return key, nil
}
