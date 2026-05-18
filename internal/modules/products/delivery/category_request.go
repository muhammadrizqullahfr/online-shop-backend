package delivery

type CreateCategoryRequest struct {
	Name     string `json:"name" validate:"required,min=2,max=50"`
	Slug     string `json:"slug" validate:"required,min=2,max=50"`
	ParentID *uint  `json:"parent_id"`
}

type UpdateCategoryRequest struct {
	Name     string `json:"name" validate:"required,min=2,max=50"`
	Slug     string `json:"slug" validate:"required,min=2,max=50"`
	ParentID *uint  `json:"parent_id"`
}
