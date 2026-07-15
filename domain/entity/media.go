package entity

import (
	"go-api/domain/enum"
	"time"

	"github.com/google/uuid"
)

type Media struct {
	ID uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`

	UserID uuid.UUID `gorm:"type:uuid;not null" json:"user_id"`
	User   User      `gorm:"foreignKey:UserID" json:"user"`

	InsightID *uuid.UUID `gorm:"type:uuid;default:null;index:idx_media_insight_id" json:"insight_id"`
	Insight   *Insight   `gorm:"foreignKey:InsightID" json:"insight"`

	Key         string `gorm:"uniqueIndex;not null" json:"key"`
	Thumbnail   string `gorm:"not null" json:"thumbnail"`
	ContentType string `gorm:"not null" json:"content_type"`
	Size        int64  `gorm:"not null" json:"size"`

	Signals []Signal `gorm:"foreignKey:MediaID" json:"signals"`

	FinalScore         float64         `gorm:"default:-1" json:"final_score"`
	AnalysisConfidence ConfidenceLevel `gorm:"type:varchar(20);default:'unknown'" json:"confidence"`
	Verdict            string          `gorm:"type:varchar(20);default:''" json:"verdict"`

	Status   enum.MediaStatus   `gorm:"type:varchar(20);not null;check:status IN ('processing','uploaded','analyzed');index:idx_media_status" json:"status"`
	Statuses []enum.MediaStatus `gorm:"serializer:json;type:jsonb;default:'[]'" json:"statuses"`

	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

func (Media) TableName() string {
	return "medias"
}
