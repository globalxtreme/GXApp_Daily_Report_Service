package config

import "os"

type OAuth2Config struct {
	BaseURL      string
	ClientID     string
	ClientSecret string
	RedirectURL  string
	FrontendURL  string
}

var OAuth2 *OAuth2Config

func InitOAuth2() {
	OAuth2 = &OAuth2Config{
		BaseURL:      os.Getenv("GX_OAUTH_BASE_URL"),
		ClientID:     os.Getenv("GX_OAUTH_CLIENT_ID"),
		ClientSecret: os.Getenv("GX_OAUTH_CLIENT_SECRET"),
		RedirectURL:  os.Getenv("GX_OAUTH_REDIRECT_URL"),
		FrontendURL:  os.Getenv("FRONTEND_URL"),
	}
}
