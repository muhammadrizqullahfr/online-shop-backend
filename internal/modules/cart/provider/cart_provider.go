package provider

import (
	"github.com/RakaMurdiarta/online-shop-system/internal/modules/cart/handlers"
	cartRepo "github.com/RakaMurdiarta/online-shop-system/internal/modules/cart/repository"
	productRepo "github.com/RakaMurdiarta/online-shop-system/internal/modules/products/repository"

	cartService "github.com/RakaMurdiarta/online-shop-system/internal/modules/cart/services"
	"github.com/RakaMurdiarta/online-shop-system/pkg/database"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v5"
)

func CartProvider(
	tm *database.TransactionManagerImpl,
	privateRoute *echo.Group,
	productRepo productRepo.ProductRepository,
	cartRepo cartRepo.CartRepository,
	cartService cartService.CartService,
) {
	v := validator.New()

	handler := handlers.NewCartHandler(cartService, v)

	g := privateRoute.Group("/cart")
	g.GET("", handler.GetCart)
	g.DELETE("", handler.ClearCart)
	g.POST("/items", handler.AddItem)
	g.PUT("/items/:itemId", handler.UpdateItemQuantity)
	g.DELETE("/items/:itemId", handler.DeleteItem)
	g.DELETE("/items/clear", handler.ClearCart)
}
