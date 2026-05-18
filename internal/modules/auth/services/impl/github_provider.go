package impl

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/RakaMurdiarta/online-shop-system/internal/config"
	"github.com/RakaMurdiarta/online-shop-system/internal/modules/auth/services"
	"github.com/RakaMurdiarta/online-shop-system/internal/modules/users/repository"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
)

type GithubProvider struct {
	userRepo repository.UserRepository
	cfg      *config.Config
}

func NewGithubProvider(userRepo repository.UserRepository, cfg *config.Config) services.OauthService {
	return &GithubProvider{
		userRepo: userRepo,
		cfg:      cfg,
	}
}

func (g *GithubProvider) Config() *oauth2.Config {

	return &oauth2.Config{
		ClientID:     g.cfg.GithubClientID,
		ClientSecret: g.cfg.GithubSecret,
		Endpoint:     github.Endpoint,
		Scopes:       []string{"read:user", "user:email"},
		RedirectURL:  g.cfg.GithubCallbackUrl,
	}
}

func (g *GithubProvider) Callback(ctx context.Context, code string) (*services.UserProviderResponse, error) {
	token, err := g.Config().Exchange(ctx, code)
	if err != nil {
		return nil, fmt.Errorf("exchange token failed: %w", err)
	}

	client := g.Config().Client(ctx, token)

	resp, err := client.Get("https://api.github.com/user")
	if err != nil {
		return nil, fmt.Errorf("failed get github user: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("github user api error: %s", resp.Status)
	}

	var ghUser struct {
		ID        int    `json:"id"`
		Login     string `json:"login"`
		Name      string `json:"name"`
		AvatarURL string `json:"avatar_url"`
		Email     string `json:"email"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&ghUser); err != nil {
		return nil, fmt.Errorf("decode user failed: %w", err)
	}

	if ghUser.Email == "" {
		respEmail, err := client.Get("https://api.github.com/user/emails")
		if err != nil {
			return nil, fmt.Errorf("failed get emails: %w", err)
		}
		defer respEmail.Body.Close()

		if respEmail.StatusCode != 200 {
			return nil, fmt.Errorf("github email api error: %s", respEmail.Status)
		}

		var emails []struct {
			Email    string `json:"email"`
			Primary  bool   `json:"primary"`
			Verified bool   `json:"verified"`
		}

		if err := json.NewDecoder(respEmail.Body).Decode(&emails); err != nil {
			return nil, fmt.Errorf("decode emails failed: %w", err)
		}

		for _, e := range emails {
			if e.Primary && e.Verified {
				ghUser.Email = e.Email
				break
			}
		}
	}

	gUser := &services.UserProviderResponse{
		Sub:     fmt.Sprintf("%d", ghUser.ID),
		Email:   ghUser.Email,
		Name:    ghUser.Name,
		Picture: ghUser.AvatarURL,
	}

	return gUser, nil
}

func (g *GithubProvider) LoginRedirectUrl(state string) string {
	return g.Config().AuthCodeURL(state)
}
