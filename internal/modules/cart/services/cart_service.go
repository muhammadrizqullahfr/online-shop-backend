package services

import (
	"context"

	"github.com/RakaMurdiarta/online-shop-system/internal/modules/cart/delivery"
	"github.com/google/uuid"
)

type CartService interface {
	GetCart(ctx context.Context, userID uint) (*delivery.CartResponse, error)
	AddItem(ctx context.Context, userID uint, req *delivery.AddItemRequest) (*delivery.CartResponse, error)
	UpdateItemQuantity(ctx context.Context, userID uint, itemID uuid.UUID, req *delivery.UpdateItemQuantityRequest) (*delivery.CartResponse, error)
	DeleteItem(ctx context.Context, userID uint, itemID uuid.UUID) error
	ClearCart(ctx context.Context, userID uint) error
}
