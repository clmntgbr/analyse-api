package entity

import (
	"time"

	"github.com/google/uuid"
)

type BillingInterval string

const (
	BillingIntervalMonth BillingInterval = "month"
	BillingIntervalYear  BillingInterval = "year"
)

type Currency string

const (
	CurrencyUSD Currency = "USD"
	CurrencyEUR Currency = "EUR"
)

type Plan struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Slug        string    `json:"slug" gorm:"type:varchar(255);not null"`

	StripePriceID string `json:"stripe_price_id" gorm:"type:varchar(255);not null"`

	IsActive bool `json:"is_active" gorm:"type:boolean;not null"`

	BillingInterval BillingInterval `json:"billing_interval" gorm:"type:varchar(255);not null;default:'month'"`
	Price           float64         `json:"price" gorm:"type:decimal(10,2);not null;default:0"`
	Currency        Currency        `json:"currency" gorm:"type:varchar(255);not null;default:'EUR'"`

	QuotaID uuid.UUID `gorm:"type:uuid;not null;index:idx_plan_quota_id" json:"quota_id"`
	Quota   Quota     `gorm:"foreignKey:QuotaID" json:"quota"`

	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

func (Plan) TableName() string {
	return "plans"
}
