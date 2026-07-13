package rabbitmq

type MessagePayload struct {
	SecretKey string `json:"secret_key"`
	Message   Event  `json:"message"`
}

type Event interface {
	isMessage()
}

func (MetadataEvent) isMessage() {}

type MetadataEvent struct {
	UserID       string `json:"user_id"`
	MediaKey     string `json:"media_key"`
	ThumbnailKey string `json:"thumbnail_key"`
}
