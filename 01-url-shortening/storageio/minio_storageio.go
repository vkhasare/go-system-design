package storageio

import (
	"bytes"
	"context"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

// MinioFileStorageHandler is an implementation of FileStorageHandler using MinIO.
type MinioFileStorageHandler struct {
	client *minio.Client
}

// NewMinioFileStorageHandler creates a new MinioFileStorageHandler instance.
func NewMinioFileStorageHandler(endpoint, accessKey, secretKey string, useSSL bool) (FileStorageHandler, error) {
	// Initialize minio client
	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		return nil, err
	}

	return &MinioFileStorageHandler{
		client: minioClient,
	}, nil
}

func (m *MinioFileStorageHandler) UploadFile(ctx context.Context, bucketName, objectName string, data []byte, contentType string) error {
	// Ensure bucket exists. If not, create it (optional).
	exists, err := m.client.BucketExists(ctx, bucketName)
	if err != nil {
		return err
	}
	if !exists {
		err = m.client.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{})
		if err != nil {
			return err
		}
	}

	// Upload the file
	reader := bytes.NewReader(data)
	_, err = m.client.PutObject(ctx, bucketName, objectName, reader, int64(len(data)), minio.PutObjectOptions{
		ContentType: contentType,
	})
	return err
}
