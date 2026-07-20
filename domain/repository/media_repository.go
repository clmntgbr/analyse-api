package repository

import (
	"context"
	"go-api/domain/entity"

	"github.com/google/uuid"
)

type MediaRepository interface {
	Create(ctx context.Context, media *entity.Media) error
	Update(ctx context.Context, media *entity.Media) error
	Delete(ctx context.Context, id uuid.UUID) error
	GetByID(ctx context.Context, id uuid.UUID) (*entity.Media, error)
	GetByKey(ctx context.Context, key string) (*entity.Media, error)
}
