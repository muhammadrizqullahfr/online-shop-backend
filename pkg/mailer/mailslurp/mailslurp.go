// Package mailslurp implements mailer.Transport on top of the MailSlurp SDK.
package mailslurp

import (
	"context"
	"errors"
	"fmt"

	"github.com/RakaMurdiarta/online-shop-system/pkg/mailer"
	mailslurpsdk "github.com/mailslurp/mailslurp-client-go"
)

type Transport struct {
	api     *mailslurpsdk.APIClient
	apiKey  string
	inboxID string
}

func New(apiKey, inboxID string) *Transport {
	cfg := mailslurpsdk.NewConfiguration()
	cfg.AddDefaultHeader("x-api-key", apiKey)
	return &Transport{
		api:     mailslurpsdk.NewAPIClient(cfg),
		apiKey:  apiKey,
		inboxID: inboxID,
	}
}

func (t *Transport) Send(ctx context.Context, msg mailer.Message) (mailer.Receipt, error) {
	if t.inboxID == "" {
		return mailer.Receipt{}, errors.New("mailslurp: inbox id is required (set MAILER_INBOX_ID)")
	}
	if len(msg.To) == 0 {
		return mailer.Receipt{}, errors.New("mailslurp: at least one recipient is required")
	}

	opts := mailslurpsdk.SendEmailOptions{
		To:      &msg.To,
		Subject: &msg.Subject,
		Body:    &msg.Body,
		IsHTML:  &msg.IsHTML,
	}
	if len(msg.Cc) > 0 {
		opts.Cc = &msg.Cc
	}
	if len(msg.Bcc) > 0 {
		opts.Bcc = &msg.Bcc
	}
	if msg.From != "" {
		opts.From = &msg.From
	}
	if msg.FromName != "" {
		opts.FromName = &msg.FromName
	}
	if msg.ReplyTo != "" {
		opts.ReplyTo = &msg.ReplyTo
	}

	authCtx := context.WithValue(ctx, mailslurpsdk.ContextAPIKey, mailslurpsdk.APIKey{Key: t.apiKey})

	sent, _, err := t.api.InboxControllerApi.SendEmailAndConfirm(authCtx, t.inboxID, opts)
	if err != nil {
		return mailer.Receipt{}, fmt.Errorf("mailslurp: send failed: %w", err)
	}

	receipt := mailer.Receipt{
		ID:     sent.Id,
		SentAt: sent.SentAt,
	}
	if sent.To != nil {
		receipt.To = *sent.To
	}
	if sent.Subject != nil {
		receipt.Subject = *sent.Subject
	}
	return receipt, nil
}
