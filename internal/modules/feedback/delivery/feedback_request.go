package delivery

type CreateFeedbackRequest struct {
	Name    string `json:"name" validate:"required,max=255"`
	Email   string `json:"email" validate:"required,email,max=255"`
	Subject string `json:"subject" validate:"required,max=255"`
	Message string `json:"message" validate:"required"`
}
