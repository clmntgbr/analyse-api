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

func NewMediaKey(userID uuid.UUID, filename string) string {
	ext := filepath.Ext(filename)
	return userID.String() + "/" + uuid.NewString() + ext
}
