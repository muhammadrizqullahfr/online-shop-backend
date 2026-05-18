package impl

import (
	"context"

	"github.com/RakaMurdiarta/online-shop-system/internal/models"
	"github.com/RakaMurdiarta/online-shop-system/internal/modules/products/repository"
	"github.com/RakaMurdiarta/online-shop-system/pkg/database"
)

type categoryRepositoryImpl struct {
	*database.TransactionManagerImpl
}

func NewCategoryRepository(db *database.TransactionManagerImpl) repository.CategoryRepository {
	return &categoryRepositoryImpl{
		TransactionManagerImpl: db,
	}
}

func (r *categoryRepositoryImpl) Create(ctx context.Context, category *models.Category) error {
	return r.GetTx(ctx).WithContext(ctx).Create(category).Error
}

func (r *categoryRepositoryImpl) GetByID(ctx context.Context, id uint) (*models.Category, error) {
	var category models.Category

	err := r.GetTx(ctx).WithContext(ctx).
		Preload("Children").
		First(&category, id).
		Error
	if err != nil {
		return nil, err
	}

	return &category, nil
}

func (r *categoryRepositoryImpl) GetBySlug(ctx context.Context, slug string) (*models.Category, error) {
	var category models.Category

	err := r.GetTx(ctx).WithContext(ctx).
		Where("slug = ?", slug).
		First(&category).
		Error
	if err != nil {
		return nil, err
	}

	return &category, nil
}

func (r *categoryRepositoryImpl) GetAll(ctx context.Context) ([]models.Category, error) {
	var categories []models.Category

	err := r.GetTx(ctx).WithContext(ctx).
		Order("id ASC").
		Find(&categories).
		Error

	return categories, err
}

func (r *categoryRepositoryImpl) GetRoots(ctx context.Context) ([]models.Category, error) {
	var categories []models.Category

	err := r.GetTx(ctx).WithContext(ctx).
		Where("parent_id IS NULL").
		Order("id ASC").
		Find(&categories).
		Error

	return categories, err
}

func (r *categoryRepositoryImpl) GetChildren(ctx context.Context, parentID uint) ([]models.Category, error) {
	var categories []models.Category

	err := r.GetTx(ctx).WithContext(ctx).
		Where("parent_id = ?", parentID).
		Order("id ASC").
		Find(&categories).
		Error

	return categories, err
}

func (r *categoryRepositoryImpl) Update(ctx context.Context, category *models.Category) error {
	return r.GetTx(ctx).WithContext(ctx).Save(category).Error
}

func (r *categoryRepositoryImpl) Delete(ctx context.Context, id uint) error {
	return r.GetTx(ctx).WithContext(ctx).Delete(&models.Category{}, id).Error
}

func (r *categoryRepositoryImpl) ExistsByName(ctx context.Context, name string) (bool, error) {
	var count int64
	err := r.GetTx(ctx).WithContext(ctx).
		Model(&models.Category{}).
		Where("name = ?", name).
		Count(&count).
		Error

	return count > 0, err
}

func (r *categoryRepositoryImpl) ExistsBySlug(ctx context.Context, slug string) (bool, error) {
	var count int64
	err := r.GetTx(ctx).WithContext(ctx).
		Model(&models.Category{}).
		Where("slug = ?", slug).
		Count(&count).
		Error

	return count > 0, err
}
