package services

import (
	"context"

	"github.com/RakaMurdiarta/online-shop-system/internal/modules/mailer/delivery"
)

type MailerService interface {
	SendEmail(ctx context.Context, req delivery.SendEmailRequest) (*delivery.SendEmailResult, error)
}
