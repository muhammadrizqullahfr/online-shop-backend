package delivery

import (
	"time"

	"github.com/RakaMurdiarta/online-shop-system/internal/models"
	"github.com/RakaMurdiarta/online-shop-system/pkg/common/response"
)

type SelectedVariantResponse struct {
	Value string `json:"value"`
	Name  string `json:"name"`
	Type  string `json:"type"`
}

type OrderItemResponse struct {
	ID               string                    `json:"id"`
	OrderID          string                    `json:"order_id"`
	ProductID        string                    `json:"product_id"`
	Name             string                    `json:"name"`
	Price            float64                   `json:"price"`
	Quantity         int                       `json:"quantity"`
	Image            string                    `json:"image"`
	SelectedVariants []SelectedVariantResponse `json:"selected_variants"`
	CreatedAt        time.Time                 `json:"created_at"`
	UpdatedAt        time.Time                 `json:"updated_at"`
}

type OrderResponse struct {
	ID              string              `json:"id"`
	CustomerName    string              `json:"customer_name"`
	CustomerEmail   string              `json:"customer_email"`
	CustomerPhone   string              `json:"customer_phone"`
	Total           float64             `json:"total"`
	Status          string              `json:"status"`
	PaymentURL      string              `json:"payment_url"`
	UserID          string              `json:"user_id"`
	CustomerAddress string              `json:"customer_address"`
	City            string              `json:"city"`
	State           string              `json:"state"`
	ZipCode         string              `json:"zip_code"`
	Items           []OrderItemResponse `json:"items"`
	CreatedAt       time.Time           `json:"created_at"`
	UpdatedAt       time.Time           `json:"updated_at"`
}

type OrderListResponse struct {
	Items  []OrderResponse         `json:"items"`
	Paging response.PagingResponse `json:"paging"`
}

func derefString(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

func ToOrderItemResponse(item *models.OrderItem) OrderItemResponse {
	variants := make([]SelectedVariantResponse, 0, len(item.SelectedVariants))
	for _, v := range item.SelectedVariants {
		variants = append(variants, SelectedVariantResponse{
			Value: v.Value,
			Name:  v.Name,
			Type:  v.Type,
		})
	}
	return OrderItemResponse{
		ID:               item.ID,
		OrderID:          item.OrderID,
		ProductID:        item.ProductID,
		Name:             item.Name,
		Price:            item.Price,
		Quantity:         item.Quantity,
		Image:            derefString(item.Image),
		SelectedVariants: variants,
		CreatedAt:        item.CreatedAt,
		UpdatedAt:        item.UpdatedAt,
	}
}

func ToOrderResponse(order *models.Order) *OrderResponse {
	if order == nil {
		return nil
	}
	items := make([]OrderItemResponse, 0, len(order.Items))
	for i := range order.Items {
		items = append(items, ToOrderItemResponse(&order.Items[i]))
	}
	return &OrderResponse{
		ID:              order.ID,
		CustomerName:    order.CustomerName,
		CustomerEmail:   order.CustomerEmail,
		CustomerPhone:   order.CustomerPhone,
		Total:           order.Total,
		Status:          derefString(order.Status),
		PaymentURL:      derefString(order.PaymentURL),
		UserID:          derefString(order.UserID),
		CustomerAddress: order.CustomerAddress,
		City:            order.City,
		State:           order.State,
		ZipCode:         order.ZipCode,
		Items:           items,
		CreatedAt:       order.CreatedAt,
		UpdatedAt:       order.UpdatedAt,
	}
}
