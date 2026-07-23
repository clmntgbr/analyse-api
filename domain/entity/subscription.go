package entity

import (
	"time"

	"github.com/google/uuid"
)

type SubscriptionStatus string

const (
	SubscriptionStatusActive    SubscriptionStatus = "active"
	SubscriptionStatusInactive  SubscriptionStatus = "inactive"
	SubscriptionStatusCancelled SubscriptionStatus = "cancelled"
	SubscriptionStatusPending   SubscriptionStatus = "pending"
	SubscriptionStatusPastDue   SubscriptionStatus = "past_due"
)

type Subscription struct {
	ID     uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	UserID uuid.UUID `gorm:"type:uuid;not null" json:"user_id"`
	User   User      `gorm:"foreignKey:UserID" json:"user"`

	PlanID uuid.UUID `gorm:"type:uuid;not null" json:"plan_id"`
	Plan   Plan      `gorm:"foreignKey:PlanID" json:"plan"`

	StripeCustomerID     string `json:"stripe_customer_id"`
	StripeSubscriptionID string `json:"stripe_subscription_id"`

	SubscriptionStatus    SubscriptionStatus `json:"subscription_status" gorm:"type:varchar(255);not null"`
	SubscriptionStartDate time.Time          `json:"subscription_start_date" gorm:"type:timestamp;not null"`
	SubscriptionEndDate   time.Time          `json:"subscription_end_date" gorm:"type:timestamp;not null"`

	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

func (Subscription) TableName() string {
	return "subscriptions"
}
