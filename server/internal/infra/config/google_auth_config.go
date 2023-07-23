package config

import (
	"fmt"
	"os"
)

type GoogleAuth struct {
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	RedirectURL  string `json:"redirect_url"`
}

func LoadGoogleAuthConfig() (GoogleAuth, error) {
	clientID := os.Getenv("GOOGLE_OAUTH_CLIENT_ID")
	if len(clientID) == 0 {
		return GoogleAuth{}, fmt.Errorf("missing GOOGLE_OAUTH_CLIENT_ID")
	}
	clientSecret := os.Getenv("GOOGLE_OAUTH_CLIENT_SECRET")
	if len(clientSecret) == 0 {
		return GoogleAuth{}, fmt.Errorf("missing GOOGLE_OAUTH_CLIENT_SECRET")
	}
	redirectURL := os.Getenv("GOOGLE_OAUTH_REDIRECT_URL")
	if len(redirectURL) == 0 {
		return GoogleAuth{}, fmt.Errorf("missing GOOGLE_OAUTH_REDIRECT_URL")
	}
	return GoogleAuth{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURL:  redirectURL,
	}, nil
}
