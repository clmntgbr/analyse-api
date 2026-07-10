package repository

import (
	"context"
	"go-api/domain/entity"
	"go-api/infrastructure/paginate"

	"github.com/google/uuid"
)

type MediaRepository interface {
	Create(ctx context.Context, media *entity.Media) error
	Update(ctx context.Context, media *entity.Media) error
	Delete(ctx context.Context, id uuid.UUID) error
	GetByUserID(ctx context.Context, userID uuid.UUID, query paginate.PaginateQuery) ([]*entity.Media, int64, error)
	GetByID(ctx context.Context, id uuid.UUID) (*entity.Media, error)
}
