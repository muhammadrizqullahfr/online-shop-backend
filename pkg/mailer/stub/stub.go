// Package stub provides an in-memory mailer.Transport that records sends
// without contacting any provider. Useful for local development and tests.
package stub

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/RakaMurdiarta/online-shop-system/pkg/mailer"
	"github.com/google/uuid"
)

type Transport struct {
	mu   sync.Mutex
	sent []mailer.Message
}

func New() *Transport {
	return &Transport{}
}

func (t *Transport) Send(_ context.Context, msg mailer.Message) (mailer.Receipt, error) {
	t.mu.Lock()
	t.sent = append(t.sent, msg)
	t.mu.Unlock()

	fmt.Printf("[mailer:stub] to=%v subject=%q bytes=%d\n", msg.To, msg.Subject, len(msg.Body))

	return mailer.Receipt{
		ID:      "stub-" + uuid.NewString(),
		To:      msg.To,
		Subject: msg.Subject,
		SentAt:  time.Now().UTC(),
	}, nil
}

// Sent returns a snapshot of every message handed to this transport.
// Handy for tests that want to assert on outbound mail.
func (t *Transport) Sent() []mailer.Message {
	t.mu.Lock()
	defer t.mu.Unlock()
	out := make([]mailer.Message, len(t.sent))
	copy(out, t.sent)
	return out
}
