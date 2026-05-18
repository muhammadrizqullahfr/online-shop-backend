package impl

import (
	"context"
	"errors"
	"fmt"
	"mime/multipart"
	"path/filepath"
	"strings"
	"time"

	"github.com/RakaMurdiarta/online-shop-system/internal/modules/upload/delivery"
	"github.com/RakaMurdiarta/online-shop-system/internal/modules/upload/services"
	"github.com/RakaMurdiarta/online-shop-system/pkg/shared"
)

const (
	uploadBucket    = "e-commerce"
	uploadDir       = "uploads"
	maxUploadSize   = 5 * 1024 * 1024
	imageMimePrefix = "image/"
)

type uploadServiceImpl struct {
	storageClient *shared.SupabaseStorageClient
}

func NewUploadService(storageClient *shared.SupabaseStorageClient) services.UploadService {
	return &uploadServiceImpl{storageClient: storageClient}
}

func (s *uploadServiceImpl) UploadImage(ctx context.Context, file *multipart.FileHeader) (*delivery.UploadImageResponse, error) {
	if file.Size > maxUploadSize {
		return nil, fmt.Errorf("file too large: max %d bytes", maxUploadSize)
	}

	contentType := file.Header.Get("Content-Type")
	if !strings.HasPrefix(contentType, imageMimePrefix) {
		return nil, errors.New("invalid file type: only images are allowed")
	}

	src, err := file.Open()
	if err != nil {
		return nil, fmt.Errorf("failed to open uploaded file: %w", err)
	}
	defer func() {
		if cerr := src.Close(); cerr != nil {
			fmt.Printf("warning: failed to close upload file: %v\n", cerr)
		}
	}()

	ext := filepath.Ext(file.Filename)
	objectPath := fmt.Sprintf("%s/%d_%s%s", uploadDir, time.Now().UnixNano(), strings.TrimSuffix(filepath.Base(file.Filename), ext), ext)

	if _, err := s.storageClient.UploadFile(uploadBucket, objectPath, src, contentType); err != nil {
		return nil, fmt.Errorf("failed to upload file to storage: %w", err)
	}

	return &delivery.UploadImageResponse{
		URL: s.storageClient.GetPublicURL(uploadBucket, objectPath),
	}, nil
}
