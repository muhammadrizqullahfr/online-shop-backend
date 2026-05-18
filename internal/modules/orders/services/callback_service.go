package services

import (
	"context"

	"github.com/RakaMurdiarta/online-shop-system/internal/modules/orders/delivery"
)

type OrderCallbackService interface {
	HandleCallback(ctx context.Context, req *delivery.XenditCallbackRequest) error
}
