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

func GetConnection(ctx context.Context) (*redis.Client, error) { // if connection fails, it returned client is nil
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

func CheckConnection(ctx context.Context, client *redis.Client) error {
	timeout := time.Second * 3 // 3 secs timeout
	timeoutCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	if _, err := client.Ping(timeoutCtx).Result(); err != nil {
		return err
	}
	return nil
}

// generates API key, inserts it to Redis, and writes it to eigen/api_key.json
// If apiKey != "", the custom key is used as the API key.
// If apiKey == "", a random API key is generated and overwrites any existing API key.
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
