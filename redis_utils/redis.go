package redis_utils

import (
	"context"
	"crypto/rand"
	"eigen_db/constants"
	"encoding/hex"
	"fmt"
	"os"

	"github.com/redis/go-redis/v9"
)

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

func CheckConnection(ctx context.Context, client *redis.Client) error {
	if _, err := client.Ping(ctx).Result(); err != nil {
		return err
	}
	return nil
}

func SetupAPIKey(ctx context.Context, client *redis.Client, apiKey string) (string, error) {
	if err := CheckConnection(ctx, client); err != nil {
		return "", err
	}

	if _, err := client.Get(ctx, constants.REDIS_API_KEY_NAME).Result(); err != nil { // key does not exist
		keyBytes := make([]byte, 32)
		if _, err := rand.Read(keyBytes); err != nil {
			return "", err
		}
		apiKey = hex.EncodeToString(keyBytes)
	}

	if apiKey != "" {
		if status := client.Set(ctx, "apiKey", apiKey, 0); status.Err() != nil {
			return "", status.Err()
		}
	}

	val, err := client.Get(ctx, constants.REDIS_API_KEY_NAME).Result()
	if err != nil {
		return "", err
	}

	return val, nil
}
