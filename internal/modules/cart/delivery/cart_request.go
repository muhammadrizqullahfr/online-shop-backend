package delivery

import "github.com/google/uuid"

type SelectedVariantInput struct {
	Value string `json:"value" validate:"max=100"`
	Name  string `json:"name" validate:"required,max=100"`
	Type  string `json:"type" validate:"required,max=100"`
}

type AddItemRequest struct {
	ProductID        uuid.UUID              `json:"productId" validate:"required"`
	Quantity         int                    `json:"quantity" validate:"required,gt=0"`
	SelectedVariants []SelectedVariantInput `json:"selectedVariants" validate:"dive"`
}

type UpdateItemQuantityRequest struct {
	Quantity int `json:"quantity" validate:"required,gt=0"`
}
