package config

import (
	"os"
)

type config struct {
	GrantType     string
	ClientId      string
	ClientSecret  string
	Username      string
	Password      string
	OAuthEndpoint string
}

func New() (*config, error) {

	var cfg config
	cfg.ClientId = os.Getenv("ClientId")
	cfg.ClientSecret = os.Getenv("ClientSecret")
	cfg.GrantType = os.Getenv("GrantType")
	cfg.Username = os.Getenv("Username")
	cfg.Password = os.Getenv("Password")
	cfg.OAuthEndpoint = os.Getenv("OAuthEndpoint")

	return &cfg, nil

}
