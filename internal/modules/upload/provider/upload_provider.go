package provider

import (
	"github.com/RakaMurdiarta/online-shop-system/internal/modules/upload/handlers"
	uploadServiceImpl "github.com/RakaMurdiarta/online-shop-system/internal/modules/upload/services/Impl"
	"github.com/RakaMurdiarta/online-shop-system/pkg/shared"

	"github.com/labstack/echo/v5"
)

func UploadProvider(privateRoute *echo.Group, storageClient *shared.SupabaseStorageClient) {
	uploadService := uploadServiceImpl.NewUploadService(storageClient)
	uploadHandler := handlers.NewUploadHandler(uploadService)

	upload := privateRoute.Group("/uploads")
	upload.POST("/image", uploadHandler.UploadImage)
}
