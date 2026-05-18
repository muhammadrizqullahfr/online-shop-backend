package shared

import "github.com/xendit/xendit-go/v7"

type XenditClient struct {
	API *xendit.APIClient
}

func NewXenditClient(secretKey string) *XenditClient {
	return &XenditClient{
		API: xendit.NewClient(secretKey),
	}
}
