package impl

import (
	"context"

	"github.com/RakaMurdiarta/online-shop-system/internal/models"
	"github.com/RakaMurdiarta/online-shop-system/internal/modules/products/repository"
	"github.com/RakaMurdiarta/online-shop-system/pkg/database"
	"github.com/google/uuid"
)

type newProductRepositoryImpl struct {
	*database.TransactionManagerImpl
}

func NewProductRepository(tm *database.TransactionManagerImpl) repository.ProductRepository {
	return &newProductRepositoryImpl{TransactionManagerImpl: tm}
}

func (r *newProductRepositoryImpl) Create(ctx context.Context, product *models.Product) error {
	return r.GetTx(ctx).Create(product).Error
}

func (r *newProductRepositoryImpl) CreateVariants(ctx context.Context, variants []models.ProductVariant) error {
	if len(variants) == 0 {
		return nil
	}
	return r.GetTx(ctx).Create(&variants).Error
}

func (r *newProductRepositoryImpl) FindByID(ctx context.Context, id uuid.UUID) (*models.Product, error) {
	var product models.Product
	err := r.GetTx(ctx).
		Preload("Category").
		Preload("Variants").
		Where("id = ?", id).
		First(&product).Error
	if err != nil {
		return nil, err
	}
	return &product, nil
}

func (r *newProductRepositoryImpl) FindAll(ctx context.Context, limit, offset int, search, category string) ([]models.Product, int64, error) {
	var (
		products []models.Product
		total    int64
	)

	db := r.GetTx(ctx).Model(&models.Product{})

	if search != "" {
		db = db.Where("name ILIKE ?", "%"+search+"%")
	}
	if category != "" {
		db = db.Joins("JOIN categories ON categories.id = products.category_id").
			Where("categories.slug = ?", category)
	}

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := db.
		Preload("Category").
		Preload("Variants").
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&products).Error
	if err != nil {
		return nil, 0, err
	}

	return products, total, nil
}

func (r *newProductRepositoryImpl) Update(ctx context.Context, product *models.Product) error {
	return r.GetTx(ctx).Save(product).Error
}

func (r *newProductRepositoryImpl) DeleteVariantsByProductID(ctx context.Context, productID uuid.UUID) error {
	return r.GetTx(ctx).Where("product_id = ?", productID).Delete(&models.ProductVariant{}).Error
}

func (r *newProductRepositoryImpl) Delete(ctx context.Context, id uuid.UUID) error {
	return r.GetTx(ctx).Where("id = ?", id).Delete(&models.Product{}).Error
}
