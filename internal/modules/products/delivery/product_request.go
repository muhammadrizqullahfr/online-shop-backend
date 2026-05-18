package delivery

type CreateNewVariantRequest struct {
	Name            string   `json:"name" validate:"required,max=100"`
	Type            string   `json:"type" validate:"required,max=50"`
	Value           string   `json:"value" validate:"required,max=100"`
	Price           *float64 `json:"price,omitempty" validate:"omitempty,gte=0"`
	PriceAdjustment *float64 `json:"priceAdjustment,omitempty"`
	Stock           int      `json:"stock" validate:"gte=0"`
}

type CreateNewProductRequest struct {
	Name          string                    `json:"name" validate:"required,max=255"`
	Description   string                    `json:"description"`
	CategorySlug  string                    `json:"category" validate:"required,max=50"`
	Price         *float64                  `json:"price" validate:"required"`
	OriginalPrice *float64                  `json:"originalPrice,omitempty" validate:"omitempty,gte=0"`
	Image         string                    `json:"image" validate:"required,url"`
	Rating        float64                   `json:"rating" validate:"gte=0,lte=5"`
	Reviews       int                       `json:"reviews" validate:"gte=0"`
	Stock         int                       `json:"stock" validate:"gte=0"`
	InStock       bool                      `json:"inStock"`
	Featured      bool                      `json:"featured"`
	Tags          []string                  `json:"tags"`
	Variants      []CreateNewVariantRequest `json:"variants" validate:"dive"`
}

type UpdateNewProductRequest struct {
	Name          string                    `json:"name" validate:"required,max=255"`
	Description   string                    `json:"description"`
	CategorySlug  string                    `json:"category" validate:"required,max=50"`
	Price         *float64                  `json:"price" validate:"required,gte=0"`
	OriginalPrice *float64                  `json:"originalPrice,omitempty" validate:"omitempty,gte=0"`
	Image         string                    `json:"image" validate:"required,url"`
	Rating        float64                   `json:"rating" validate:"gte=0,lte=5"`
	Reviews       int                       `json:"reviews" validate:"gte=0"`
	Stock         int                       `json:"stock" validate:"gte=0"`
	InStock       bool                      `json:"inStock"`
	Featured      bool                      `json:"featured"`
	Tags          []string                  `json:"tags"`
	Variants      []CreateNewVariantRequest `json:"variants" validate:"dive"`
}
