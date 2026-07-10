package media

import (
	"context"
	"errors"
	"go-api/domain/entity"
	"go-api/domain/enum"
	"go-api/domain/repository"

	"github.com/google/uuid"
)

type CreateMediaUseCase struct {
	mediaRepo *repository.MediaRepository
}

func NewCreateMediaUseCase(mediaRepo *repository.MediaRepository) *CreateMediaUseCase {
	return &CreateMediaUseCase{mediaRepo: mediaRepo}
}

func (u *CreateMediaUseCase) Execute(ctx context.Context, userID uuid.UUID, key string, contentType string, size int64) (*entity.Media, error) {
	media := entity.Media{
		UserID:      userID,
		Key:         key,
		ContentType: contentType,
		Size:        size,
		Status:      enum.MediaStatusProcessing,
	}

	err := (*u.mediaRepo).Create(ctx, &media)
	if err != nil {
		return nil, errors.New("failed to create media")
	}

	return &media, nil
}
