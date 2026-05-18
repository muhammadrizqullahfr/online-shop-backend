package delivery

import (
	"time"

	"github.com/RakaMurdiarta/online-shop-system/internal/models"
)

type CategoryResponse struct {
	ID        uint               `json:"id"`
	Name      string             `json:"name"`
	Slug      string             `json:"slug"`
	ParentID  *uint              `json:"parent_id"`
	Children  []CategoryResponse `json:"children,omitempty"`
	CreatedAt time.Time          `json:"created_at"`
	UpdatedAt time.Time          `json:"updated_at"`
}

func ToCategoryResponse(category *models.Category) *CategoryResponse {
	if category == nil {
		return nil
	}

	resp := &CategoryResponse{
		ID:        category.ID,
		Name:      category.Name,
		Slug:      category.Slug,
		ParentID:  category.ParentID,
		CreatedAt: category.CreatedAt,
		UpdatedAt: category.UpdatedAt,
	}

	if len(category.Children) > 0 {
		resp.Children = make([]CategoryResponse, 0, len(category.Children))
		for _, child := range category.Children {
			c := child
			resp.Children = append(resp.Children, *ToCategoryResponse(&c))
		}
	}

	return resp

}

func BuildCategoryTree(categories []models.Category) []CategoryResponse {

	categoryMap := make(map[uint]*CategoryResponse)

	for _, category := range categories {
		cat := category

		categoryMap[cat.ID] = &CategoryResponse{
			ID:        cat.ID,
			Name:      cat.Name,
			Slug:      cat.Slug,
			ParentID:  cat.ParentID,
			CreatedAt: cat.CreatedAt,
			UpdatedAt: cat.UpdatedAt,
			Children:  []CategoryResponse{},
		}
	}

	rootPointers := []*CategoryResponse{}

	for _, category := range categories {

		current := categoryMap[category.ID]

		if category.ParentID == nil {
			rootPointers = append(rootPointers, current)
			continue
		}

		parent, ok := categoryMap[*category.ParentID]
		if ok {
			parent.Children = append(parent.Children, *current)
		}
	}

	result := make([]CategoryResponse, 0, len(rootPointers))

	for _, root := range rootPointers {
		result = append(result, *root)
	}

	return result
}
