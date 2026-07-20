package analysis

import (
	"context"
	"errors"
	"go-api/domain/entity"
	"go-api/domain/repository"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type GetAnalysisUseCase struct {
	analysisRepo *repository.AnalysisRepository
}

func NewGetAnalysisUseCase(analysisRepo *repository.AnalysisRepository) *GetAnalysisUseCase {
	return &GetAnalysisUseCase{analysisRepo: analysisRepo}
}

func (u *GetAnalysisUseCase) Execute(ctx context.Context, userID uuid.UUID, analysisID uuid.UUID) (*entity.Analysis, error) {
	analysis, err := (*u.analysisRepo).GetByID(ctx, analysisID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("analysis not found")
		}

		return nil, errors.New("failed to get analysis")
	}

	if analysis.UserID != userID {
		return nil, errors.New("analysis not found")
	}

	return analysis, nil
}
