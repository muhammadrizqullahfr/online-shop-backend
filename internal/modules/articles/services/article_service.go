package services

import (
	"context"

	"github.com/RakaMurdiarta/online-shop-system/internal/modules/articles/delivery"
)

type ArticleService interface {
	Create(ctx context.Context, req *delivery.CreateArticleRequest) (*delivery.ArticleResponse, error)
	Update(ctx context.Context, id string, req *delivery.UpdateArticleRequest) (*delivery.ArticleResponse, error)
	GetByID(ctx context.Context, id string) (*delivery.ArticleResponse, error)
	GetAll(ctx context.Context, page, limit int, search, category string) (*delivery.ArticleListResponse, error)
	Delete(ctx context.Context, id string) error
}
