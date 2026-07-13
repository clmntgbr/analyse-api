package rabbitmq

import (
	"github.com/google/uuid"
)

type MessagePayload struct {
	SecretKey string        `json:"secret_key"`
	Message   MetadataEvent `json:"message"`
}

type MetadataEvent struct {
	UserID       uuid.UUID `json:"user_id"`
	MediaID      uuid.UUID `json:"media_id"`
	MediaKey     string    `json:"media_key"`
	ThumbnailKey string    `json:"thumbnail_key"`
}
