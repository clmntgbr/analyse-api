package gorm

import (
	"context"
	"errors"
	"go-api/domain/entity"
	"go-api/domain/repository"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type insightRepository struct {
	db *gorm.DB
}

func NewInsightRepository(db *gorm.DB) repository.InsightRepository {
	return &insightRepository{db: db}
}

func (r *insightRepository) Create(ctx context.Context, insight *entity.Insight) error {
	return dbWithContext(ctx, r.db).Create(insight).Error
}

func (r *insightRepository) Update(ctx context.Context, insight *entity.Insight) error {
	return dbWithContext(ctx, r.db).Save(insight).Error
}

func (r *insightRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return dbWithContext(ctx, r.db).Delete(&entity.Insight{}, id).Error
}

func (r *insightRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.Insight, error) {
	var insight entity.Insight
	err := dbWithContext(ctx, r.db).Where("id = ?", id).First(&insight).Error
	if err != nil {
		return nil, err
	}
	if insight.ID == uuid.Nil {
		return nil, errors.New("insight not found")
	}

	return &insight, nil
}
