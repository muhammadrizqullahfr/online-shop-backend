package impl

import (
	"context"
	"fmt"

	"github.com/RakaMurdiarta/online-shop-system/internal/models"
	"github.com/RakaMurdiarta/online-shop-system/internal/modules/feedback/delivery"
	"github.com/RakaMurdiarta/online-shop-system/internal/modules/feedback/repository"
	"github.com/RakaMurdiarta/online-shop-system/internal/modules/feedback/services"
	mailerDelivery "github.com/RakaMurdiarta/online-shop-system/internal/modules/mailer/delivery"
	mailerServices "github.com/RakaMurdiarta/online-shop-system/internal/modules/mailer/services"
)

type feedbackServiceImpl struct {
	repo   repository.FeedbackRepository
	mailer mailerServices.MailerService
}

func NewFeedbackService(repo repository.FeedbackRepository, mailer mailerServices.MailerService) services.FeedbackService {
	return &feedbackServiceImpl{repo: repo, mailer: mailer}
}

func (s *feedbackServiceImpl) Create(ctx context.Context, req *delivery.CreateFeedbackRequest) (*delivery.FeedbackResponse, error) {
	feedback := &models.Feedback{
		Name:    req.Name,
		Email:   req.Email,
		Subject: req.Subject,
		Message: req.Message,
	}

	if err := s.repo.WithTransaction(ctx, func(txCtx context.Context) error {
		return s.repo.Create(txCtx, feedback)
	}); err != nil {
		return nil, fmt.Errorf("failed to create feedback: %w", err)
	}

	// Best-effort acknowledgement email. Feedback is already persisted, so we
	// don't fail the request if the mail provider is down — but we surface
	// the failure in the response so the caller knows the email didn't ship.
	res := delivery.ToFeedbackResponse(feedback)
	if _, err := s.mailer.SendEmail(ctx, buildAckEmail(feedback)); err != nil {
		fmt.Printf("warning: failed to send feedback ack email to %s: %v\n", feedback.Email, err)
		res.EmailSent = false
		res.EmailError = err.Error()
	} else {
		res.EmailSent = true
	}

	return res, nil
}

func buildAckEmail(f *models.Feedback) mailerDelivery.SendEmailRequest {
	body := fmt.Sprintf(
		"Hi %s,\n\nThanks for reaching out — we've received your feedback and the team will get back to you shortly.\n\n"+
			"Subject: %s\n\nYour message:\n%s\n\n— Online Shop Volt",
		f.Name, f.Subject, f.Message,
	)
	return mailerDelivery.SendEmailRequest{
		To:      []string{f.Email},
		Subject: "We received your feedback: " + f.Subject,
		Body:    body,
		IsHTML:  false,
	}
}
