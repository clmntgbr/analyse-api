package presenter

import (
	"go-api/domain/entity"
	"time"
)

type GeneratePresignedUploadUrlDetailResponse struct {
	UploadURL string `json:"uploadUrl"`
}

type SignalResponse struct {
	Name       string   `json:"name"`
	Score      int      `json:"score"`
	Confidence string   `json:"confidence"`
	Details    []string `json:"details"`
}

type MediaListResponse struct {
	ID         string           `json:"id"`
	Key        string           `json:"key"`
	Thumbnail  string           `json:"thumbnail"`
	Status     string           `json:"status"`
	FinalScore float64          `json:"finalScore,omitempty"`
	Confidence string           `json:"confidence,omitempty"`
	Verdict    string           `json:"verdict,omitempty"`
	Signals    []SignalResponse `json:"signals,omitempty"`
	CreatedAt  time.Time        `json:"createdAt"`
	UpdatedAt  time.Time        `json:"updatedAt"`
}

func NewGeneratePresignedUploadUrlDetailResponse(url string) GeneratePresignedUploadUrlDetailResponse {
	return GeneratePresignedUploadUrlDetailResponse{
		UploadURL: url,
	}
}

func NewMediaListResponse(media *entity.Media) *MediaListResponse {
	response := &MediaListResponse{
		ID:        media.ID.String(),
		Key:       media.Key,
		Thumbnail: thumbnailURL(media),
		Status:    string(media.Status),
		CreatedAt: media.CreatedAt,
		UpdatedAt: media.UpdatedAt,
	}

	if media.Verdict != "" {
		response.FinalScore = media.FinalScore
		response.Confidence = string(media.AnalysisConfidence)
		response.Verdict = media.Verdict
	}

	if len(media.Signals) > 0 {
		response.Signals = NewSignalResponses(media.Signals)
	}

	return response
}

func NewSignalResponses(signals []entity.Signal) []SignalResponse {
	responses := make([]SignalResponse, 0, len(signals))
	for _, signal := range signals {
		responses = append(responses, SignalResponse{
			Name:       signal.Name,
			Score:      signal.Score,
			Confidence: string(signal.Confidence),
			Details:    signal.Details,
		})
	}

	return responses
}

func thumbnailURL(media *entity.Media) string {
	if media.Thumbnail == "" {
		return ""
	}

	return "/api/medias/" + media.ID.String() + "/thumbnail"
}

func NewMediaListResponses(medias []*entity.Media) []*MediaListResponse {
	mediaListResponses := make([]*MediaListResponse, len(medias))
	for i, media := range medias {
		mediaListResponses[i] = NewMediaListResponse(media)
	}
	return mediaListResponses
}
