package provider

import (
	"github.com/RakaMurdiarta/online-shop-system/internal/middlewares"
	"github.com/RakaMurdiarta/online-shop-system/internal/modules/orders/handlers"
	orderRepo "github.com/RakaMurdiarta/online-shop-system/internal/modules/orders/repository"
	orderService "github.com/RakaMurdiarta/online-shop-system/internal/modules/orders/services"
	"github.com/RakaMurdiarta/online-shop-system/pkg/database"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v5"
)

func OrderProvider(
	tm *database.TransactionManagerImpl,
	privateRoute *echo.Group,
	publicRoute *echo.Group,
	orderRepo orderRepo.OrderRepository,
	orderService orderService.OrderService,
	callbackService orderService.OrderCallbackService,
) {
	v := validator.New()

	handler := handlers.NewOrderHandler(orderService, v)
	callbackHandler := handlers.NewOrderCallbackHandler(callbackService)

	g := privateRoute.Group("/orders")
	g.POST("", handler.CreateOrder)
	g.GET("", handler.ListOrders, middlewares.IsAdmin)
	g.GET("/my-orders", handler.GetMyOrders)
	g.GET("/by-user/:userId", handler.GetOrdersByUserID, middlewares.IsAdmin)
	g.GET("/:id", handler.GetOrderByID)
	g.PUT("/:id", handler.UpdateOrder, middlewares.IsAdmin)
	g.DELETE("/:id", handler.DeleteOrder, middlewares.IsAdmin)

	g.GET("/:id/items", handler.ListItems)
	g.POST("/:id/items", handler.AddItem, middlewares.IsAdmin)
	g.PUT("/:id/items/:itemId", handler.UpdateItem, middlewares.IsAdmin)
	g.DELETE("/:id/items/:itemId", handler.DeleteItem, middlewares.IsAdmin)

	paymentGroup := publicRoute.Group("/payments")
	paymentGroup.POST("/cb", callbackHandler.Callback)
}
