// Package mailer defines a provider-agnostic email transport abstraction.
// Concrete providers (MailSlurp, SendGrid, SMTP, ...) live in subpackages
// and implement Transport so consumers depend on this interface only.
package mailer

import (
	"context"
	"time"
)

type Transport interface {
	Send(ctx context.Context, msg Message) (Receipt, error)
}

type Message struct {
	To       []string
	Cc       []string
	Bcc      []string
	From     string
	FromName string
	ReplyTo  string
	Subject  string
	Body     string
	IsHTML   bool
}

type Receipt struct {
	ID      string
	To      []string
	Subject string
	SentAt  time.Time
}
