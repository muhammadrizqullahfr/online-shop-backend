package impl

import (
	"context"

	"github.com/RakaMurdiarta/online-shop-system/internal/models"
	"github.com/RakaMurdiarta/online-shop-system/internal/modules/articles/repository"
	"github.com/RakaMurdiarta/online-shop-system/pkg/database"
)

type articleRepositoryImpl struct {
	*database.TransactionManagerImpl
}

func NewArticleRepository(tm *database.TransactionManagerImpl) repository.ArticleRepository {
	return &articleRepositoryImpl{TransactionManagerImpl: tm}
}

func (r *articleRepositoryImpl) Create(ctx context.Context, a *models.Article) error {
	return r.GetTx(ctx).Create(a).Error
}

func (r *articleRepositoryImpl) Update(ctx context.Context, a *models.Article) error {
	return r.GetTx(ctx).Save(a).Error
}

func (r *articleRepositoryImpl) FindByID(ctx context.Context, id string) (*models.Article, error) {
	var a models.Article
	if err := r.GetTx(ctx).Where("id = ?", id).First(&a).Error; err != nil {
		return nil, err
	}
	return &a, nil
}

func (r *articleRepositoryImpl) FindAll(ctx context.Context, limit, offset int, search, category string) ([]models.Article, int64, error) {
	var articles []models.Article
	var total int64

	query := r.GetTx(ctx).Model(&models.Article{})

	if search != "" {
		query = query.Where("title ILIKE ?", "%"+search+"%")
	}
	if category != "" {
		query = query.Where("category = ?", category)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := query.
		Order("published_at DESC NULLS LAST, created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&articles).Error
	return articles, total, err
}

func (r *articleRepositoryImpl) Delete(ctx context.Context, id string) error {
	return r.GetTx(ctx).Where("id = ?", id).Delete(&models.Article{}).Error
}
