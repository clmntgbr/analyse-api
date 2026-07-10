package media

import (
	"context"
	"errors"
	"go-api/domain/entity"
	"go-api/domain/repository"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type GetMediaUseCase struct {
	mediaRepo *repository.MediaRepository
}

func NewGetMediaUseCase(mediaRepo *repository.MediaRepository) *GetMediaUseCase {
	return &GetMediaUseCase{mediaRepo: mediaRepo}
}

func (u *GetMediaUseCase) Execute(ctx context.Context, userID uuid.UUID, mediaID uuid.UUID) (*entity.Media, error) {
	media, err := (*u.mediaRepo).GetByID(ctx, mediaID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("media not found")
		}

		return nil, errors.New("failed to get media")
	}

	if media.UserID != userID {
		return nil, errors.New("media not found")
	}

	return media, nil
}
