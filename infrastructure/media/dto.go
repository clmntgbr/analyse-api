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

var allowedExtensions = map[string]struct{}{
	"jpg":  {},
	"jpeg": {},
	"png":  {},
	"webp": {},
	"mp4":  {},
	"mov":  {},
	"avi":  {},
	"mkv":  {},
	"m4v":  {},
	"mpeg": {},
	"mpg":  {},
	"wmv":  {},
	"asf":  {},
	"flv":  {},
	"webm": {},
	"ogg":  {},
	"ogv":  {},
	"mka":  {},
}

func ValidatePresignUploadInput(input PresignUploadInput) error {
	ext := strings.ToLower(strings.TrimPrefix(filepath.Ext(input.Filename), "."))
	if ext == "" {
		return fmt.Errorf("filename must have a supported extension")
	}

	if _, ok := allowedExtensions[ext]; !ok {
		return fmt.Errorf("unsupported file type: .%s", ext)
	}

	return nil
}

func NewFileKey(filename string) string {
	ext := filepath.Ext(filename)
	return uuid.NewString() + ext
}

func NewObjectKey(userID uuid.UUID, fileKey string) string {
	return userID.String() + "/" + fileKey
}

func NewObjectKeyFromFilename(userID uuid.UUID, filename string) string {
	return NewObjectKey(userID, NewFileKey(filename))
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

func FileKeyFromObjectKey(encodedKey string) (string, error) {
	key, err := DecodeObjectKey(encodedKey)
	if err != nil {
		return "", err
	}

	_, fileKey, ok := strings.Cut(key, "/")
	if !ok {
		return "", fmt.Errorf("invalid media key: %q", key)
	}

	return fileKey, nil
}
