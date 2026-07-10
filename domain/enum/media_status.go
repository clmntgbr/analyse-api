package enum

type MediaStatus string

const (
	MediaStatusPending    MediaStatus = "pending"
	MediaStatusProcessing MediaStatus = "processing"
	MediaStatusCompleted  MediaStatus = "completed"
	MediaStatusFailed     MediaStatus = "failed"
)
