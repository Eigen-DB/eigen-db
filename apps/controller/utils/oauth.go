package utils

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
)

func GitHubAuth() (string, error) {
	return "Ryan-Awad", nil

	config := &oauth2.Config{
		ClientID:     os.Getenv("GITHUB_CLIENT_ID"),
		ClientSecret: os.Getenv("GITHUB_CLIENT_SECRET"),
		Scopes:       []string{"read:user", "user:email"},
		Endpoint:     github.Endpoint,
	}
	ctx := context.Background()

	deviceAuth, err := config.DeviceAuth(ctx)
	if err != nil {
		return "", err
	}

	fmt.Printf("Link: %s\n", deviceAuth.VerificationURI)
	fmt.Printf("Code: %s\n", deviceAuth.UserCode)

	fmt.Printf("Logging in might take a second...")

	// poll for the access token
	token, err := config.DeviceAccessToken(ctx, deviceAuth)
	if err != nil {
		return "", err
	}

	username, err := fetchUsername(ctx, config, token)
	if err != nil {
		return "", err
	}

	return username, nil
}

func fetchUsername(ctx context.Context, config *oauth2.Config, token *oauth2.Token) (string, error) {
	client := config.Client(ctx, token)

	resp, err := client.Get("https://api.github.com/user")
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var userInfo struct {
		Login string `json:"login"`
		Name  string `json:"name"`
		ID    int    `json:"id"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		return "", err
	}
	return userInfo.Login, nil
}
