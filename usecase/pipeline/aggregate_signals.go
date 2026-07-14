package pipeline

import (
	"context"
	"errors"

	"go-api/domain/aggregate"
	"go-api/domain/entity"
	"go-api/domain/enum"
	"go-api/domain/repository"

	"github.com/google/uuid"
)

var requiredSignalNames = []string{"metadata", "heuristics", "ai_model"}

type AggregateAnalysisUseCase struct {
	mediaRepo  *repository.MediaRepository
	signalRepo *repository.SignalRepository
}

func NewAggregateAnalysisUseCase(
	mediaRepo *repository.MediaRepository,
	signalRepo *repository.SignalRepository,
) *AggregateAnalysisUseCase {
	return &AggregateAnalysisUseCase{
		mediaRepo:  mediaRepo,
		signalRepo: signalRepo,
	}
}

func (u *AggregateAnalysisUseCase) Execute(ctx context.Context, mediaID uuid.UUID) error {
	media, err := (*u.mediaRepo).GetByID(ctx, mediaID)
	if err != nil {
		return errors.New("media not found")
	}

	signals, err := (*u.signalRepo).GetByMediaID(ctx, mediaID)
	if err != nil {
		return errors.New("failed to load signals")
	}

	if !hasAllRequiredSignals(signals) {
		return errors.New("not all signals are ready")
	}

	result := aggregate.Compute(toEntitySignals(signals))

	media.FinalScore = result.FinalScore
	media.AnalysisConfidence = result.Confidence
	media.Verdict = result.Verdict
	media.Statuses = append(media.Statuses, enum.MediaStatusAnalyzed)
	media.Status = enum.MediaStatusAnalyzed

	return (*u.mediaRepo).Update(ctx, media)
}

func hasAllRequiredSignals(signals []*entity.Signal) bool {
	found := make(map[string]struct{}, len(requiredSignalNames))
	for _, signal := range signals {
		found[signal.Name] = struct{}{}
	}

	for _, name := range requiredSignalNames {
		if _, ok := found[name]; !ok {
			return false
		}
	}

	return true
}

func toEntitySignals(signals []*entity.Signal) []entity.Signal {
	result := make([]entity.Signal, 0, len(signals))
	for _, signal := range signals {
		if signal == nil {
			continue
		}
		result = append(result, *signal)
	}

	return result
}
