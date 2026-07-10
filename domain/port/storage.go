package port

import (
	"context"
	"io"
	"time"
)

type StoredObject struct {
	Key         string
	Size        int64
	ContentType string
}

type Storage interface {
	Put(ctx context.Context, key string, r io.Reader, size int64, contentType string) (*StoredObject, error)
	Get(ctx context.Context, key string) (io.ReadCloser, error)
	Delete(ctx context.Context, key string) error
	PresignedPutURL(ctx context.Context, key string, expiry time.Duration) (string, error)
	PresignedGetURL(ctx context.Context, key string, expiry time.Duration) (string, error)
}