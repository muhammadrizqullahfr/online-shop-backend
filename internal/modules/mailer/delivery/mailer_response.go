package delivery

import "time"

type SendEmailResult struct {
	ID      string    `json:"id"`
	To      []string  `json:"to"`
	Subject string    `json:"subject"`
	SentAt  time.Time `json:"sentAt"`
}
