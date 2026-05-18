package repository

import (
	"context"

	"github.com/RakaMurdiarta/online-shop-system/internal/models"
	"github.com/RakaMurdiarta/online-shop-system/pkg/database"
)

type OrderRepository interface {
	CreateOrder(ctx context.Context, order *models.Order) error
	CreateItems(ctx context.Context, items []models.OrderItem) error
	UpdateOrder(ctx context.Context, order *models.Order) error
	UpdateOrderStatus(ctx context.Context, id, status string) error
	FindOrderByID(ctx context.Context, id string) (*models.Order, error)
	FindOrderWithItems(ctx context.Context, id string) (*models.Order, error)
	ListOrders(ctx context.Context, limit, offset int, search, status, userID string) ([]models.Order, int64, error)
	DeleteOrder(ctx context.Context, id string) error

	CreateItem(ctx context.Context, item *models.OrderItem) error
	UpdateItem(ctx context.Context, item *models.OrderItem) error
	FindItemByID(ctx context.Context, id string) (*models.OrderItem, error)
	FindItemsByOrderID(ctx context.Context, orderID string) ([]models.OrderItem, error)
	DeleteItem(ctx context.Context, id string) error
	DeleteItemsByOrderID(ctx context.Context, orderID string) error

	database.TransactionManager
}
