package gorm

import (
	"context"
	"errors"
	"go-api/domain/entity"
	"go-api/domain/repository"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type mediaRepository struct {
	db *gorm.DB
}

func NewMediaRepository(db *gorm.DB) repository.MediaRepository {
	return &mediaRepository{db: db}
}

func (r *mediaRepository) Create(ctx context.Context, media *entity.Media) error {
	return dbWithContext(ctx, r.db).Create(media).Error
}

func (r *mediaRepository) Update(ctx context.Context, media *entity.Media) error {
	return dbWithContext(ctx, r.db).Save(media).Error
}

func (r *mediaRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return dbWithContext(ctx, r.db).Delete(&entity.Media{}, id).Error
}

func (r *mediaRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.Media, error) {
	var media entity.Media
	err := dbWithContext(ctx, r.db).
		Where("id = ?", id).
		Preload("Signals").
		Preload("Insight").
		Preload("Analysis").
		First(&media).Error
	if err != nil {
		return nil, err
	}
	if media.ID == uuid.Nil {
		return nil, errors.New("media not found")
	}
	return &media, nil
}

func (r *mediaRepository) GetByKey(ctx context.Context, key string) (*entity.Media, error) {
	var media entity.Media
	err := dbWithContext(ctx, r.db).
		Where("key = ?", key).
		Preload("Analysis").
		First(&media).Error
	if err != nil {
		return nil, err
	}
	if media.ID == uuid.Nil {
		return nil, errors.New("media not found")
	}
	return &media, nil
}
