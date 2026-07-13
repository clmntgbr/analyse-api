package rabbitmq

import "github.com/google/uuid"

type MessagePayload struct {
	SecretKey string `json:"secret_key"`
	Message   any    `json:"message"`
}

type AnalyzeMessage struct {
	UserID       uuid.UUID `json:"user_id"`
	MediaID      uuid.UUID `json:"media_id"`
	MediaKey     string    `json:"media_key"`
	ThumbnailKey string    `json:"thumbnail_key"`
}

type StageDoneMessage struct {
	MediaID uuid.UUID `json:"media_id"`
	Stage   string    `json:"stage"`
}

type FailedMessage struct {
	MediaID uuid.UUID `json:"media_id"`
	Stage   string    `json:"stage"`
	Error   string    `json:"error"`
}
