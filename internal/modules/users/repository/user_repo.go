package repository

import (
	"context"

	"github.com/RakaMurdiarta/online-shop-system/internal/models"
	"github.com/RakaMurdiarta/online-shop-system/pkg/common/constants"
	"github.com/RakaMurdiarta/online-shop-system/pkg/database"
)

type UserRepository interface {
	CreateUser(ctx context.Context, user *models.User) error
	GetUserByID(ctx context.Context, id uint) (*models.User, error)
	GetUserByEmail(ctx context.Context, email string) (*models.User, error)
	UpdateUser(ctx context.Context, user *models.User) error
	DeleteUser(ctx context.Context, id uint) error
	ListUsers(ctx context.Context, limit, offset int, search string, role constants.UserRole) ([]models.User, int64, error)
	IsSeller(ctx context.Context, userID uint) (bool, error)
	FindByProvider(ctx context.Context, provider constants.Provider, providerID string) (*models.User, error)
	FindBySession(ctx context.Context, token string) (*models.User, error)

	database.TransactionManager
}
