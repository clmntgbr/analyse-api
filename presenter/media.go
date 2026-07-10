package presenter

type MediaDetailResponse struct {
	UploadURL string `json:"uploadUrl"`
}

func NewMediaDetailResponse(url string) MediaDetailResponse {
	return MediaDetailResponse{
		UploadURL: url,
	}
}
