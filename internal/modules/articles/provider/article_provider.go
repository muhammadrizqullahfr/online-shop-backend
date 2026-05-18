package provider

import (
	"github.com/RakaMurdiarta/online-shop-system/internal/middlewares"
	"github.com/RakaMurdiarta/online-shop-system/internal/modules/articles/handlers"
	articleRepoImpl "github.com/RakaMurdiarta/online-shop-system/internal/modules/articles/repository/Impl"
	articleServiceImpl "github.com/RakaMurdiarta/online-shop-system/internal/modules/articles/services/Impl"
	"github.com/RakaMurdiarta/online-shop-system/pkg/database"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v5"
)

func ArticleProvider(
	tm *database.TransactionManagerImpl,
	privateRoute *echo.Group,
	publicRoute *echo.Group,
) {
	v := validator.New()

	repo := articleRepoImpl.NewArticleRepository(tm)
	service := articleServiceImpl.NewArticleService(repo)
	handler := handlers.NewArticleHandler(service, v)

	publicGroup := publicRoute.Group("/articles")
	publicGroup.GET("", handler.GetAll)
	publicGroup.GET("/:id", handler.GetByID)

	privateGroup := privateRoute.Group("/articles")
	privateGroup.POST("", handler.Create, middlewares.IsAdmin)
	privateGroup.PUT("/:id", handler.Update, middlewares.IsAdmin)
	privateGroup.DELETE("/:id", handler.Delete, middlewares.IsAdmin)
}
