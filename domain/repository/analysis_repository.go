package repository

import (
	"context"
	"go-api/domain/entity"
	"go-api/infrastructure/paginate"

	"github.com/google/uuid"
)

type AnalysisRepository interface {
	Create(ctx context.Context, analysis *entity.Analysis) error
	Update(ctx context.Context, analysis *entity.Analysis) error
	Delete(ctx context.Context, id uuid.UUID) error
	GetByID(ctx context.Context, id uuid.UUID) (*entity.Analysis, error)
	GetByUserID(ctx context.Context, userID uuid.UUID, query paginate.PaginateQuery) ([]*entity.Analysis, int64, error)
	GetStatisticsByUserID(ctx context.Context, userID uuid.UUID) (*entity.MediaStatistics, error)
}
