package delivery

import (
	"time"

	"github.com/RakaMurdiarta/online-shop-system/internal/models"
	"github.com/RakaMurdiarta/online-shop-system/pkg/common/response"
)

type ArticleResponse struct {
	ID          string     `json:"id"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Content     string     `json:"content"`
	Image       string     `json:"image"`
	Author      string     `json:"author"`
	Category    string     `json:"category"`
	PublishedAt *time.Time `json:"publishedAt,omitempty"`
	Featured    bool       `json:"featured"`
	ReadTime    int        `json:"readTime"`
	Excerpt     string     `json:"excerpt"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

type ArticleListResponse struct {
	Items  []ArticleResponse       `json:"items"`
	Paging response.PagingResponse `json:"paging"`
}

func ToArticleResponse(a *models.Article) *ArticleResponse {
	if a == nil {
		return nil
	}

	return &ArticleResponse{
		ID:          a.ID,
		Title:       a.Title,
		Description: deref(a.Description),
		Content:     deref(a.Content),
		Image:       deref(a.Image),
		Author:      deref(a.Author),
		Category:    deref(a.Category),
		PublishedAt: a.PublishedAt,
		Featured:    a.Featured,
		ReadTime:    a.ReadTime,
		Excerpt:     deref(a.Excerpt),
		CreatedAt:   a.CreatedAt,
		UpdatedAt:   a.UpdatedAt,
	}
}

func deref(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}
