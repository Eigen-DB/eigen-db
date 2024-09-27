package redis_utils

import (
	"context"
	"testing"

	"github.com/alicebob/miniredis"
	"github.com/redis/go-redis/v9"
)

func setupTestRedisServer(t *testing.T) *redis.Client {
	server, err := miniredis.Run()
	if err != nil {
		t.Fatalf("Error starting miniredis server: %s", err.Error())
	}
	return redis.NewClient(&redis.Options{
		Addr: server.Addr(),
	})
}

func TestSetupAPIKey_random(t *testing.T) {
	client := setupTestRedisServer(t)
	defer client.Close()

	ctx := context.Background()
	_, err := SetupAPIKey(ctx, client, "")
	if err != nil {
		t.Fatalf("Failed setting up API key: %s", err.Error())
	}
}

func TestSetupAPIKey_custom(t *testing.T) {
	client := setupTestRedisServer(t)
	defer client.Close()

	ctx := context.Background()
	key, err := SetupAPIKey(ctx, client, "test123")
	if err != nil {
		t.Fatalf("Failed setting up API key: %s", err.Error())
	}

	if key != "test123" {
		t.Fatalf("API key not set to custom value")
	}
}

func TestSetupAPIKey_custom_to_random(t *testing.T) {
	client := setupTestRedisServer(t)
	defer client.Close()

	ctx := context.Background()
	key, err := SetupAPIKey(ctx, client, "test123")
	if err != nil {
		t.Fatalf("Failed setting up API key: %s", err.Error())
	}

	if key != "test123" {
		t.Fatalf("API key not set to custom value")
	}

	key, err = SetupAPIKey(ctx, client, "")
	if err != nil {
		t.Fatalf("Failed setting up API key: %s", err.Error())
	}

	if key == "test123" {
		t.Fatalf("API key did not go back to being random")
	}
}

func TestSetupAPIKey_bad_redis_connection(t *testing.T) {
	badClient := redis.NewClient(&redis.Options{
		Addr: "0.0.0.0:0", // invalid server address
	})
	defer badClient.Close()
	ctx := context.Background()

	_, err := SetupAPIKey(ctx, badClient, "")
	if err == nil {
		t.Fatalf("No error occured with an invalid Redis server address")
	}
}
