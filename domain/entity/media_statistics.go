package entity

type MediaStatistics struct {
	AnalysesCount  int64   `json:"analyses_count"`
	RealImageCount int64   `json:"real_image_count"`
	AIImageCount   int64   `json:"ai_image_count"`
	AverageScore   float64 `json:"average_score"`
}
