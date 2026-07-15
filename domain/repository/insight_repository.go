package repository

import (
	"context"
	"go-api/domain/entity"

	"github.com/google/uuid"
)

type InsightRepository interface {
	Create(ctx context.Context, insight *entity.Insight) error
	Update(ctx context.Context, insight *entity.Insight) error
	Delete(ctx context.Context, id uuid.UUID) error
	GetByID(ctx context.Context, id uuid.UUID) (*entity.Insight, error)
}
