package media

import (
	"context"
	mediadto "go-api/infrastructure/media"
	"go-api/infrastructure/storage"
	"time"

	"github.com/google/uuid"
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

func (uc *GeneratePresignedUploadUrlUseCase) Execute(ctx context.Context, userID uuid.UUID, input mediadto.PresignUploadInput) (string, error) {
	objectKey := mediadto.NewObjectKeyFromFilename(userID, input.Filename)

	url, err := uc.storage.PresignedPutURL(ctx, objectKey, 15*time.Minute)
	if err != nil {
		return "", err
	}

	return url, nil
}
