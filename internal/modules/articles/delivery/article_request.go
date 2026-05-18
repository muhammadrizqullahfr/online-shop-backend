package delivery

type CreateArticleRequest struct {
	Title       string `json:"title" validate:"required,max=255"`
	Description string `json:"description"`
	Content     string `json:"content"`
	Image       string `json:"image" validate:"omitempty,url,max=255"`
	Author      string `json:"author" validate:"max=255"`
	Category    string `json:"category" validate:"max=255"`
	Featured    bool   `json:"featured"`
	ReadTime    int    `json:"readTime" validate:"gte=0"`
	Excerpt     string `json:"excerpt"`
}

type UpdateArticleRequest struct {
	Title       string `json:"title" validate:"required,max=255"`
	Description string `json:"description"`
	Content     string `json:"content"`
	Image       string `json:"image" validate:"omitempty,url,max=255"`
	Author      string `json:"author" validate:"max=255"`
	Category    string `json:"category" validate:"max=255"`
	Featured    bool   `json:"featured"`
	ReadTime    int    `json:"readTime" validate:"gte=0"`
	Excerpt     string `json:"excerpt"`
}
