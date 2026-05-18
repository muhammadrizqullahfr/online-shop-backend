package provider

import (
	userRepo "github.com/RakaMurdiarta/online-shop-system/internal/modules/users/repository"

	"github.com/RakaMurdiarta/online-shop-system/internal/config"
	"github.com/RakaMurdiarta/online-shop-system/internal/modules/products/handlers"
	categoryRepo "github.com/RakaMurdiarta/online-shop-system/internal/modules/products/repository"
	productService "github.com/RakaMurdiarta/online-shop-system/internal/modules/products/services"

	"github.com/RakaMurdiarta/online-shop-system/pkg/database"
	"github.com/RakaMurdiarta/online-shop-system/pkg/shared"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v5"
	"gorm.io/gorm"
)

func ProductProvider(
	db *gorm.DB,
	privateRoute *echo.Group,
	publicRoute *echo.Group,
	tx *database.TransactionManagerImpl,
	userRepo userRepo.UserRepository,
	categoryRepo categoryRepo.CategoryRepository,
	categoryService productService.CategoryService,
	productService productService.ProductService,
	conf *config.Config,
	storageClient *shared.SupabaseStorageClient,
) {
	v := validator.New()

	//categories
	categoryHandler := handlers.NewCategoryHandler(categoryService, v)
	handler := handlers.NewProductHandler(productService, v)

	categories := privateRoute.Group("/categories")
	categoriesPub := publicRoute.Group("/categories")
	categories.POST("", categoryHandler.CreateCategory)
	categoriesPub.GET("", categoryHandler.GetAllCategories)
	categories.GET("/tree", categoryHandler.GetCategoryTree)
	categories.GET("/slug/:slug", categoryHandler.GetCategoryBySlug)
	categoriesPub.GET("/:id", categoryHandler.GetCategoryByID)
	categories.PUT("/:id", categoryHandler.UpdateCategory)
	categories.DELETE("/:id", categoryHandler.DeleteCategory)

	publicGroup := publicRoute.Group("/products")
	publicGroup.GET("", handler.GetAll)
	publicGroup.GET("/:id", handler.GetByID)

	privateGroup := privateRoute.Group("/products")
	privateGroup.POST("", handler.Create)
	privateGroup.PUT("/:id", handler.Update)
	privateGroup.DELETE("/:id", handler.Delete)
}
