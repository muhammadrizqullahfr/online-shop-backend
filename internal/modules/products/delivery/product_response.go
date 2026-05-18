package delivery

import (
	"time"

	"github.com/RakaMurdiarta/online-shop-system/internal/models"
	"github.com/RakaMurdiarta/online-shop-system/pkg/common/response"
)

type NewVariantResponse struct {
	ID              string   `json:"id"`
	Name            string   `json:"name"`
	Type            string   `json:"type"`
	Value           string   `json:"value"`
	Price           *float64 `json:"price,omitempty"`
	PriceAdjustment *float64 `json:"priceAdjustment,omitempty"`
	Stock           int      `json:"stock"`
}

type NewProductResponse struct {
	ID            string               `json:"id"`
	Name          string               `json:"name"`
	Description   string               `json:"description"`
	Category      string               `json:"category"`
	Price         *float64             `json:"price"`
	OriginalPrice *float64             `json:"originalPrice,omitempty"`
	Image         string               `json:"image"`
	Rating        float64              `json:"rating"`
	Reviews       int                  `json:"reviews"`
	Stock         int                  `json:"stock"`
	InStock       bool                 `json:"inStock"`
	Featured      bool                 `json:"featured"`
	Tags          []string             `json:"tags"`
	CreatedAt     time.Time            `json:"created_at"`
	Variants      []NewVariantResponse `json:"variants,omitempty"`
}

type NewProductListResponse struct {
	Items  []NewProductResponse    `json:"items"`
	Paging response.PagingResponse `json:"paging"`
}

func ToNewVariantResponse(v *models.ProductVariant) NewVariantResponse {
	return NewVariantResponse{
		ID:              v.ID.String(),
		Name:            v.Name,
		Type:            v.Type,
		Value:           v.Value,
		Price:           v.Price,
		PriceAdjustment: v.PriceAdjustment,
		Stock:           v.Stock,
	}
}

func ToNewProductResponse(p *models.Product) *NewProductResponse {
	if p == nil {
		return nil
	}

	tags := []string(p.Tags)
	if tags == nil {
		tags = []string{}
	}

	resp := &NewProductResponse{
		ID:            p.ID.String(),
		Name:          p.Name,
		Description:   p.Description,
		Category:      p.Category.Slug,
		Price:         p.Price,
		OriginalPrice: p.OriginalPrice,
		Image:         p.Image,
		Rating:        p.Rating,
		Reviews:       p.Reviews,
		Stock:         p.Stock,
		InStock:       p.InStock,
		Featured:      p.Featured,
		Tags:          tags,
		CreatedAt:     p.CreatedAt,
	}

	if len(p.Variants) > 0 {
		resp.Variants = make([]NewVariantResponse, 0, len(p.Variants))
		for i := range p.Variants {
			resp.Variants = append(resp.Variants, ToNewVariantResponse(&p.Variants[i]))
		}
	}

	return resp
}
