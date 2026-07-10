package presenter

import (
	"go-api/domain/entity"
	"time"
)

type GeneratePresignedUploadUrlDetailResponse struct {
	UploadURL string `json:"uploadUrl"`
}

type MediaListResponse struct {
	ID        string    `json:"id"`
	Key       string    `json:"key"`
	Thumbnail string    `json:"thumbnail"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func NewGeneratePresignedUploadUrlDetailResponse(url string) GeneratePresignedUploadUrlDetailResponse {
	return GeneratePresignedUploadUrlDetailResponse{
		UploadURL: url,
	}
}

func NewMediaListResponse(media *entity.Media) *MediaListResponse {
	return &MediaListResponse{
		ID:        media.ID.String(),
		Key:       media.Key,
		Thumbnail: thumbnailURL(media),
		CreatedAt: media.CreatedAt,
		UpdatedAt: media.UpdatedAt,
	}
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
