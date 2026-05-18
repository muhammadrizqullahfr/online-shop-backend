package repository

import (
	"context"

	"github.com/RakaMurdiarta/online-shop-system/internal/models"
	"github.com/RakaMurdiarta/online-shop-system/pkg/database"
	"github.com/google/uuid"
)

type ProductRepository interface {
	Create(ctx context.Context, product *models.Product) error
	CreateVariants(ctx context.Context, variants []models.ProductVariant) error

	FindByID(ctx context.Context, id uuid.UUID) (*models.Product, error)
	FindAll(ctx context.Context, limit, offset int, search, category string) ([]models.Product, int64, error)

	Update(ctx context.Context, product *models.Product) error
	DeleteVariantsByProductID(ctx context.Context, productID uuid.UUID) error
	Delete(ctx context.Context, id uuid.UUID) error

	database.TransactionManager
}
