package s3

import (
	"context"
	"fmt"
	"io"
	"net/url"
	"strings"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"

	"github.com/optivor/optivor/internal/config"
	"github.com/optivor/optivor/internal/storage"
	"go.opentelemetry.io/otel"
)

// Driver implements storage.StorageDriver for S3-compatible object storage using minio-go.
type Driver struct {
	client *minio.Client
	bucket string
}

// New initializes a new S3 StorageDriver instance.
func New(cfg config.S3Config) (*Driver, error) {
	endpoint := cfg.Endpoint
	secure := true

	if strings.HasPrefix(endpoint, "http://") {
		secure = false
		endpoint = strings.TrimPrefix(endpoint, "http://")
	} else if strings.HasPrefix(endpoint, "https://") {
		endpoint = strings.TrimPrefix(endpoint, "https://")
	}

	if u, err := url.Parse("http://" + endpoint); err == nil && u.Host != "" {
		endpoint = u.Host
	}

	opts := &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.AccessKeyID, cfg.SecretAccessKey, ""),
		Secure: secure,
		Region: cfg.Region,
	}

	client, err := minio.New(endpoint, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize S3 client: %w", err)
	}

	return &Driver{
		client: client,
		bucket: cfg.Bucket,
	}, nil
}

// Get fetches an object from S3-compatible storage by key.
func (d *Driver) Get(ctx context.Context, key string) (io.ReadCloser, error) {
	ctx, span := otel.Tracer("optivor").Start(ctx, "storage.GetObject")
	defer span.End()
	// Trim leading slash if present
	cleanKey := strings.TrimPrefix(key, "/")

	obj, err := d.client.GetObject(ctx, d.bucket, cleanKey, minio.GetObjectOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to get object from S3: %w", err)
	}

	// Verify existence with Stat() call
	_, err = obj.Stat()
	if err != nil {
		_ = obj.Close()
		errResp := minio.ToErrorResponse(err)
		if errResp.Code == "NoSuchKey" || errResp.Code == "NotFound" || errResp.StatusCode == 404 {
			return nil, storage.ErrNotFound
		}
		return nil, fmt.Errorf("failed to stat object in S3: %w", err)
	}

	return obj, nil
}
