package services

import (
	"context"
	"mime/multipart"

	"github.com/RakaMurdiarta/online-shop-system/internal/modules/upload/delivery"
)

type UploadService interface {
	UploadImage(ctx context.Context, file *multipart.FileHeader) (*delivery.UploadImageResponse, error)
}
