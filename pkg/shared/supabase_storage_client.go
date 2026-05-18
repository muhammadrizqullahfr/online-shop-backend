package shared

import (
	"fmt"
	"io"

	"github.com/RakaMurdiarta/online-shop-system/internal/config"
	storage_go "github.com/supabase-community/storage-go"
)

type SupabaseStorageClient struct {
	client *storage_go.Client
}

// NewSupabaseStorageClient initializes a new storage client using the project configuration
func NewSupabaseStorageClient(config *config.Config) *SupabaseStorageClient {
	storageUrl := fmt.Sprintf("%v/storage/v1", config.SupabaseURL)
	storageClient := storage_go.NewClient(storageUrl, config.SupabaseKey, nil)
	return &SupabaseStorageClient{client: storageClient}
}

// UploadFile uploads a file stream to a specific bucket
// bucketName: The target bucket (e.g., "products")
// path: The destination path inside the bucket (e.g., "images/item-01.png")
// file: The file stream (usually from echo.Context.FormFile)
// contentType: The MIME type of the file (e.g., "image/jpeg")
func (s *SupabaseStorageClient) UploadFile(bucketName string, path string, file io.Reader, contentType string) (string, error) {
	resp, err := s.client.UploadFile(bucketName, path, file, storage_go.FileOptions{
		ContentType: &contentType,
		Upsert:      boolPtr(true), // Overwrites the file if the path already exists
	})
	if err != nil {
		return "", err
	}

	return resp.Key, nil
}

// DeleteFile removes one or multiple files from a bucket
func (s *SupabaseStorageClient) DeleteFile(bucketName string, filePaths []string) error {
	_, err := s.client.RemoveFile(bucketName, filePaths)
	if err != nil {
		return err
	}
	return nil
}

// GetPublicURL generates a permanent link for files in a Public Bucket
func (s *SupabaseStorageClient) GetPublicURL(bucketName string, path string) string {
	resp := s.client.GetPublicUrl(bucketName, path)
	return resp.SignedURL
}

// GetSignedURL generates a temporary link for files in a Private Bucket
// expiresIn: Time in seconds before the link expires
func (s *SupabaseStorageClient) GetSignedURL(bucketName string, path string, expiresIn int) (string, error) {
	resp, err := s.client.CreateSignedUrl(bucketName, path, expiresIn)
	if err != nil {
		return "", err
	}
	return resp.SignedURL, nil
}

// boolPtr is a helper function to return a pointer to a boolean
func boolPtr(b bool) *bool {
	return &b
}
