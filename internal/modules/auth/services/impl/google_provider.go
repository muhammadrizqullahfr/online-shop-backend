package impl

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/RakaMurdiarta/online-shop-system/internal/config"
	"github.com/RakaMurdiarta/online-shop-system/internal/modules/auth/services"
	"github.com/RakaMurdiarta/online-shop-system/internal/modules/users/repository"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type GoogleProvider struct {
	userRepo repository.UserRepository
	cfg      *config.Config
}

func NewGoogleProvider(userRepo repository.UserRepository, cfg *config.Config) services.OauthService {
	return &GoogleProvider{
		userRepo: userRepo,
		cfg:      cfg,
	}
}

func (g *GoogleProvider) Config() *oauth2.Config {

	return &oauth2.Config{
		ClientID:     g.cfg.GoogleClientID,
		ClientSecret: g.cfg.GoogleSecret,
		Endpoint:     google.Endpoint,
		Scopes:       []string{"email", "profile"},
		RedirectURL:  g.cfg.GoogleCallbackUrl,
	}
}

func (g *GoogleProvider) Callback(ctx context.Context, code string) (*services.UserProviderResponse, error) {
	token, err := g.Config().Exchange(ctx, code)
	if err != nil {
		return nil, fmt.Errorf("exchange token failed: %w", err)
	}

	client := g.Config().Client(ctx, token)

	resp, err := client.Get("https://www.googleapis.com/oauth2/v3/userinfo")
	if err != nil {
		return nil, fmt.Errorf("failed get google user: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("google user api error: %s", resp.Status)
	}

	var googleUser struct {
		Sub           string `json:"sub"`
		Email         string `json:"email"`
		Name          string `json:"name"`
		Picture       string `json:"picture"`
		EmailVerified bool   `json:"email_verified"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&googleUser); err != nil {
		return nil, fmt.Errorf("decode google user failed: %w", err)
	}

	if !googleUser.EmailVerified {
		return nil, fmt.Errorf("email not verified")
	}

	gUser := &services.UserProviderResponse{
		Sub:     googleUser.Sub,
		Email:   googleUser.Email,
		Name:    googleUser.Name,
		Picture: googleUser.Picture,
	}

	return gUser, nil
}

func (g *GoogleProvider) LoginRedirectUrl(state string) string {
	return g.Config().AuthCodeURL(state)
}
