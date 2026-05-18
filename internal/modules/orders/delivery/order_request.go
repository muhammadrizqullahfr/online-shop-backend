package delivery

type SelectedVariantInput struct {
	Value string `json:"value" validate:"max=100"`
	Name  string `json:"name" validate:"required,max=100"`
	Type  string `json:"type" validate:"required,max=100"`
}

type OrderItemInput struct {
	Category         string                 `json:"category" validate:"omitempty,max=255"`
	ID               string                 `json:"id" validate:"required"`
	Image            string                 `json:"image" validate:"omitempty,max=255"`
	Name             string                 `json:"name" validate:"required,max=255"`
	Price            float64                `json:"price" validate:"gte=0"`
	ProductID        string                 `json:"product_id" validate:"required"`
	Quantity         int                    `json:"quantity" validate:"required,gt=0"`
	SelectedVariants []SelectedVariantInput `json:"selectedVariants" validate:"dive"`
}

type CreateOrderRequest struct {
	CustomerName    string           `json:"customer_name" validate:"required,max=255"`
	CustomerEmail   string           `json:"customer_email" validate:"required,email,max=255"`
	CustomerPhone   string           `json:"customer_phone" validate:"required,max=255"`
	CustomerAddress string           `json:"customer_address" validate:"required"`
	City            string           `json:"city" validate:"required,max=255"`
	State           string           `json:"state" validate:"required,max=255"`
	ZipCode         string           `json:"zip_code" validate:"required,max=255"`
	Status          string           `json:"status" validate:"omitempty,max=255"`
	PaymentURL      string           `json:"payment_url" validate:"omitempty,max=255"`
	Total           int              `json:"total" validate:"required,gt=0"`
	Items           []OrderItemInput `json:"items" validate:"required,min=1,dive"`
}

type UpdateOrderRequest struct {
	CustomerName    string `json:"customer_name" validate:"omitempty,max=255"`
	CustomerEmail   string `json:"customer_email" validate:"omitempty,email,max=255"`
	CustomerPhone   string `json:"customer_phone" validate:"omitempty,max=255"`
	CustomerAddress string `json:"customer_address" validate:"omitempty"`
	City            string `json:"city" validate:"omitempty,max=255"`
	State           string `json:"state" validate:"omitempty,max=255"`
	ZipCode         string `json:"zip_code" validate:"omitempty,max=255"`
	Status          string `json:"status" validate:"omitempty,max=255"`
	PaymentURL      string `json:"payment_url" validate:"omitempty,max=255"`
}

type CreateOrderItemRequest struct {
	OrderItemInput
}

type UpdateOrderItemRequest struct {
	Name             string                 `json:"name" validate:"omitempty,max=255"`
	Price            *float64               `json:"price" validate:"omitempty,gte=0"`
	Quantity         *int                   `json:"quantity" validate:"omitempty,gt=0"`
	Image            string                 `json:"image" validate:"omitempty,max=255"`
	SelectedVariants []SelectedVariantInput `json:"selected_variants" validate:"omitempty,dive"`
}
