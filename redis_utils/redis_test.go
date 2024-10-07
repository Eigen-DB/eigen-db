package redis_utils

import (
	"context"
	"errors"
	"os"
	"testing"

	"github.com/alicebob/miniredis"
	"github.com/redis/go-redis/v9"
)

const TEST_API_KEY_FILEPATH string = "/tmp/api_key.txt"

func setupTestRedisServer(t *testing.T) *redis.Client {
	server, err := miniredis.Run()
	if err != nil {
		t.Fatalf("Error starting miniredis server: %s", err.Error())
	}
	return redis.NewClient(&redis.Options{
		Addr: server.Addr(),
	})
}

func cleanup(t *testing.T) {
	if err := os.Remove(TEST_API_KEY_FILEPATH); err != nil {
		t.Fatalf("Error cleaning up: %s", err.Error())
	}
}

// checks if the api_key.txt file exists after a test
func checkApiKeyFile(t *testing.T) {
	if _, err := os.Stat(TEST_API_KEY_FILEPATH); errors.Is(err, os.ErrNotExist) {
		t.Fatal("api_key.txt file not created")
	}
}

func TestSetupAPIKey_random(t *testing.T) {
	client := setupTestRedisServer(t)
	defer client.Close()

	ctx := context.Background()
	_, err := SetupAPIKey(ctx, client, "", TEST_API_KEY_FILEPATH)
	if err != nil {
		t.Fatalf("Failed setting up API key: %s", err.Error())
	}

	checkApiKeyFile(t)
	cleanup(t)
}

func TestSetupAPIKey_custom(t *testing.T) {
	client := setupTestRedisServer(t)
	defer client.Close()

	ctx := context.Background()
	key, err := SetupAPIKey(ctx, client, "test123", TEST_API_KEY_FILEPATH)
	if err != nil {
		t.Fatalf("Failed setting up API key: %s", err.Error())
	}

	if key != "test123" {
		t.Fatalf("API key not set to custom value")
	}

	checkApiKeyFile(t)
	cleanup(t)
}

func TestSetupAPIKey_custom_to_random(t *testing.T) {
	client := setupTestRedisServer(t)
	defer client.Close()

	ctx := context.Background()
	key, err := SetupAPIKey(ctx, client, "test123", TEST_API_KEY_FILEPATH)
	if err != nil {
		t.Fatalf("Failed setting up API key: %s", err.Error())
	}

	if key != "test123" {
		t.Fatalf("API key not set to custom value")
	}

	key, err = SetupAPIKey(ctx, client, "", TEST_API_KEY_FILEPATH)
	if err != nil {
		t.Fatalf("Failed setting up API key: %s", err.Error())
	}

	if key == "test123" {
		t.Fatalf("API key did not go back to being random")
	}

	checkApiKeyFile(t)
	cleanup(t)
}

func TestSetupAPIKey_bad_redis_connection(t *testing.T) {
	badClient := redis.NewClient(&redis.Options{
		Addr: "x.x.x.x:0", // invalid server address
	})
	defer badClient.Close()
	ctx := context.Background()

	_, err := SetupAPIKey(ctx, badClient, "", TEST_API_KEY_FILEPATH)
	if err == nil {
		t.Fatalf("No error occured with an invalid Redis server address")
	}
}
