package media

import (
	"fmt"
	"net/url"
	"path/filepath"
	"strings"
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

func DecodeObjectKey(key string) (string, error) {
	decoded, err := url.QueryUnescape(key)
	if err != nil {
		return "", fmt.Errorf("invalid media key: %w", err)
	}

	return decoded, nil
}

func UserIDFromKey(encodedKey string) (uuid.UUID, error) {
	key, err := DecodeObjectKey(encodedKey)
	if err != nil {
		return uuid.Nil, err
	}

	userID, _, ok := strings.Cut(key, "/")
	if !ok {
		return uuid.Nil, fmt.Errorf("invalid media key: %q", key)
	}

	return uuid.Parse(userID)
}
