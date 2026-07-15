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

type InsightResponse struct {
	Noise       float64 `json:"noise"`
	Compression float64 `json:"compression"`
	Frequency   float64 `json:"frequency"`
	Histogram   float64 `json:"histogram"`
}

type MediaListResponse struct {
	ID         string    `json:"id"`
	Key        string    `json:"key"`
	Thumbnail  string    `json:"thumbnail"`
	Status     string    `json:"status"`
	FinalScore float64   `json:"finalScore,omitempty"`
	Confidence string    `json:"confidence,omitempty"`
	Verdict    string    `json:"verdict,omitempty"`
	Size       int64     `json:"size,omitempty"`
	CreatedAt  time.Time `json:"createdAt"`
	UpdatedAt  time.Time `json:"updatedAt"`
}

type MediaDetailResponse struct {
	ID         string           `json:"id"`
	Key        string           `json:"key"`
	Thumbnail  string           `json:"thumbnail"`
	Status     string           `json:"status"`
	FinalScore float64          `json:"finalScore,omitempty"`
	Confidence string           `json:"confidence,omitempty"`
	Verdict    string           `json:"verdict,omitempty"`
	Signals    []SignalResponse `json:"signals,omitempty"`
	Insight    InsightResponse  `json:"insight,omitempty"`
	Size       int64            `json:"size,omitempty"`
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
		Size:      media.Size,
		CreatedAt: media.CreatedAt,
		UpdatedAt: media.UpdatedAt,
	}

	if media.Verdict != "" {
		response.FinalScore = media.FinalScore
		response.Confidence = string(media.AnalysisConfidence)
		response.Verdict = media.Verdict
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

func NewInsightResponse(insight *entity.Insight) InsightResponse {
	return InsightResponse{
		Noise:       insight.Noise,
		Compression: insight.Compression,
		Frequency:   insight.Frequency,
		Histogram:   insight.Histogram,
	}
}

func NewMediaDetailResponse(media *entity.Media) *MediaDetailResponse {
	return &MediaDetailResponse{
		ID:         media.ID.String(),
		Key:        media.Key,
		Thumbnail:  thumbnailURL(media),
		Status:     string(media.Status),
		FinalScore: media.FinalScore,
		Confidence: string(media.AnalysisConfidence),
		Verdict:    media.Verdict,
		Signals:    NewSignalResponses(media.Signals),
		Insight:    NewInsightResponse(media.Insight),
		Size:       media.Size,
		CreatedAt:  media.CreatedAt,
		UpdatedAt:  media.UpdatedAt,
	}
}

func NewMediaDetailResponses(medias []*entity.Media) []*MediaDetailResponse {
	mediaDetailResponses := make([]*MediaDetailResponse, len(medias))
	for i, media := range medias {
		mediaDetailResponses[i] = NewMediaDetailResponse(media)
	}
	return mediaDetailResponses
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
