package repository

import (
	"context"

	"github.com/RakaMurdiarta/online-shop-system/internal/models"
	"github.com/RakaMurdiarta/online-shop-system/pkg/database"
	"github.com/google/uuid"
)

type CartRepository interface {
	FindCartByUserID(ctx context.Context, userID uint) (*models.Cart, error)
	CreateCart(ctx context.Context, cart *models.Cart) error
	FindCartWithItems(ctx context.Context, cartID uuid.UUID) (*models.Cart, error)

	FindItemByMatch(ctx context.Context, cartID, productID uuid.UUID, variantsJSON string) (*models.CartItem, error)
	FindItemByID(ctx context.Context, id uuid.UUID) (*models.CartItem, error)
	CreateItem(ctx context.Context, item *models.CartItem) error
	UpdateItem(ctx context.Context, item *models.CartItem) error
	DeleteItem(ctx context.Context, id uuid.UUID) error
	DeleteItemsByCartID(ctx context.Context, cartID uuid.UUID) error

	database.TransactionManager
}
