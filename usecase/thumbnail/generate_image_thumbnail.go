package thumbnail

import (
	"bytes"
	"context"
	"errors"
	"go-api/infrastructure/storage"
	"io"

	"github.com/disintegration/imaging"
)

type GenerateImageThumbnailUseCase struct {
	storage *storage.MinIOStorage
}

func NewGenerateImageThumbnailUseCase(storage *storage.MinIOStorage) *GenerateImageThumbnailUseCase {
	return &GenerateImageThumbnailUseCase{storage: storage}
}

func (uc *GenerateImageThumbnailUseCase) Execute(ctx context.Context, src io.Reader, maxWidth int) ([]byte, error) {
	img, err := imaging.Decode(src, imaging.AutoOrientation(true))
	if err != nil {
		return nil, errors.New("failed to decode image")
	}

	thumb := imaging.Resize(img, maxWidth, 0, imaging.Lanczos)

	var buf bytes.Buffer
	if err := imaging.Encode(&buf, thumb, imaging.JPEG, imaging.JPEGQuality(80)); err != nil {
		return nil, errors.New("failed to encode thumbnail")
	}

	return buf.Bytes(), nil
}
