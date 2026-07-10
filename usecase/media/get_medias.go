package media

import (
	"context"
	"go-api/domain/entity"
	"go-api/domain/repository"
	"go-api/infrastructure/paginate"

	"github.com/google/uuid"
)

type GetMediasUseCase struct {
	mediaRepo *repository.MediaRepository
}

func NewGetMediasUseCase(mediaRepo *repository.MediaRepository) *GetMediasUseCase {
	return &GetMediasUseCase{mediaRepo: mediaRepo}
}

func (u *GetMediasUseCase) Execute(ctx context.Context, userID uuid.UUID, query paginate.PaginateQuery) ([]*entity.Media, int64, error) {
	medias, total, err := (*u.mediaRepo).GetByUserID(ctx, userID, query)
	if err != nil {
		return []*entity.Media{}, 0, err
	}

	return medias, total, nil
}
