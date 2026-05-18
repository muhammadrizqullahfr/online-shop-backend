package services

import (
	"context"

	"github.com/RakaMurdiarta/online-shop-system/internal/modules/products/delivery"
)

type CategoryService interface {
	Create(ctx context.Context, req *delivery.CreateCategoryRequest) (*delivery.CategoryResponse, error)
	GetByID(ctx context.Context, id uint) (*delivery.CategoryResponse, error)
	GetBySlug(ctx context.Context, slug string) (*delivery.CategoryResponse, error)
	GetAll(ctx context.Context) ([]delivery.CategoryResponse, error)
	GetTree(ctx context.Context) ([]delivery.CategoryResponse, error)
	Update(ctx context.Context, id uint, req *delivery.UpdateCategoryRequest) (*delivery.CategoryResponse, error)
	Delete(ctx context.Context, id uint) error
}
