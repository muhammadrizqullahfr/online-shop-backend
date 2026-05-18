package repository

import (
	"context"

	"github.com/RakaMurdiarta/online-shop-system/internal/models"
	"github.com/RakaMurdiarta/online-shop-system/pkg/database"
)

type FeedbackRepository interface {
	Create(ctx context.Context, f *models.Feedback) error
	FindByID(ctx context.Context, id int) (*models.Feedback, error)

	database.TransactionManager
}
