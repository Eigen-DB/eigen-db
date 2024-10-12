package auth

import (
	"errors"
	"fmt"
	"os"
	"testing"
)

// clean up this file if it still exists prior to running the tests
const TEST_API_KEY_FILEPATH string = "/tmp/api_key.txt"

func TestMain(m *testing.M) {
	if err := setup(); err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	code := m.Run()
	if err := cleanup(); err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	os.Exit(code)
}

func setup() error {
	if err := cleanup(); err != nil && !os.IsNotExist(err) { // ignore the file not existing
		return err
	}
	return nil
}

func cleanup() error {
	return os.Remove(TEST_API_KEY_FILEPATH)
}

// checks if the api_key.txt file exists after a test
func checkApiKeyFile(t *testing.T) {
	if _, err := os.Stat(TEST_API_KEY_FILEPATH); errors.Is(err, os.ErrNotExist) {
		t.Fatal("api_key.txt file not created")
	}
}

func TestSetupAPIKey_random(t *testing.T) {
	_, err := SetupAPIKey("", false, TEST_API_KEY_FILEPATH)
	if err != nil {
		t.Fatalf("Failed setting up API key: %s", err.Error())
	}

	checkApiKeyFile(t)
}

func TestSetupAPIKey_regen(t *testing.T) {
	key1, err := SetupAPIKey("", false, TEST_API_KEY_FILEPATH)
	if err != nil {
		t.Fatalf("Failed setting up API key: %s", err.Error())
	}

	key2, err := SetupAPIKey("", true, TEST_API_KEY_FILEPATH)
	if err != nil {
		t.Fatalf("Failed regenerating API key: %s", err.Error())
	}

	if key1 == key2 {
		t.Fatal("Key was not regenerated, but stayed identical to the inital key")
	}

	checkApiKeyFile(t)
}

func TestSetupAPIKey_custom(t *testing.T) {
	key, err := SetupAPIKey("12345", true, TEST_API_KEY_FILEPATH)
	if err != nil {
		t.Fatalf("Failed to generate custom API key: %s", err.Error())
	}

	if key != "12345" {
		t.Fatalf("API key has not taken the custom value. API key = %s", key)
	}

	checkApiKeyFile(t)
}

func TestSetupAPIKey_custom_to_regen(t *testing.T) {
	key1, err := SetupAPIKey("12345", true, TEST_API_KEY_FILEPATH)
	if err != nil {
		t.Fatalf("Failed to generate custom API key: %s", err.Error())
	}

	key2, err := SetupAPIKey("", true, TEST_API_KEY_FILEPATH)
	if err != nil {
		t.Fatalf("Failed regenerating API key: %s", err.Error())
	}

	if key1 == key2 {
		t.Fatal("Key was not regenerated, but stayed identical to the inital custom key")
	}

	checkApiKeyFile(t)
}
