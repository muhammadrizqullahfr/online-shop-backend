package services

import (
	"context"

	"github.com/RakaMurdiarta/online-shop-system/internal/modules/orders/delivery"
)

type OrderService interface {
	CreateOrder(ctx context.Context, userID *string, req *delivery.CreateOrderRequest) (*delivery.OrderResponse, error)
	GetOrderByID(ctx context.Context, id string) (*delivery.OrderResponse, error)
	ListOrders(ctx context.Context, page, limit int, search, status, userID string) (*delivery.OrderListResponse, error)
	UpdateOrder(ctx context.Context, id string, req *delivery.UpdateOrderRequest) (*delivery.OrderResponse, error)
	DeleteOrder(ctx context.Context, id string) error

	AddItem(ctx context.Context, orderID string, req *delivery.CreateOrderItemRequest) (*delivery.OrderItemResponse, error)
	UpdateItem(ctx context.Context, orderID, itemID string, req *delivery.UpdateOrderItemRequest) (*delivery.OrderItemResponse, error)
	DeleteItem(ctx context.Context, orderID, itemID string) error
	ListItems(ctx context.Context, orderID string) ([]delivery.OrderItemResponse, error)
}
