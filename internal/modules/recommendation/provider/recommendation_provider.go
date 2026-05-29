package provider

import (
	"github.com/RakaMurdiarta/online-shop-system/internal/modules/recommendation/client"
	"github.com/RakaMurdiarta/online-shop-system/internal/modules/recommendation/handlers"
	"github.com/labstack/echo/v5"
	"gorm.io/gorm"
)

func RecommendationProvider(
	e *echo.Echo,
	db *gorm.DB,
	client client.RecommendationClient,
) {
	handler := handlers.NewRecommendationHandler(client, db)

	// Register directly on Echo router to match example and unit_7 paths perfectly
	group := e.Group("/ml-service")
	group.GET("/recommendations", handler.RecommendProductHandler)
	group.POST("/products", handler.CreateProductHandler)
	group.POST("/products/bulk", handler.BulkCreateProductHandler)
	group.POST("/track", handler.TrackInteractionHandler)
	group.POST("/chatbot", handler.ChatbotQueryHandler)
}
