package impl

import (
	"context"

	"github.com/RakaMurdiarta/online-shop-system/internal/modules/mailer/delivery"
	"github.com/RakaMurdiarta/online-shop-system/internal/modules/mailer/services"
	"github.com/RakaMurdiarta/online-shop-system/pkg/mailer"
)

type mailerServiceImpl struct {
	transport mailer.Transport
}

func NewMailerService(transport mailer.Transport) services.MailerService {
	return &mailerServiceImpl{transport: transport}
}

func (s *mailerServiceImpl) SendEmail(ctx context.Context, req delivery.SendEmailRequest) (*delivery.SendEmailResult, error) {
	msg := mailer.Message{
		To:       req.To,
		Cc:       req.Cc,
		Bcc:      req.Bcc,
		From:     req.From,
		FromName: req.FromName,
		ReplyTo:  req.ReplyTo,
		Subject:  req.Subject,
		Body:     req.Body,
		IsHTML:   req.IsHTML,
	}

	receipt, err := s.transport.Send(ctx, msg)
	if err != nil {
		return nil, err
	}

	return &delivery.SendEmailResult{
		ID:      receipt.ID,
		To:      receipt.To,
		Subject: receipt.Subject,
		SentAt:  receipt.SentAt,
	}, nil
}
