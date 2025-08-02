package image_store

import (
    "context"
    "fmt"
    "os"
    "github.com/minio/minio-go/v7"
    "github.com/minio/minio-go/v7/pkg/credentials"
)

// NewMinioClientFromEnv creates a MinIO client and ensures the bucket exists
func NewMinioClientFromEnv() (*minio.Client, error) {
    endpoint := os.Getenv("MINIO_ENDPOINT")
    accessKeyID := os.Getenv("MINIO_ACCESS_KEY")
    secretAccessKey := os.Getenv("MINIO_SECRET_KEY")
    useSSL := os.Getenv("MINIO_USE_SSL") == "true"

    client, err := minio.New(endpoint, &minio.Options{
        Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
        Secure: useSSL,
    })
    if err != nil {
        return nil, err
    }

    // Ensure the bucket exists (auto-create if missing)
    bucket := os.Getenv("MINIO_BUCKET")
    if bucket == "" {
        bucket = "listing-images"
    }
    ctx := context.Background()
    exists, err := client.BucketExists(ctx, bucket)
    if err != nil {
        return nil, fmt.Errorf("could not check bucket: %w", err)
    }
    if !exists {
        err := client.MakeBucket(ctx, bucket, minio.MakeBucketOptions{})
        if err != nil {
            return nil, fmt.Errorf("could not create bucket: %w", err)
        }
        fmt.Println("Created bucket:", bucket)
    }

    return client, nil
}
