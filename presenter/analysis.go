package presenter

import (
	"go-api/domain/entity"
	"time"
)

type GeneratePresignedUploadUrlDetailResponse struct {
	UploadURL string `json:"uploadUrl"`
}

type AnalysisListResponse struct {
	ID         string              `json:"id"`
	Status     string              `json:"status"`
	FinalScore float64             `json:"finalScore,omitempty"`
	Confidence string              `json:"confidence,omitempty"`
	Verdict    string              `json:"verdict,omitempty"`
	Filename   string              `json:"filename,omitempty"`
	Thumbnail  string              `json:"thumbnail,omitempty"`
	Medias     []MediaItemResponse `json:"medias"`
	CreatedAt  time.Time           `json:"createdAt"`
	UpdatedAt  time.Time           `json:"updatedAt"`
}

type AnalysisDetailResponse struct {
	ID         string              `json:"id"`
	Status     string              `json:"status"`
	FinalScore float64             `json:"finalScore,omitempty"`
	Confidence string              `json:"confidence,omitempty"`
	Verdict    string              `json:"verdict,omitempty"`
	Medias     []MediaItemResponse `json:"medias"`
	CreatedAt  time.Time           `json:"createdAt"`
	UpdatedAt  time.Time           `json:"updatedAt"`
}

func NewGeneratePresignedUploadUrlDetailResponse(url string) GeneratePresignedUploadUrlDetailResponse {
	return GeneratePresignedUploadUrlDetailResponse{
		UploadURL: url,
	}
}

func primaryMedia(analysis *entity.Analysis) *entity.Media {
	if analysis == nil || len(analysis.Medias) == 0 {
		return nil
	}
	return &analysis.Medias[0]
}

func analysisStatus(analysis *entity.Analysis) string {
	media := primaryMedia(analysis)
	if media == nil {
		return ""
	}
	return string(media.Status)
}

func NewAnalysisListResponse(analysis *entity.Analysis) *AnalysisListResponse {
	response := &AnalysisListResponse{
		ID:        analysis.ID.String(),
		Status:    analysisStatus(analysis),
		Medias:    NewMediaItemResponses(analysis.Medias),
		CreatedAt: analysis.CreatedAt,
		UpdatedAt: analysis.UpdatedAt,
	}

	if media := primaryMedia(analysis); media != nil {
		response.Filename = media.Filename
		response.Thumbnail = thumbnailURL(*media)
	}

	if analysis.Verdict != "" {
		response.FinalScore = analysis.FinalScore
		response.Confidence = string(analysis.AnalysisConfidence)
		response.Verdict = analysis.Verdict
	}

	return response
}

func NewAnalysisDetailResponse(analysis *entity.Analysis) *AnalysisDetailResponse {
	response := &AnalysisDetailResponse{
		ID:        analysis.ID.String(),
		Status:    analysisStatus(analysis),
		Medias:    NewMediaItemResponses(analysis.Medias),
		CreatedAt: analysis.CreatedAt,
		UpdatedAt: analysis.UpdatedAt,
	}

	if analysis.Verdict != "" {
		response.FinalScore = analysis.FinalScore
		response.Confidence = string(analysis.AnalysisConfidence)
		response.Verdict = analysis.Verdict
	}

	return response
}

func NewAnalysisListResponses(analyses []*entity.Analysis) []*AnalysisListResponse {
	responses := make([]*AnalysisListResponse, len(analyses))
	for i, analysis := range analyses {
		responses[i] = NewAnalysisListResponse(analysis)
	}
	return responses
}
