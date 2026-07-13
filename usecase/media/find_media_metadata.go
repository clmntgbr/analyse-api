package media

import (
	"context"
	"errors"
	"go-api/domain/repository"
	"go-api/infrastructure/config"
	"go-api/infrastructure/messaging/rabbitmq"

	"github.com/google/uuid"
)

type FindMediaMetadataUseCase struct {
	mediaRepo *repository.MediaRepository
	publisher rabbitmq.Publisher
	config    *config.Config
}

func NewFindMediaMetadataUseCase(
	mediaRepo *repository.MediaRepository,
	publisher rabbitmq.Publisher,
	config *config.Config,
) *FindMediaMetadataUseCase {
	return &FindMediaMetadataUseCase{
		mediaRepo: mediaRepo,
		publisher: publisher,
		config:    config,
	}
}

func (u *FindMediaMetadataUseCase) Execute(ctx context.Context, mediaID uuid.UUID) error {
	media, err := (*u.mediaRepo).GetByID(ctx, mediaID)
	if err != nil {
		return errors.New("failed to get media")
	}

	event := rabbitmq.MetadataEvent{
		UserID:       media.UserID.String(),
		MediaKey:     media.Key,
		ThumbnailKey: media.Thumbnail,
	}

	err = u.publisher.PublishMetadataEvent(ctx, u.config, event)
	if err != nil {
		return errors.New("failed to publish metadata event")
	}

	return nil
}
