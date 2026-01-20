package badges

import (
	"time"

	"github.com/google/uuid"
)

// Badge represents a user badge
type Badge struct {
	ID         uuid.UUID `json:"id" db:"id"`
	UserID     uuid.UUID `json:"user_id" db:"user_id"`
	BadgeID    string    `json:"badge_id" db:"badge_id"`
	Name       string    `json:"name" db:"name"`
	Description string   `json:"description,omitempty" db:"description"`
	Icon       string    `json:"icon,omitempty" db:"icon"`
	UnlockedAt time.Time `json:"unlocked_at" db:"unlocked_at"`
	XPReward   int       `json:"xp_reward" db:"xp_reward"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
}

// UnlockBadgeRequest represents a request to unlock a badge
type UnlockBadgeRequest struct {
	BadgeID    string `json:"badge_id" validate:"required"`
	Name       string `json:"name" validate:"required"`
	Description string `json:"description,omitempty"`
	Icon       string `json:"icon,omitempty"`
	XPReward   int    `json:"xp_reward,omitempty"`
}

// AvailableBadge represents an available badge definition
type AvailableBadge struct {
	BadgeID    string `json:"badge_id"`
	Name       string `json:"name"`
	Description string `json:"description"`
	Icon       string `json:"icon"`
	XPReward   int    `json:"xp_reward"`
}
