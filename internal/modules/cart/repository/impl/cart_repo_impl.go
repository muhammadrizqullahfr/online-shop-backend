package impl

import (
	"context"

	"github.com/RakaMurdiarta/online-shop-system/internal/models"
	"github.com/RakaMurdiarta/online-shop-system/internal/modules/cart/repository"
	"github.com/RakaMurdiarta/online-shop-system/pkg/database"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type newCartRepositoryImpl struct {
	*database.TransactionManagerImpl
}

func NewCartRepository(tm *database.TransactionManagerImpl) repository.CartRepository {
	return &newCartRepositoryImpl{TransactionManagerImpl: tm}
}

func (r *newCartRepositoryImpl) FindCartByUserID(ctx context.Context, userID uint) (*models.Cart, error) {
	var cart models.Cart
	err := r.GetTx(ctx).
		Where("user_id = ?", userID).
		First(&cart).Error
	if err != nil {
		return nil, err
	}
	return &cart, nil
}

func (r *newCartRepositoryImpl) CreateCart(ctx context.Context, cart *models.Cart) error {
	return r.GetTx(ctx).Create(cart).Error
}

func (r *newCartRepositoryImpl) FindCartWithItems(ctx context.Context, cartID uuid.UUID) (*models.Cart, error) {
	var cart models.Cart
	err := r.GetTx(ctx).
		Preload("Items", func(db *gorm.DB) *gorm.DB { return db.Order("created_at ASC") }).
		Preload("Items.Product.Category").
		Preload("Items.Product.Variants").
		Where("id = ?", cartID).
		First(&cart).Error
	if err != nil {
		return nil, err
	}
	return &cart, nil
}

func (r *newCartRepositoryImpl) FindItemByMatch(ctx context.Context, cartID, productID uuid.UUID, variantsJSON string) (*models.CartItem, error) {
	var item models.CartItem
	err := r.GetTx(ctx).
		Where("cart_id = ? AND product_id = ? AND selected_variants = ?::jsonb", cartID, productID, variantsJSON).
		First(&item).Error
	if err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *newCartRepositoryImpl) FindItemByID(ctx context.Context, id uuid.UUID) (*models.CartItem, error) {
	var item models.CartItem
	err := r.GetTx(ctx).
		Where("id = ?", id).
		First(&item).Error
	if err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *newCartRepositoryImpl) CreateItem(ctx context.Context, item *models.CartItem) error {
	return r.GetTx(ctx).Create(item).Error
}

func (r *newCartRepositoryImpl) UpdateItem(ctx context.Context, item *models.CartItem) error {
	return r.GetTx(ctx).Save(item).Error
}

func (r *newCartRepositoryImpl) DeleteItem(ctx context.Context, id uuid.UUID) error {
	return r.GetTx(ctx).Where("id = ?", id).Delete(&models.CartItem{}).Error
}

func (r *newCartRepositoryImpl) DeleteItemsByCartID(ctx context.Context, cartID uuid.UUID) error {
	return r.GetTx(ctx).Where("cart_id = ?", cartID).Delete(&models.CartItem{}).Error
}
