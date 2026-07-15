package presenter

import "go-api/domain/entity"

type MediaStatisticsResponse struct {
	AnalysesCount  int64   `json:"analysesCount"`
	RealImageCount int64   `json:"realImageCount"`
	AIImageCount   int64   `json:"aiImageCount"`
	AverageScore   float64 `json:"averageScore"`
}

func NewMediaStatisticsResponse(stats *entity.MediaStatistics) MediaStatisticsResponse {
	return MediaStatisticsResponse{
		AnalysesCount:  stats.AnalysesCount,
		RealImageCount: stats.RealImageCount,
		AIImageCount:   stats.AIImageCount,
		AverageScore:   stats.AverageScore,
	}
}
