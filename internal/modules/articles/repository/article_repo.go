package repository

import (
	"context"

	"github.com/RakaMurdiarta/online-shop-system/internal/models"
	"github.com/RakaMurdiarta/online-shop-system/pkg/database"
)

type ArticleRepository interface {
	Create(ctx context.Context, a *models.Article) error
	Update(ctx context.Context, a *models.Article) error
	FindByID(ctx context.Context, id string) (*models.Article, error)
	FindAll(ctx context.Context, limit, offset int, search, category string) ([]models.Article, int64, error)
	Delete(ctx context.Context, id string) error

	//for database transaction
	database.TransactionManager
}
