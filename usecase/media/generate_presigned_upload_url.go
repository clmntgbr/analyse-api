package media

import (
	"context"
	mediadto "go-api/infrastructure/media"
	"go-api/infrastructure/storage"
	"time"
)

type GeneratePresignedUploadUrlUseCase struct {
	storage *storage.MinIOStorage
}

func NewGeneratePresignedUploadUrlUseCase(
	storage *storage.MinIOStorage,
) *GeneratePresignedUploadUrlUseCase {
	return &GeneratePresignedUploadUrlUseCase{
		storage: storage,
	}
}

func (uc *GeneratePresignedUploadUrlUseCase) Execute(ctx context.Context, input mediadto.PresignUploadInput) (string, error) {
	key := mediadto.NewMediaKey(input.Filename)

	url, err := uc.storage.PresignedPutURL(ctx, key, 15*time.Minute)
	if err != nil {
		return "", err
	}

	return url, nil
}
