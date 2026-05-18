package services

import (
	"context"

	"github.com/RakaMurdiarta/online-shop-system/internal/models"
	"github.com/RakaMurdiarta/online-shop-system/internal/modules/auth/delivery"
)

type AuthService interface {
	RegisterLocal(ctx context.Context, req *delivery.SingUpRequest) (*delivery.RegisterResponse, error)
	LoginLocal(ctx context.Context, req *delivery.LoginRequest) (*delivery.LoginResponse, error)
	HandleOAuthCallback(ctx context.Context, data *delivery.OAuthUserRequest) (*models.User, error)
	RefreshToken(ctx context.Context, refreshToken string) (*delivery.RefreshTokenResponse, error)
	Me(ctx context.Context, userID uint) (*delivery.MeResponse, error)
}
