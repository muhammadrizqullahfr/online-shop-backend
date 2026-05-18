package impl

import (
	"context"
	"errors"
	"fmt"

	"github.com/RakaMurdiarta/online-shop-system/internal/modules/orders/delivery"
	"github.com/RakaMurdiarta/online-shop-system/internal/modules/orders/repository"
	"github.com/RakaMurdiarta/online-shop-system/internal/modules/orders/services"
	"github.com/RakaMurdiarta/online-shop-system/pkg/shared"
	"github.com/xendit/xendit-go/v7/invoice"
	"gorm.io/gorm"
)

type orderCallbackServiceImpl struct {
	repo   repository.OrderRepository
	xendit *shared.XenditClient
}

func NewOrderCallbackService(repo repository.OrderRepository, xendit *shared.XenditClient) services.OrderCallbackService {
	return &orderCallbackServiceImpl{repo: repo, xendit: xendit}
}

func mapInvoiceStatus(status invoice.InvoiceStatus) string {
	switch status {
	case invoice.INVOICESTATUS_PAID, invoice.INVOICESTATUS_SETTLED:
		return "paid"
	case invoice.INVOICESTATUS_EXPIRED:
		return "expired"
	case invoice.INVOICESTATUS_PENDING:
		return "pending"
	default:
		return string(status)
	}
}

func (s *orderCallbackServiceImpl) HandleCallback(ctx context.Context, req *delivery.XenditCallbackRequest) error {
	if req.ExternalID == "" {
		return fmt.Errorf("external_id is required")
	}
	if req.ID == "" {
		return fmt.Errorf("invoice id is required")
	}

	return s.repo.WithTransaction(ctx, func(txCtx context.Context) error {
		order, err := s.repo.FindOrderByID(txCtx, req.ExternalID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return fmt.Errorf("order not found for external_id %s", req.ExternalID)
			}
			return err
		}

		if order.Status != nil && *order.Status == "paid" {
			return nil
		}

		invoiceResp, _, errXendit := s.xendit.API.InvoiceApi.GetInvoiceById(ctx, req.ID).Execute()
		if errXendit != nil {
			return fmt.Errorf("failed to verify xendit invoice: %w", errXendit)
		}

		if invoiceResp.ExternalId != order.ID {
			return fmt.Errorf("external_id mismatch: callback=%s xendit=%s", order.ID, invoiceResp.ExternalId)
		}

		newStatus := mapInvoiceStatus(invoiceResp.Status)
		if err := s.repo.UpdateOrderStatus(txCtx, order.ID, newStatus); err != nil {
			return fmt.Errorf("failed to update order status: %w", err)
		}
		return nil
	})
}
