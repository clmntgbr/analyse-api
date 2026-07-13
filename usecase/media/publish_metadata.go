package media

import (
	"context"
	"errors"
	"go-api/domain/repository"
	"go-api/infrastructure/config"
	"go-api/infrastructure/messaging/rabbitmq"

	"github.com/google/uuid"
)

type PublishMetadataUseCase struct {
	mediaRepo *repository.MediaRepository
	publisher rabbitmq.Publisher
	config    *config.Config
}

func NewPublishMetadataUseCase(
	mediaRepo *repository.MediaRepository,
	publisher rabbitmq.Publisher,
	config *config.Config,
) *PublishMetadataUseCase {
	return &PublishMetadataUseCase{
		mediaRepo: mediaRepo,
		publisher: publisher,
		config:    config,
	}
}

func (u *PublishMetadataUseCase) Execute(ctx context.Context, mediaID uuid.UUID) error {
	media, err := (*u.mediaRepo).GetByID(ctx, mediaID)
	if err != nil {
		return errors.New("failed to get media")
	}

	event := rabbitmq.AnalyzeMessage{
		UserID:       media.UserID,
		MediaID:      mediaID,
		MediaKey:     media.Key,
		ThumbnailKey: media.Thumbnail,
	}

	err = u.publisher.Publish(ctx, u.config.AnalyzeRequestQueueName, event)
	if err != nil {
		return errors.New("failed to publish metadata event")
	}

	return nil
}
