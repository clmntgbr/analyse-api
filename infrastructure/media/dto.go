package media

import (
	"path/filepath"
	"time"

	"github.com/google/uuid"
)

type MediaStatus string

const (
	MediaStatusPending  MediaStatus = "pending"
	MediaStatusUploaded MediaStatus = "uploaded"
	MediaStatusAnalyzed MediaStatus = "analyzed"
	MediaStatusFailed   MediaStatus = "failed"
)

type Media struct {
	ID          string
	Key         string
	ContentType string
	Status      MediaStatus
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type PresignUploadInput struct {
	Filename    string
	ContentType string
}

type PresignUploadOutput struct {
	MediaID   string
	UploadURL string
}

type ConfirmUploadInput struct {
	MediaID string
}

type ConfirmUploadOutput struct {
	MediaID string
	Status  MediaStatus
}

func NewMediaKey(filename string) string {
	ext := filepath.Ext(filename)
	return "media/" + uuid.NewString() + ext
}

func NewPendingMedia(key, contentType string) *Media {
	now := time.Now()
	return &Media{
		ID:          uuid.NewString(),
		Key:         key,
		ContentType: contentType,
		Status:      MediaStatusPending,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}
