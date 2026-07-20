package analysis

import (
	"context"
	"errors"
	"go-api/domain/entity"
	"go-api/domain/repository"

	"github.com/google/uuid"
)

type GetStatisticsUseCase struct {
	analysisRepo *repository.AnalysisRepository
}

func NewGetStatisticsUseCase(analysisRepo *repository.AnalysisRepository) *GetStatisticsUseCase {
	return &GetStatisticsUseCase{analysisRepo: analysisRepo}
}

func (u *GetStatisticsUseCase) Execute(ctx context.Context, userID uuid.UUID) (*entity.MediaStatistics, error) {
	stats, err := (*u.analysisRepo).GetStatisticsByUserID(ctx, userID)
	if err != nil {
		return nil, errors.New("failed to get analysis statistics")
	}

	return stats, nil
}
