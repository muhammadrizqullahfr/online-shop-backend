package delivery

import (
	"time"

	"github.com/RakaMurdiarta/online-shop-system/internal/models"
)

type FeedbackResponse struct {
	ID         int       `json:"id"`
	Name       string    `json:"name"`
	Email      string    `json:"email"`
	Subject    string    `json:"subject"`
	Message    string    `json:"message"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
	EmailSent  bool      `json:"email_sent"`
	EmailError string    `json:"email_error,omitempty"`
}

func ToFeedbackResponse(f *models.Feedback) *FeedbackResponse {
	if f == nil {
		return nil
	}
	return &FeedbackResponse{
		ID:        f.ID,
		Name:      f.Name,
		Email:     f.Email,
		Subject:   f.Subject,
		Message:   f.Message,
		CreatedAt: f.CreatedAt,
		UpdatedAt: f.UpdatedAt,
	}
}
