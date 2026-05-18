package delivery

import (
	"time"

	"github.com/RakaMurdiarta/online-shop-system/internal/models"
	"github.com/google/uuid"
)

type SelectedVariantResponse struct {
	Value string `json:"value"`
	Name  string `json:"name"`
	Type  string `json:"type"`
}

type CartItemProductSnapshot struct {
	ID       uuid.UUID `json:"id"`
	Name     string    `json:"name"`
	Image    string    `json:"image"`
	Price    *float64  `json:"price"`
	Stock    int       `json:"stock"`
	Category string    `json:"category"`
}

type CartItemResponse struct {
	ID               uuid.UUID                 `json:"id"`
	ProductID        uuid.UUID                 `json:"product_id"`
	Quantity         int                       `json:"quantity"`
	SelectedVariants []SelectedVariantResponse `json:"selected_variants"`
	Product          *CartItemProductSnapshot  `json:"product,omitempty"`
	CreatedAt        time.Time                 `json:"created_at"`
	UpdatedAt        time.Time                 `json:"updated_at"`
}

type CartResponse struct {
	ID        uuid.UUID          `json:"id"`
	UserID    uint               `json:"user_id"`
	Items     []CartItemResponse `json:"items"`
	CreatedAt time.Time          `json:"created_at"`
	UpdatedAt time.Time          `json:"updated_at"`
}

func ToCartItemResponse(item *models.CartItem) CartItemResponse {
	variants := make([]SelectedVariantResponse, 0, len(item.SelectedVariants))
	for _, v := range item.SelectedVariants {
		variants = append(variants, SelectedVariantResponse{
			Value: v.Value,
			Name:  v.Name,
			Type:  v.Type,
		})
	}

	resp := CartItemResponse{
		ID:               item.ID,
		ProductID:        item.ProductID,
		Quantity:         item.Quantity,
		SelectedVariants: variants,
		CreatedAt:        item.CreatedAt,
		UpdatedAt:        item.UpdatedAt,
	}

	if item.Product.ID != uuid.Nil {
		resp.Product = &CartItemProductSnapshot{
			ID:       item.Product.ID,
			Name:     item.Product.Name,
			Image:    item.Product.Image,
			Price:    resolveUnitPrice(item),
			Stock:    resolveStock(item),
			Category: item.Product.Category.Name,
		}
	}

	return resp
}

// resolveStock returns the available stock for the cart line. If no variants
// were selected, it returns the product's stock. If variants were selected, it
// returns the minimum stock across all matched NewVariant rows (the limiting
// SKU). Unmatched selected variants fall back to product stock.
func resolveStock(item *models.CartItem) int {
	if len(item.SelectedVariants) == 0 {
		return item.Product.Stock
	}

	stock := -1
	for _, sv := range item.SelectedVariants {
		for i := range item.Product.Variants {
			v := &item.Product.Variants[i]
			if v.Name != sv.Name || v.Type != sv.Type || v.Value != sv.Value {
				continue
			}
			if stock == -1 || v.Stock < stock {
				stock = v.Stock
			}
			break
		}
	}
	if stock == -1 {
		return item.Product.Stock
	}
	return stock
}

// resolveUnitPrice returns the price to show on the cart item's product
// snapshot. If no variants were selected, it returns the product's original
// price. If variants were selected, the matched NewVariant's Price overrides
// the product price; if the variant only carries a PriceAdjustment, that
// adjustment is added to the product price.
func resolveUnitPrice(item *models.CartItem) *float64 {
	if item.Product.ID == uuid.Nil || item.Product.Price == nil {
		return nil
	}

	if len(item.SelectedVariants) == 0 {
		p := *item.Product.Price
		return &p
	}

	unit := *item.Product.Price
	for _, sv := range item.SelectedVariants {
		for i := range item.Product.Variants {
			v := &item.Product.Variants[i]
			if v.Name != sv.Name || v.Type != sv.Type || v.Value != sv.Value {
				continue
			}
			switch {
			case v.Price != nil:
				unit = *v.Price
			case v.PriceAdjustment != nil:
				unit += *v.PriceAdjustment
			}
			break
		}
	}
	return &unit
}

func ToCartResponse(cart *models.Cart) *CartResponse {
	if cart == nil {
		return nil
	}

	items := make([]CartItemResponse, 0, len(cart.Items))
	for i := range cart.Items {
		items = append(items, ToCartItemResponse(&cart.Items[i]))
	}

	return &CartResponse{
		ID:        cart.ID,
		UserID:    cart.UserID,
		Items:     items,
		CreatedAt: cart.CreatedAt,
		UpdatedAt: cart.UpdatedAt,
	}
}
