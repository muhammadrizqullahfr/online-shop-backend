package delivery

import "time"

type XenditCallbackRequest struct {
	ID                 string    `json:"id"`
	ExternalID         string    `json:"external_id"`
	Status             string    `json:"status"`
	Amount             float64   `json:"amount"`
	PaidAmount         float64   `json:"paid_amount"`
	PaidAt             time.Time `json:"paid_at"`
	Created            time.Time `json:"created"`
	Updated            time.Time `json:"updated"`
	UserID             string    `json:"user_id"`
	Currency           string    `json:"currency"`
	PaymentMethod      string    `json:"payment_method"`
	PaymentChannel     string    `json:"payment_channel"`
	PaymentDestination string    `json:"payment_destination"`
}
