package impl

import (
	"context"
	"errors"
	"fmt"
	"math"

	"github.com/RakaMurdiarta/online-shop-system/internal/models"
	"github.com/RakaMurdiarta/online-shop-system/internal/modules/orders/delivery"
	"github.com/RakaMurdiarta/online-shop-system/internal/modules/orders/repository"
	"github.com/RakaMurdiarta/online-shop-system/internal/modules/orders/services"
	"github.com/RakaMurdiarta/online-shop-system/pkg/common/response"
	"github.com/RakaMurdiarta/online-shop-system/pkg/shared"
	"github.com/google/uuid"
	"github.com/xendit/xendit-go/v7/invoice"
	"gorm.io/gorm"
)

type orderServiceImpl struct {
	repo   repository.OrderRepository
	xendit *shared.XenditClient
}

func NewOrderService(repo repository.OrderRepository, xendit *shared.XenditClient) services.OrderService {
	return &orderServiceImpl{repo: repo, xendit: xendit}
}

func optionalString(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

func toVariants(inputs []delivery.SelectedVariantInput) models.SelectedVariants {
	out := make(models.SelectedVariants, 0, len(inputs))
	for _, v := range inputs {
		out = append(out, models.SelectedVariant{
			Value: v.Value,
			Name:  v.Name,
			Type:  v.Type,
		})
	}
	return out
}

func computeTotal(items []models.OrderItem) float64 {
	var total float64
	for i := range items {
		total += items[i].Price * float64(items[i].Quantity)
	}
	return total
}

func (s *orderServiceImpl) CreateOrder(ctx context.Context, userID *string, req *delivery.CreateOrderRequest) (*delivery.OrderResponse, error) {
	orderID := uuid.NewString()

	items := make([]models.OrderItem, 0, len(req.Items))
	for _, it := range req.Items {
		items = append(items, models.OrderItem{
			ID:               it.ID,
			OrderID:          orderID,
			ProductID:        it.ProductID,
			Name:             it.Name,
			Price:            it.Price,
			Quantity:         it.Quantity,
			Image:            optionalString(it.Image),
			SelectedVariants: toVariants(it.SelectedVariants),
		})
	}

	customerEmail := req.CustomerEmail
	invoiceRequest := invoice.CreateInvoiceRequest{
		ExternalId: orderID,
		Amount:     float64(req.Total),
		Customer: &invoice.CustomerObject{
			Email: *invoice.NewNullableString(&customerEmail),
		},
	}

	resp, _, errXendit := s.xendit.API.InvoiceApi.CreateInvoice(ctx).
		CreateInvoiceRequest(invoiceRequest).
		Execute()
	if errXendit != nil {
		return nil, fmt.Errorf("failed to create xendit invoice: %w", errXendit)
	}

	paymentURL := resp.InvoiceUrl

	order := &models.Order{
		City:            req.City,
		ID:              orderID,
		CustomerAddress: req.CustomerAddress,
		CustomerEmail:   req.CustomerEmail,
		CustomerName:    req.CustomerName,
		CustomerPhone:   req.CustomerPhone,
		State:           req.State,
		ZipCode:         req.ZipCode,
		PaymentURL:      &paymentURL,
		UserID:          userID,
		Total:           float64(req.Total),
	}

	if err := s.repo.WithTransaction(ctx, func(tx context.Context) error {
		if err := s.repo.CreateOrder(tx, order); err != nil {
			return fmt.Errorf("failed to create order: %w", err)
		}
		if err := s.repo.CreateItems(tx, items); err != nil {
			return fmt.Errorf("failed to create order items: %w", err)
		}
		return nil
	}); err != nil {
		return nil, err
	}

	return s.GetOrderByID(ctx, orderID)
}

func (s *orderServiceImpl) GetOrderByID(ctx context.Context, id string) (*delivery.OrderResponse, error) {
	order, err := s.repo.FindOrderWithItems(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("order not found")
		}
		return nil, err
	}
	return delivery.ToOrderResponse(order), nil
}

func (s *orderServiceImpl) ListOrders(ctx context.Context, page, limit int, search, status, userID string) (*delivery.OrderListResponse, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	offset := (page - 1) * limit

	orders, total, err := s.repo.ListOrders(ctx, limit, offset, search, status, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch orders: %w", err)
	}

	items := make([]delivery.OrderResponse, 0, len(orders))
	for i := range orders {
		items = append(items, *delivery.ToOrderResponse(&orders[i]))
	}

	totalPage := 0
	if total > 0 {
		totalPage = int(math.Ceil(float64(total) / float64(limit)))
	}

	return &delivery.OrderListResponse{
		Items: items,
		Paging: response.PagingResponse{
			CurrentPage: page,
			TotalPage:   totalPage,
			TotalItems:  total,
		},
	}, nil
}

func (s *orderServiceImpl) UpdateOrder(ctx context.Context, id string, req *delivery.UpdateOrderRequest) (*delivery.OrderResponse, error) {
	err := s.repo.WithTransaction(ctx, func(tx context.Context) error {
		order, err := s.repo.FindOrderByID(tx, id)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return fmt.Errorf("order not found")
			}
			return err
		}

		if req.CustomerName != "" {
			order.CustomerName = req.CustomerName
		}
		if req.CustomerEmail != "" {
			order.CustomerEmail = req.CustomerEmail
		}
		if req.CustomerPhone != "" {
			order.CustomerPhone = req.CustomerPhone
		}
		if req.CustomerAddress != "" {
			order.CustomerAddress = req.CustomerAddress
		}
		if req.City != "" {
			order.City = req.City
		}
		if req.State != "" {
			order.State = req.State
		}
		if req.ZipCode != "" {
			order.ZipCode = req.ZipCode
		}
		if req.Status != "" {
			order.Status = optionalString(req.Status)
		}
		if req.PaymentURL != "" {
			order.PaymentURL = optionalString(req.PaymentURL)
		}

		return s.repo.UpdateOrder(tx, order)
	})
	if err != nil {
		return nil, err
	}
	return s.GetOrderByID(ctx, id)
}

func (s *orderServiceImpl) DeleteOrder(ctx context.Context, id string) error {
	return s.repo.WithTransaction(ctx, func(tx context.Context) error {
		if _, err := s.repo.FindOrderByID(tx, id); err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return fmt.Errorf("order not found")
			}
			return err
		}
		if err := s.repo.DeleteItemsByOrderID(tx, id); err != nil {
			return fmt.Errorf("failed to delete order items: %w", err)
		}
		if err := s.repo.DeleteOrder(tx, id); err != nil {
			return fmt.Errorf("failed to delete order: %w", err)
		}
		return nil
	})
}

func (s *orderServiceImpl) recalculateTotal(ctx context.Context, orderID string) error {
	items, err := s.repo.FindItemsByOrderID(ctx, orderID)
	if err != nil {
		return err
	}
	order, err := s.repo.FindOrderByID(ctx, orderID)
	if err != nil {
		return err
	}
	order.Total = computeTotal(items)
	return s.repo.UpdateOrder(ctx, order)
}

func (s *orderServiceImpl) AddItem(ctx context.Context, orderID string, req *delivery.CreateOrderItemRequest) (*delivery.OrderItemResponse, error) {
	item := &models.OrderItem{
		ID:               uuid.NewString(),
		OrderID:          orderID,
		ProductID:        req.ProductID,
		Name:             req.Name,
		Price:            req.Price,
		Quantity:         req.Quantity,
		Image:            optionalString(req.Image),
		SelectedVariants: toVariants(req.SelectedVariants),
	}

	err := s.repo.WithTransaction(ctx, func(tx context.Context) error {
		if _, err := s.repo.FindOrderByID(tx, orderID); err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return fmt.Errorf("order not found")
			}
			return err
		}
		if err := s.repo.CreateItem(tx, item); err != nil {
			return fmt.Errorf("failed to create item: %w", err)
		}
		return s.recalculateTotal(tx, orderID)
	})
	if err != nil {
		return nil, err
	}

	resp := delivery.ToOrderItemResponse(item)
	return &resp, nil
}

func (s *orderServiceImpl) UpdateItem(ctx context.Context, orderID, itemID string, req *delivery.UpdateOrderItemRequest) (*delivery.OrderItemResponse, error) {
	var updated *models.OrderItem

	err := s.repo.WithTransaction(ctx, func(tx context.Context) error {
		item, err := s.repo.FindItemByID(tx, itemID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return fmt.Errorf("order item not found")
			}
			return err
		}
		if item.OrderID != orderID {
			return fmt.Errorf("order item not found")
		}

		if req.Name != "" {
			item.Name = req.Name
		}
		if req.Price != nil {
			item.Price = *req.Price
		}
		if req.Quantity != nil {
			item.Quantity = *req.Quantity
		}
		if req.Image != "" {
			item.Image = optionalString(req.Image)
		}
		if req.SelectedVariants != nil {
			item.SelectedVariants = toVariants(req.SelectedVariants)
		}

		if err := s.repo.UpdateItem(tx, item); err != nil {
			return fmt.Errorf("failed to update item: %w", err)
		}
		updated = item
		return s.recalculateTotal(tx, orderID)
	})
	if err != nil {
		return nil, err
	}

	resp := delivery.ToOrderItemResponse(updated)
	return &resp, nil
}

func (s *orderServiceImpl) DeleteItem(ctx context.Context, orderID, itemID string) error {
	return s.repo.WithTransaction(ctx, func(tx context.Context) error {
		item, err := s.repo.FindItemByID(tx, itemID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return fmt.Errorf("order item not found")
			}
			return err
		}
		if item.OrderID != orderID {
			return fmt.Errorf("order item not found")
		}

		if err := s.repo.DeleteItem(tx, itemID); err != nil {
			return fmt.Errorf("failed to delete item: %w", err)
		}
		return s.recalculateTotal(tx, orderID)
	})
}

func (s *orderServiceImpl) ListItems(ctx context.Context, orderID string) ([]delivery.OrderItemResponse, error) {
	if _, err := s.repo.FindOrderByID(ctx, orderID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("order not found")
		}
		return nil, err
	}

	items, err := s.repo.FindItemsByOrderID(ctx, orderID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch items: %w", err)
	}

	out := make([]delivery.OrderItemResponse, 0, len(items))
	for i := range items {
		out = append(out, delivery.ToOrderItemResponse(&items[i]))
	}
	return out, nil
}
