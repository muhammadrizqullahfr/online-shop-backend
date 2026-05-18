package services

import (
	"context"

	"golang.org/x/oauth2"
)

type UserProviderResponse struct {
	Sub     string `json:"sub"`
	Email   string `json:"email"`
	Name    string `json:"name"`
	Picture string `json:"picture"`
}

type OauthService interface {
	Callback(ctx context.Context, code string) (*UserProviderResponse, error)
	Config() *oauth2.Config
	LoginRedirectUrl(state string) string
}
