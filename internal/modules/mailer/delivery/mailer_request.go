package delivery

type SendEmailRequest struct {
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
