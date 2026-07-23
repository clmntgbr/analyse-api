package presenter

import (
	"go-api/domain/entity"
	"time"
)

type QuotaResponse struct {
	ID                string        `json:"id"`
	MaxImagesPerMonth int           `json:"maxImagesPerMonth"`
	MaxVideosPerMonth int           `json:"maxVideosPerMonth"`
	MaxFileSizeImage  int64         `json:"maxFileSizeImage"`
	MaxFileSizeVideo  int64         `json:"maxFileSizeVideo"`
	FullPipeline      bool          `json:"fullPipeline"`
	HistoryRetention  time.Duration `json:"historyRetention"`
	CreatedAt         time.Time     `json:"createdAt"`
	UpdatedAt         time.Time     `json:"updatedAt"`
}

func NewQuotaResponse(quota *entity.Quota) *QuotaResponse {
	return &QuotaResponse{
		ID:                quota.ID.String(),
		MaxImagesPerMonth: quota.MaxImagesPerMonth,
		MaxVideosPerMonth: quota.MaxVideosPerMonth,
		MaxFileSizeImage:  quota.MaxFileSizeImage,
		MaxFileSizeVideo:  quota.MaxFileSizeVideo,
		FullPipeline:      quota.FullPipeline,
		HistoryRetention:  quota.HistoryRetention,
		CreatedAt:         quota.CreatedAt,
		UpdatedAt:         quota.UpdatedAt,
	}
}

func NewQuotaResponses(quotas []*entity.Quota) []*QuotaResponse {
	responses := make([]*QuotaResponse, 0, len(quotas))
	for _, quota := range quotas {
		responses = append(responses, NewQuotaResponse(quota))
	}
	return responses
}
