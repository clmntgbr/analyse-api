package media

import (
	"context"
	"errors"
	"go-api/domain/entity"
	"go-api/domain/repository"

	"github.com/google/uuid"
)

type GetMediaStatisticsUseCase struct {
	mediaRepo *repository.MediaRepository
}

func NewGetMediaStatisticsUseCase(mediaRepo *repository.MediaRepository) *GetMediaStatisticsUseCase {
	return &GetMediaStatisticsUseCase{mediaRepo: mediaRepo}
}

func (u *GetMediaStatisticsUseCase) Execute(ctx context.Context, userID uuid.UUID) (*entity.MediaStatistics, error) {
	stats, err := (*u.mediaRepo).GetStatisticsByUserID(ctx, userID)
	if err != nil {
		return nil, errors.New("failed to get media statistics")
	}

	return stats, nil
}
