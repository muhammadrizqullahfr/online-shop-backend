package impl

import (
	"context"

	"github.com/RakaMurdiarta/online-shop-system/internal/models"
	"github.com/RakaMurdiarta/online-shop-system/internal/modules/orders/repository"
	"github.com/RakaMurdiarta/online-shop-system/pkg/database"
	"gorm.io/gorm"
)

type newOrderRepositoryImpl struct {
	*database.TransactionManagerImpl
}

func NewNewOrderRepository(tm *database.TransactionManagerImpl) repository.OrderRepository {
	return &newOrderRepositoryImpl{TransactionManagerImpl: tm}
}

func (r *newOrderRepositoryImpl) CreateOrder(ctx context.Context, order *models.Order) error {
	return r.GetTx(ctx).Create(order).Error
}

func (r *newOrderRepositoryImpl) CreateItems(ctx context.Context, items []models.OrderItem) error {
	if len(items) == 0 {
		return nil
	}
	return r.GetTx(ctx).Create(&items).Error
}

func (r *newOrderRepositoryImpl) UpdateOrder(ctx context.Context, order *models.Order) error {
	return r.GetTx(ctx).Save(order).Error
}

func (r *newOrderRepositoryImpl) UpdateOrderStatus(ctx context.Context, id, status string) error {
	return r.GetTx(ctx).
		Model(&models.Order{}).
		Where("id = ?", id).
		Update("status", status).Error
}

func (r *newOrderRepositoryImpl) FindOrderByID(ctx context.Context, id string) (*models.Order, error) {
	var order models.Order
	if err := r.GetTx(ctx).Where("id = ?", id).First(&order).Error; err != nil {
		return nil, err
	}
	return &order, nil
}

func (r *newOrderRepositoryImpl) FindOrderWithItems(ctx context.Context, id string) (*models.Order, error) {
	var order models.Order
	err := r.GetTx(ctx).
		Preload("Items", func(db *gorm.DB) *gorm.DB { return db.Order("created_at ASC") }).
		Where("id = ?", id).
		First(&order).Error
	if err != nil {
		return nil, err
	}
	return &order, nil
}

func (r *newOrderRepositoryImpl) ListOrders(ctx context.Context, limit, offset int, search, status, userID string) ([]models.Order, int64, error) {
	var orders []models.Order
	var total int64

	query := r.GetTx(ctx).Model(&models.Order{})

	if search != "" {
		like := "%" + search + "%"
		query = query.Where("customer_name ILIKE ? OR customer_email ILIKE ? OR id ILIKE ?", like, like, like)
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}
	if userID != "" {
		query = query.Where("user_id = ?", userID)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := query.
		Preload("Items", func(db *gorm.DB) *gorm.DB { return db.Order("created_at ASC") }).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&orders).Error
	return orders, total, err
}

func (r *newOrderRepositoryImpl) DeleteOrder(ctx context.Context, id string) error {
	return r.GetTx(ctx).Where("id = ?", id).Delete(&models.Order{}).Error
}

func (r *newOrderRepositoryImpl) CreateItem(ctx context.Context, item *models.OrderItem) error {
	return r.GetTx(ctx).Create(item).Error
}

func (r *newOrderRepositoryImpl) UpdateItem(ctx context.Context, item *models.OrderItem) error {
	return r.GetTx(ctx).Save(item).Error
}

func (r *newOrderRepositoryImpl) FindItemByID(ctx context.Context, id string) (*models.OrderItem, error) {
	var item models.OrderItem
	if err := r.GetTx(ctx).Where("id = ?", id).First(&item).Error; err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *newOrderRepositoryImpl) FindItemsByOrderID(ctx context.Context, orderID string) ([]models.OrderItem, error) {
	var items []models.OrderItem
	err := r.GetTx(ctx).
		Where("order_id = ?", orderID).
		Order("created_at ASC").
		Find(&items).Error
	return items, err
}

func (r *newOrderRepositoryImpl) DeleteItem(ctx context.Context, id string) error {
	return r.GetTx(ctx).Where("id = ?", id).Delete(&models.OrderItem{}).Error
}

func (r *newOrderRepositoryImpl) DeleteItemsByOrderID(ctx context.Context, orderID string) error {
	return r.GetTx(ctx).Where("order_id = ?", orderID).Delete(&models.OrderItem{}).Error
}
