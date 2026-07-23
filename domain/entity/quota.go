package entity

import (
	"time"

	"github.com/google/uuid"
)

type Quota struct {
	ID uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`

	MaxImagesPerMonth int           `json:"max_images_per_month"`
	MaxVideosPerMonth int           `json:"max_videos_per_month"`
	MaxFileSizeImage  int64         `json:"max_file_size_image"`
	MaxFileSizeVideo  int64         `json:"max_file_size_video"`
	FullPipeline      bool          `json:"full_pipeline"`
	HistoryRetention  time.Duration `json:"history_retention"`

	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

func (Quota) TableName() string {
	return "quotas"
}
