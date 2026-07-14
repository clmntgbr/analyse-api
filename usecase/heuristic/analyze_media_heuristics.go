package heuristic

import (
	"context"
	"errors"
	"fmt"

	mediadto "go-api/infrastructure/media"
	heuristicsinfra "go-api/infrastructure/heuristics"
	"go-api/infrastructure/storage"

	"github.com/google/uuid"
)

type AnalyzeMediaHeuristicsUseCase struct {
	storage  *storage.MinIOStorage
	analyzer *heuristicsinfra.Analyzer
}

func NewAnalyzeMediaHeuristicsUseCase(
	storage *storage.MinIOStorage,
	analyzer *heuristicsinfra.Analyzer,
) *AnalyzeMediaHeuristicsUseCase {
	return &AnalyzeMediaHeuristicsUseCase{
		storage:  storage,
		analyzer: analyzer,
	}
}

func (uc *AnalyzeMediaHeuristicsUseCase) Execute(
	ctx context.Context,
	userID uuid.UUID,
	mediaKey string,
) (*heuristicsinfra.AnalysisResult, error) {
	objectKey := mediadto.NewObjectKey(userID, mediaKey)

	reader, err := uc.storage.Get(ctx, objectKey)
	if err != nil {
		return nil, fmt.Errorf("failed to download media %q: %w", objectKey, err)
	}
	defer reader.Close()

	result, err := uc.analyzer.Analyze(reader)
	if err != nil {
		return nil, errors.New("failed to analyze media heuristics")
	}

	return &result, nil
}
