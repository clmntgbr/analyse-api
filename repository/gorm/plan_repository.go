package gorm

import (
	"context"
	"go-api/domain/entity"
	"go-api/domain/repository"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type planRepository struct {
	db *gorm.DB
}

func NewPlanRepository(db *gorm.DB) repository.PlanRepository {
	return &planRepository{db: db}
}

func (r *planRepository) Create(ctx context.Context, plan *entity.Plan) error {
	return dbWithContext(ctx, r.db).Create(plan).Error
}

func (r *planRepository) Update(ctx context.Context, plan *entity.Plan) error {
	return dbWithContext(ctx, r.db).Save(plan).Error
}

func (r *planRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return dbWithContext(ctx, r.db).Delete(&entity.Plan{}, id).Error
}

func (r *planRepository) GetAll(ctx context.Context) ([]*entity.Plan, error) {
	var plans []*entity.Plan
	err := dbWithContext(ctx, r.db).
		Find(&plans).Error
	if err != nil {
		return nil, err
	}
	return plans, nil
}
