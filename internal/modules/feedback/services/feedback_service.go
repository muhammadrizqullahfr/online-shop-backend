package services

import (
	"context"

	"github.com/RakaMurdiarta/online-shop-system/internal/modules/feedback/delivery"
)

type FeedbackService interface {
	Create(ctx context.Context, req *delivery.CreateFeedbackRequest) (*delivery.FeedbackResponse, error)
}
