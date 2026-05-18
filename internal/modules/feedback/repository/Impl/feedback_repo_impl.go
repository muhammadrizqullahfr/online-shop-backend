package impl

import (
	"context"

	"github.com/RakaMurdiarta/online-shop-system/internal/models"
	"github.com/RakaMurdiarta/online-shop-system/internal/modules/feedback/repository"
	"github.com/RakaMurdiarta/online-shop-system/pkg/database"
)

type feedbackRepositoryImpl struct {
	*database.TransactionManagerImpl
}

func NewFeedbackRepository(tm *database.TransactionManagerImpl) repository.FeedbackRepository {
	return &feedbackRepositoryImpl{TransactionManagerImpl: tm}
}

func (r *feedbackRepositoryImpl) Create(ctx context.Context, f *models.Feedback) error {
	return r.GetTx(ctx).Create(f).Error
}

func (r *feedbackRepositoryImpl) FindByID(ctx context.Context, id int) (*models.Feedback, error) {
	var f models.Feedback
	if err := r.GetTx(ctx).Where("id = ?", id).First(&f).Error; err != nil {
		return nil, err
	}
	return &f, nil
}
