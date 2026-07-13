package media

import (
	"context"
	"errors"
	"go-api/domain/enum"
	"go-api/domain/repository"

	"github.com/google/uuid"
)

type UpdateMediaStatusUseCase struct {
	mediaRepo *repository.MediaRepository
}

func NewUpdateMediaStatusUseCase(mediaRepo *repository.MediaRepository) *UpdateMediaStatusUseCase {
	return &UpdateMediaStatusUseCase{mediaRepo: mediaRepo}
}

func (u *UpdateMediaStatusUseCase) Execute(ctx context.Context, mediaID uuid.UUID, status enum.MediaStatus) error {
	media, err := (*u.mediaRepo).GetByID(ctx, mediaID)
	if err != nil {
		return errors.New("failed to get media")
	}

	media.Statuses = append(media.Statuses, status)
	media.Status = media.Statuses[len(media.Statuses)-1]

	return (*u.mediaRepo).Update(ctx, media)
}
