package media

import (
	"bytes"
	"context"
	"errors"
	"go-api/domain/repository"
	"go-api/infrastructure/storage"
	"go-api/usecase/thumbnail"
	"strings"

	"github.com/google/uuid"
)

type GenerateThumbnailUseCase struct {
	storage                  *storage.MinIOStorage
	mediaRepo                *repository.MediaRepository
	generateThumbnailUseCase *thumbnail.GenerateImageThumbnailUseCase
}

func NewGenerateThumbnailUseCase(storage *storage.MinIOStorage, mediaRepo *repository.MediaRepository, generateThumbnailUseCase *thumbnail.GenerateImageThumbnailUseCase) *GenerateThumbnailUseCase {
	return &GenerateThumbnailUseCase{storage: storage, mediaRepo: mediaRepo, generateThumbnailUseCase: generateThumbnailUseCase}
}

func (uc *GenerateThumbnailUseCase) Execute(ctx context.Context, userID uuid.UUID, mediaID uuid.UUID) error {
	media, err := (*uc.mediaRepo).GetByID(ctx, mediaID)
	if err != nil {
		return errors.New("media not found")
	}

	original, err := uc.storage.Get(ctx, userID.String()+"/"+media.Key)
	if err != nil {
		return errors.New("failed to fetch original")
	}
	defer original.Close()

	var thumbBytes []byte
	if strings.HasPrefix(media.ContentType, "image/") {
		thumbBytes, err = uc.generateThumbnailUseCase.Execute(ctx, original, 400)
	} else {
		return errors.New("unsupported content type")
	}
	if err != nil {
		return err
	}

	thumbKey := userID.String() + "/thumbnails/" + media.ID.String() + ".jpg"
	if err := uc.storage.Put(ctx, thumbKey, bytes.NewReader(thumbBytes), int64(len(thumbBytes)), "image/jpeg"); err != nil {
		return errors.New("failed to store thumbnail")
	}

	media.Thumbnail = thumbKey
	return (*uc.mediaRepo).Update(ctx, media)
}
