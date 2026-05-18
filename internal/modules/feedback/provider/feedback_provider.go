package provider

import (
	"github.com/RakaMurdiarta/online-shop-system/internal/modules/feedback/handlers"
	feedbackRepoImpl "github.com/RakaMurdiarta/online-shop-system/internal/modules/feedback/repository/Impl"
	feedbackServiceImpl "github.com/RakaMurdiarta/online-shop-system/internal/modules/feedback/services/Impl"
	mailerServices "github.com/RakaMurdiarta/online-shop-system/internal/modules/mailer/services"
	"github.com/RakaMurdiarta/online-shop-system/pkg/database"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v5"
)

func FeedbackProvider(
	tm *database.TransactionManagerImpl,
	publicRoute *echo.Group,
	mailer mailerServices.MailerService,
) {
	v := validator.New()

	repo := feedbackRepoImpl.NewFeedbackRepository(tm)
	service := feedbackServiceImpl.NewFeedbackService(repo, mailer)
	handler := handlers.NewFeedbackHandler(service, v)

	publicGroup := publicRoute.Group("/feedbacks")
	publicGroup.POST("", handler.Create)
}
