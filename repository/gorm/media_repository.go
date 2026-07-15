package gorm

import (
	"context"
	"errors"
	"go-api/domain/entity"
	"go-api/domain/repository"
	"go-api/infrastructure/paginate"

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

func (r *mediaRepository) GetByUserID(ctx context.Context, userID uuid.UUID, query paginate.PaginateQuery) ([]*entity.Media, int64, error) {
	var medias []*entity.Media

	db := dbWithContext(ctx, r.db).Model(&entity.Media{}).
		Where("medias.user_id = ?", userID)
	if query.Search != "" {
		db = db.Where("medias.name ILIKE ? OR medias.key ILIKE ?", "%"+query.Search+"%", "%"+query.Search+"%")
	}

	db, total, err := Paginate(db, query)
	if err != nil {
		return nil, 0, err
	}

	err = db.Find(&medias).Error
	if err != nil {
		return nil, 0, err
	}

	return medias, total, nil
}

func (r *mediaRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.Media, error) {
	var media entity.Media
	err := dbWithContext(ctx, r.db).Where("id = ?", id).Preload("Signals").Preload("Insight").First(&media).Error
	if err != nil {
		return nil, err
	}
	if media.ID == uuid.Nil {
		return nil, errors.New("media not found")
	}
	return &media, nil
}
