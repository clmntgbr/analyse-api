package presenter

import (
	"go-api/domain/entity"
	"time"
)

type PlanResponse struct {
	ID              string         `json:"id"`
	Name            string         `json:"name"`
	Price           float64        `json:"price"`
	Currency        string         `json:"currency"`
	BillingInterval string         `json:"billingInterval"`
	Description     string         `json:"description"`
	Slug            string         `json:"slug"`
	IsActive        bool           `json:"isActive"`
	Quota           *QuotaResponse `json:"quota"`
	CreatedAt       time.Time      `json:"createdAt"`
	UpdatedAt       time.Time      `json:"updatedAt"`
}

func NewPlanResponse(plan *entity.Plan) *PlanResponse {
	return &PlanResponse{
		ID:              plan.ID.String(),
		Name:            plan.Name,
		Price:           plan.Price,
		Currency:        string(plan.Currency),
		BillingInterval: string(plan.BillingInterval),
		Description:     plan.Description,
		Slug:            plan.Slug,
		IsActive:        plan.IsActive,
		Quota:           NewQuotaResponse(&plan.Quota),
		CreatedAt:       plan.CreatedAt,
		UpdatedAt:       plan.UpdatedAt,
	}
}

func NewPlanResponses(plans []*entity.Plan) []*PlanResponse {
	responses := make([]*PlanResponse, 0, len(plans))
	for _, plan := range plans {
		responses = append(responses, NewPlanResponse(plan))
	}
	return responses
}
