package impl

import (
	"context"
	"errors"
	"fmt"
	"math"
	"time"

	"github.com/RakaMurdiarta/online-shop-system/internal/models"
	"github.com/RakaMurdiarta/online-shop-system/internal/modules/articles/delivery"
	"github.com/RakaMurdiarta/online-shop-system/internal/modules/articles/repository"
	"github.com/RakaMurdiarta/online-shop-system/internal/modules/articles/services"
	"github.com/RakaMurdiarta/online-shop-system/pkg/common/response"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type articleServiceImpl struct {
	repo repository.ArticleRepository
}

func NewArticleService(repo repository.ArticleRepository) services.ArticleService {
	return &articleServiceImpl{repo: repo}
}

func (s *articleServiceImpl) Create(ctx context.Context, req *delivery.CreateArticleRequest) (*delivery.ArticleResponse, error) {
	article := &models.Article{
		ID:          uuid.NewString(),
		Title:       req.Title,
		Description: optionalString(req.Description),
		Content:     optionalString(req.Content),
		Image:       optionalString(req.Image),
		Author:      optionalString(req.Author),
		Category:    optionalString(req.Category),
		PublishedAt: TimePtr(time.Now()),
		Featured:    req.Featured,
		ReadTime:    req.ReadTime,
		Excerpt:     optionalString(req.Excerpt),
	}

	if err := s.repo.WithTransaction(ctx, func(txCtx context.Context) error {
		return s.repo.Create(txCtx, article)
	}); err != nil {
		return nil, fmt.Errorf("failed to create article: %w", err)
	}

	return s.GetByID(ctx, article.ID)
}

func (s *articleServiceImpl) Update(ctx context.Context, id string, req *delivery.UpdateArticleRequest) (*delivery.ArticleResponse, error) {
	err := s.repo.WithTransaction(ctx, func(txCtx context.Context) error {
		article, err := s.repo.FindByID(txCtx, id)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return fmt.Errorf("article not found")
			}
			return err
		}

		article.Title = req.Title
		article.Description = optionalString(req.Description)
		article.Content = optionalString(req.Content)
		article.Image = optionalString(req.Image)
		article.Author = optionalString(req.Author)
		article.Category = optionalString(req.Category)
		article.PublishedAt = TimePtr(time.Now())
		article.Featured = req.Featured
		article.ReadTime = req.ReadTime
		article.Excerpt = optionalString(req.Excerpt)

		if err := s.repo.Update(txCtx, article); err != nil {
			return fmt.Errorf("failed to update article: %w", err)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return s.GetByID(ctx, id)
}

func (s *articleServiceImpl) GetByID(ctx context.Context, id string) (*delivery.ArticleResponse, error) {
	article, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("article not found")
		}
		return nil, err
	}
	return delivery.ToArticleResponse(article), nil
}

func (s *articleServiceImpl) GetAll(ctx context.Context, page, limit int, search, category string) (*delivery.ArticleListResponse, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	offset := (page - 1) * limit

	articles, total, err := s.repo.FindAll(ctx, limit, offset, search, category)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch articles: %w", err)
	}

	items := make([]delivery.ArticleResponse, 0, len(articles))
	for i := range articles {
		items = append(items, *delivery.ToArticleResponse(&articles[i]))
	}

	totalPage := 0
	if total > 0 {
		totalPage = int(math.Ceil(float64(total) / float64(limit)))
	}

	return &delivery.ArticleListResponse{
		Items: items,
		Paging: response.PagingResponse{
			CurrentPage: page,
			TotalPage:   totalPage,
			TotalItems:  total,
		},
	}, nil
}

func (s *articleServiceImpl) Delete(ctx context.Context, id string) error {
	if err := s.repo.Delete(ctx, id); err != nil {
		return fmt.Errorf("failed to delete article: %w", err)
	}
	return nil
}

func optionalString(v string) *string {
	if v == "" {
		return nil
	}
	return &v
}

func TimePtr(t time.Time) *time.Time {
	return &t
}
