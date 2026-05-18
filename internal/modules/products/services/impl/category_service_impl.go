package impl

import (
	"context"
	"errors"
	"fmt"

	"github.com/RakaMurdiarta/online-shop-system/internal/config"
	"github.com/RakaMurdiarta/online-shop-system/internal/models"
	"github.com/RakaMurdiarta/online-shop-system/internal/modules/products/delivery"
	"github.com/RakaMurdiarta/online-shop-system/internal/modules/products/repository"
	"github.com/RakaMurdiarta/online-shop-system/internal/modules/products/services"
	"github.com/RakaMurdiarta/online-shop-system/pkg/database"
	"gorm.io/gorm"
)

type categoryServiceImpl struct {
	txManager    database.TransactionManager
	conf         *config.Config
	categoryRepo repository.CategoryRepository
}

func NewCategoryService(categoryRepo repository.CategoryRepository, tx database.TransactionManager, conf *config.Config,
) services.CategoryService {

	return &categoryServiceImpl{
		categoryRepo: categoryRepo,
		txManager:    tx,
		conf:         conf,
	}
}

func (s *categoryServiceImpl) Create(ctx context.Context, req *delivery.CreateCategoryRequest) (*delivery.CategoryResponse, error) {

	exists, err := s.categoryRepo.ExistsBySlug(ctx, req.Slug)
	if err != nil {
		return nil, err
	}

	if exists {
		return nil, fmt.Errorf("category label already exists")
	}

	if req.ParentID != nil {
		_, err := s.categoryRepo.GetByID(ctx, *req.ParentID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, fmt.Errorf("parent category not found")
			}
			return nil, err
		}
	}

	category := &models.Category{
		Name:     req.Name,
		Slug:     req.Slug,
		ParentID: req.ParentID,
	}

	if err := s.categoryRepo.Create(ctx, category); err != nil {
		return nil, err
	}

	return delivery.ToCategoryResponse(category), nil
}

func (s *categoryServiceImpl) GetByID(ctx context.Context, id uint) (*delivery.CategoryResponse, error) {
	category, err := s.categoryRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return delivery.ToCategoryResponse(category), nil
}

func (s *categoryServiceImpl) GetBySlug(ctx context.Context, slug string) (*delivery.CategoryResponse, error) {
	category, err := s.categoryRepo.GetBySlug(ctx, slug)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("category not found")
		}
		return nil, err
	}

	return delivery.ToCategoryResponse(category), nil
}

func (s *categoryServiceImpl) GetAll(ctx context.Context) ([]delivery.CategoryResponse, error) {
	categories, err := s.categoryRepo.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	result := make([]delivery.CategoryResponse, 0, len(categories))
	for _, category := range categories {
		cat := category
		result = append(result, *delivery.ToCategoryResponse(&cat))
	}

	return result, nil
}

func (s *categoryServiceImpl) GetTree(ctx context.Context) ([]delivery.CategoryResponse, error) {
	categories, err := s.categoryRepo.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	return delivery.BuildCategoryTree(categories), nil
}

func (s *categoryServiceImpl) Update(ctx context.Context, id uint, req *delivery.UpdateCategoryRequest) (*delivery.CategoryResponse, error) {
	category, err := s.categoryRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if req.ParentID != nil {
		if *req.ParentID == id {
			return nil, fmt.Errorf("category cannot be its own parent")
		}

		_, err := s.categoryRepo.GetByID(ctx, *req.ParentID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, fmt.Errorf("parent category not found")
			}
			return nil, err
		}
	}

	if category.Slug != req.Slug {
		exists, err := s.categoryRepo.ExistsBySlug(ctx, req.Slug)
		if err != nil {
			return nil, err
		}
		if exists {
			return nil, fmt.Errorf("category slug already exists")
		}
	}

	category.Name = req.Name
	category.Slug = req.Slug
	category.ParentID = req.ParentID

	if err := s.categoryRepo.Update(ctx, category); err != nil {
		return nil, err
	}

	return delivery.ToCategoryResponse(category), nil
}

func (s *categoryServiceImpl) Delete(ctx context.Context, id uint) error {
	_, err := s.categoryRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	children, err := s.categoryRepo.GetChildren(ctx, id)
	if err != nil {
		return err
	}

	if len(children) > 0 {
		return fmt.Errorf("cannot delete category that still has child categories")
	}

	return s.categoryRepo.Delete(ctx, id)
}
