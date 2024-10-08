package redis_utils

import (
	"context"
	"crypto/rand"
	"eigen_db/constants"
	"encoding/hex"
	"fmt"
	"os"
	"time"

	"github.com/redis/go-redis/v9"
)

// Creates a connection to Redis.
//
// If the connection fails, the returned client is nil. Make sure to check
// for any errors after calling this function as this can cause a lot of
// nil pointer dereference errors.
//
// Returns a Redis client or an error if the connection was unsuccessfull.
func GetConnection(ctx context.Context) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", os.Getenv("REDIS_HOST"), os.Getenv("REDIS_PORT")),
		Password: os.Getenv("REDIS_PASS"),
		DB:       0,
	})

	if err := CheckConnection(ctx, client); err != nil {
		return nil, err
	}

	return client, nil
}

// Checks the connection to Redis by pinging it with a 3 secs timeout.
//
// Returns an error if one occurs.
func CheckConnection(ctx context.Context, client *redis.Client) error {
	timeout := time.Second * 3 // 3 secs timeout
	timeoutCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	if _, err := client.Ping(timeoutCtx).Result(); err != nil {
		return err
	}
	return nil
}

// Generates the API key
//
// Once the key is created, it's inserted in Redis, and written to "apiKeyFilePath"
//
// If apiKey != "", then the API key will be set to the value of apiKey.
//
// If apiKey == "", a random API key is generated and overwrites any existing API key.
//
// Returns the API key, or an error if one occured.
func SetupAPIKey(ctx context.Context, client *redis.Client, apiKey string, apiKeyFilePath string) (string, error) {
	if err := CheckConnection(ctx, client); err != nil {
		return "", err
	}

	if apiKey != "" { // if we are explicitly setting a custom key
		if status := client.Set(ctx, "apiKey", apiKey, 0); status.Err() != nil {
			return "", status.Err()
		}
	} else { // no custom key
		keyBytes := make([]byte, 32)
		if _, err := rand.Read(keyBytes); err != nil {
			return "", err
		}
		apiKey = hex.EncodeToString(keyBytes)
		if status := client.Set(ctx, "apiKey", apiKey, 0); status.Err() != nil {
			return "", status.Err()
		}
	}

	val, err := client.Get(ctx, constants.REDIS_API_KEY_NAME).Result()
	if err != nil {
		return "", err
	}

	if err := os.WriteFile(apiKeyFilePath, []byte(val), constants.API_KEY_FILE_CHMOD); err != nil {
		return "", err
	}

	return val, nil
}
