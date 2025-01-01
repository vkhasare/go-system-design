package storageio

import (
	"context"
)

// FileStorageHandler defines the interface for uploading files.
type FileStorageHandler interface {
	// UploadFile uploads data to the specified bucket with the given object name.
	// contentType is the MIME type (e.g., "image/png").
	UploadFile(ctx context.Context, bucketName, objectName string, data []byte, contentType string) error
}
