package entity

import (
	"github.com/google/uuid"
)

type Insight struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	Noise       float64   `json:"noise"`
	Compression float64   `json:"compression"`
	Frequency   float64   `json:"frequency"`
	Histogram   float64   `json:"histogram"`
}

func (Insight) TableName() string {
	return "insights"
}
