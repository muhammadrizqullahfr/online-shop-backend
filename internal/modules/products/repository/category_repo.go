package repository

import (
	"context"

	"github.com/RakaMurdiarta/online-shop-system/internal/models"
)

type CategoryRepository interface {
	Create(ctx context.Context, category *models.Category) error
	GetByID(ctx context.Context, id uint) (*models.Category, error)
	GetBySlug(ctx context.Context, slug string) (*models.Category, error)
	GetAll(ctx context.Context) ([]models.Category, error)
	GetRoots(ctx context.Context) ([]models.Category, error)
	GetChildren(ctx context.Context, parentID uint) ([]models.Category, error)
	Update(ctx context.Context, category *models.Category) error
	Delete(ctx context.Context, id uint) error
	ExistsByName(ctx context.Context, name string) (bool, error)
	ExistsBySlug(ctx context.Context, slug string) (bool, error)
}
