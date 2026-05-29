package impl

import (
	"context"
	"errors"
	"fmt"
	"math"

	"github.com/RakaMurdiarta/online-shop-system/internal/models"
	categoryRepo "github.com/RakaMurdiarta/online-shop-system/internal/modules/products/repository"

	"github.com/RakaMurdiarta/online-shop-system/internal/modules/products/delivery"
	"github.com/RakaMurdiarta/online-shop-system/internal/modules/products/repository"
	"github.com/RakaMurdiarta/online-shop-system/internal/modules/products/services"
	"github.com/RakaMurdiarta/online-shop-system/internal/modules/recommendation/client"
	"github.com/RakaMurdiarta/online-shop-system/pkg/common/response"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type newProductServiceImpl struct {
	repo         repository.ProductRepository
	categoryRepo categoryRepo.CategoryRepository
	recClient    client.RecommendationClient
}

func NewProductService(repo repository.ProductRepository, catRepo categoryRepo.CategoryRepository, recClient client.RecommendationClient) services.ProductService {
	return &newProductServiceImpl{repo: repo, categoryRepo: catRepo, recClient: recClient}
}

func (s *newProductServiceImpl) resolveCategoryID(ctx context.Context, slug string) (uint, error) {
	cat, err := s.categoryRepo.GetBySlug(ctx, slug)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return 0, fmt.Errorf("category with slug %q not found", slug)
		}
		return 0, err
	}
	return cat.ID, nil
}

func (s *newProductServiceImpl) Create(ctx context.Context, req *delivery.CreateNewProductRequest) (*delivery.NewProductResponse, error) {
	categoryID, err := s.resolveCategoryID(ctx, req.CategorySlug)
	if err != nil {
		return nil, err
	}

	product := &models.Product{
		Name:          req.Name,
		Description:   req.Description,
		CategoryID:    categoryID,
		Price:         req.Price,
		OriginalPrice: req.OriginalPrice,
		Image:         req.Image,
		Rating:        req.Rating,
		Reviews:       req.Reviews,
		Stock:         req.Stock,
		InStock:       req.InStock,
		Featured:      req.Featured,
		Tags:          models.StringSlice(req.Tags),
	}

	err = s.repo.WithTransaction(ctx, func(txCtx context.Context) error {
		if err := s.repo.Create(txCtx, product); err != nil {
			return fmt.Errorf("failed to create product: %w", err)
		}

		if len(req.Variants) == 0 {
			return nil
		}

		variants := make([]models.ProductVariant, 0, len(req.Variants))
		for _, v := range req.Variants {
			variants = append(variants, models.ProductVariant{
				ProductID:       product.ID,
				Name:            v.Name,
				Type:            v.Type,
				Value:           v.Value,
				Price:           v.Price,
				PriceAdjustment: v.PriceAdjustment,
				Stock:           v.Stock,
			})
		}
		if err := s.repo.CreateVariants(txCtx, variants); err != nil {
			return fmt.Errorf("failed to create variants: %w", err)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	if s.recClient != nil {
		go func(p models.Product) {
			_ = s.recClient.SyncProductToFastAPI(p)
		}(*product)
	}

	return s.GetByID(ctx, product.ID)
}

func (s *newProductServiceImpl) Update(ctx context.Context, id uuid.UUID, req *delivery.UpdateNewProductRequest) (*delivery.NewProductResponse, error) {
	categoryID, err := s.resolveCategoryID(ctx, req.CategorySlug)
	if err != nil {
		return nil, err
	}

	err = s.repo.WithTransaction(ctx, func(txCtx context.Context) error {
		product, err := s.repo.FindByID(txCtx, id)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return fmt.Errorf("product not found")
			}
			return err
		}

		product.Name = req.Name
		product.Description = req.Description
		product.CategoryID = categoryID
		product.Price = req.Price
		product.OriginalPrice = req.OriginalPrice
		product.Image = req.Image
		product.Rating = req.Rating
		product.Reviews = req.Reviews
		product.Stock = req.Stock
		product.InStock = req.InStock
		product.Featured = req.Featured
		product.Tags = models.StringSlice(req.Tags)

		if err := s.repo.Update(txCtx, product); err != nil {
			return fmt.Errorf("failed to update product: %w", err)
		}

		if err := s.repo.DeleteVariantsByProductID(txCtx, product.ID); err != nil {
			return fmt.Errorf("failed to clear variants: %w", err)
		}

		if len(req.Variants) == 0 {
			return nil
		}

		variants := make([]models.ProductVariant, 0, len(req.Variants))
		for _, v := range req.Variants {
			variants = append(variants, models.ProductVariant{
				ProductID:       product.ID,
				Name:            v.Name,
				Type:            v.Type,
				Value:           v.Value,
				Price:           v.Price,
				PriceAdjustment: v.PriceAdjustment,
				Stock:           v.Stock,
			})
		}
		if err := s.repo.CreateVariants(txCtx, variants); err != nil {
			return fmt.Errorf("failed to recreate variants: %w", err)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return s.GetByID(ctx, id)
}

func (s *newProductServiceImpl) GetByID(ctx context.Context, id uuid.UUID) (*delivery.NewProductResponse, error) {
	product, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("product not found")
		}
		return nil, err
	}
	return delivery.ToNewProductResponse(product), nil
}

func (s *newProductServiceImpl) GetAll(ctx context.Context, page, limit int, search, categorySlug string) (*delivery.NewProductListResponse, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	offset := (page - 1) * limit

	products, total, err := s.repo.FindAll(ctx, limit, offset, search, categorySlug)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch products: %w", err)
	}

	items := make([]delivery.NewProductResponse, 0, len(products))
	for i := range products {
		items = append(items, *delivery.ToNewProductResponse(&products[i]))
	}

	totalPage := 0
	if total > 0 {
		totalPage = int(math.Ceil(float64(total) / float64(limit)))
	}

	return &delivery.NewProductListResponse{
		Items: items,
		Paging: response.PagingResponse{
			CurrentPage: page,
			TotalPage:   totalPage,
			TotalItems:  total,
		},
	}, nil
}

func (s *newProductServiceImpl) Delete(ctx context.Context, id uuid.UUID) error {
	if err := s.repo.Delete(ctx, id); err != nil {
		return fmt.Errorf("failed to delete product: %w", err)
	}
	return nil
}
