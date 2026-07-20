package analysis

import (
	"context"
	"go-api/domain/entity"
	"go-api/domain/repository"
	"go-api/infrastructure/paginate"

	"github.com/google/uuid"
)

type GetAnalysesUseCase struct {
	analysisRepo *repository.AnalysisRepository
}

func NewGetAnalysesUseCase(analysisRepo *repository.AnalysisRepository) *GetAnalysesUseCase {
	return &GetAnalysesUseCase{analysisRepo: analysisRepo}
}

func (u *GetAnalysesUseCase) Execute(ctx context.Context, userID uuid.UUID, query paginate.PaginateQuery) ([]*entity.Analysis, int64, error) {
	analyses, total, err := (*u.analysisRepo).GetByUserID(ctx, userID, query)
	if err != nil {
		return []*entity.Analysis{}, 0, err
	}

	return analyses, total, nil
}
