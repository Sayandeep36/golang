package oauth

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

	"Golang/common"
	"Golang/internal/config"

	"github.com/joho/godotenv"
)

const (
	loginEndpoint    = "/services/oauth2/token"
	userInfoEndpoint = "/services/oauth2/userinfo"
)

type LoginResponse struct {
	AccessToken string `json:"access_token"`
	InstanceURL string `json:"instance_url"`
	ID          string `json:"id"`
	TokenType   string `json:"token_type"`
	IssuedAt    string `json:"issued_at"`
	Signature   string `json:"signature"`
}

type UserInfoResponse struct {
	UserID         string `json:"user_id"`
	OrganizationID string `json:"organization_id"`
}

const dev = "dev"

func Login() (*LoginResponse, error) {

	var (
		envType = flag.String("env_type", "dev", "Set environment type")
	)
	flag.Parse()

	if strings.EqualFold(*envType, dev) {
		err := godotenv.Load()
		if err != nil {
			log.Println("no .env file found")
		}
	}

	cfg, err := config.New()
	if err != nil {
		log.Fatalf("error reading config ")
	}

	body := url.Values{}
	body.Set("grant_type", cfg.GrantType)
	body.Set("client_id", cfg.ClientId)
	body.Set("client_secret", cfg.ClientSecret)
	body.Set("username", cfg.Username)
	body.Set("password", cfg.Password)

	ctx, cancelFn := context.WithTimeout(context.Background(), common.OAuthDialTimeout)
	defer cancelFn()

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, cfg.OAuthEndpoint+loginEndpoint, strings.NewReader(body.Encode()))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	httpResp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer httpResp.Body.Close()

	if httpResp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("non-200 status code returned on OAuth authentication call: %v", httpResp.StatusCode)
	}

	var loginResponse LoginResponse
	err = json.NewDecoder(httpResp.Body).Decode(&loginResponse)
	if err != nil {
		return nil, err
	}

	return &loginResponse, nil
}

func UserInfo(accessToken string) (*UserInfoResponse, error) {
	ctx, cancelFn := context.WithTimeout(context.Background(), common.OAuthDialTimeout)
	defer cancelFn()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, os.Getenv("OAuthEndpoint")+userInfoEndpoint, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken))

	httpResp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer httpResp.Body.Close()

	if httpResp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("non-200 status code returned on OAuth user info call: %v", httpResp.StatusCode)
	}

	var userInfoResponse UserInfoResponse
	err = json.NewDecoder(httpResp.Body).Decode(&userInfoResponse)
	if err != nil {
		return nil, err
	}

	return &userInfoResponse, nil
}
