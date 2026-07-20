package gorm

import (
	"context"
	"errors"
	"go-api/domain/aggregate"
	"go-api/domain/entity"
	"go-api/domain/repository"
	"go-api/infrastructure/paginate"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type analysisRepository struct {
	db *gorm.DB
}

func NewAnalysisRepository(db *gorm.DB) repository.AnalysisRepository {
	return &analysisRepository{db: db}
}

func (r *analysisRepository) Create(ctx context.Context, analysis *entity.Analysis) error {
	return dbWithContext(ctx, r.db).Create(analysis).Error
}

func (r *analysisRepository) Update(ctx context.Context, analysis *entity.Analysis) error {
	return dbWithContext(ctx, r.db).Save(analysis).Error
}

func (r *analysisRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return dbWithContext(ctx, r.db).Delete(&entity.Analysis{}, id).Error
}

func (r *analysisRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.Analysis, error) {
	var analysis entity.Analysis
	err := dbWithContext(ctx, r.db).
		Where("id = ?", id).
		Preload("Medias", func(db *gorm.DB) *gorm.DB {
			return db.Order("medias.created_at ASC")
		}).
		Preload("Medias.Signals").
		Preload("Medias.Insight").
		First(&analysis).Error
	if err != nil {
		return nil, err
	}
	if analysis.ID == uuid.Nil {
		return nil, errors.New("analysis not found")
	}
	return &analysis, nil
}

func (r *analysisRepository) GetByUserID(ctx context.Context, userID uuid.UUID, query paginate.PaginateQuery) ([]*entity.Analysis, int64, error) {
	var analyses []*entity.Analysis

	db := dbWithContext(ctx, r.db).Model(&entity.Analysis{}).
		Where("analyses.user_id = ?", userID)

	if query.Search != "" {
		db = db.Joins("JOIN medias ON medias.analysis_id = analyses.id").
			Where("medias.filename ILIKE ? OR medias.key ILIKE ?", "%"+query.Search+"%", "%"+query.Search+"%").
			Distinct()
	}

	db, total, err := Paginate(db, query)
	if err != nil {
		return nil, 0, err
	}

	err = db.
		Preload("Medias", func(db *gorm.DB) *gorm.DB {
			return db.Order("medias.created_at ASC")
		}).
		Find(&analyses).Error
	if err != nil {
		return nil, 0, err
	}

	return analyses, total, nil
}

func (r *analysisRepository) GetStatisticsByUserID(ctx context.Context, userID uuid.UUID) (*entity.MediaStatistics, error) {
	var stats entity.MediaStatistics

	err := dbWithContext(ctx, r.db).Raw(`
		SELECT
			COUNT(*) FILTER (WHERE verdict <> '') AS analyses_count,
			COUNT(*) FILTER (WHERE verdict = ?) AS real_image_count,
			COUNT(*) FILTER (WHERE verdict = ?) AS ai_image_count,
			COALESCE(AVG(final_score) FILTER (WHERE verdict <> '' AND final_score >= 0), 0) AS average_score
		FROM analyses
		WHERE user_id = ?
	`,
		aggregate.VerdictLikelyReal,
		aggregate.VerdictLikelyAI,
		userID,
	).Scan(&stats).Error
	if err != nil {
		return nil, err
	}

	return &stats, nil
}
