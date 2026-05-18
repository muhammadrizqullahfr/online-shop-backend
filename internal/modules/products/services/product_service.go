package services

import (
	"context"

	"github.com/RakaMurdiarta/online-shop-system/internal/modules/products/delivery"
	"github.com/google/uuid"
)

type ProductService interface {
	Create(ctx context.Context, req *delivery.CreateNewProductRequest) (*delivery.NewProductResponse, error)
	Update(ctx context.Context, id uuid.UUID, req *delivery.UpdateNewProductRequest) (*delivery.NewProductResponse, error)
	GetByID(ctx context.Context, id uuid.UUID) (*delivery.NewProductResponse, error)
	GetAll(ctx context.Context, page, limit int, search, categorySlug string) (*delivery.NewProductListResponse, error)
	Delete(ctx context.Context, id uuid.UUID) error
}
