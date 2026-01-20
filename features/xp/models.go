package xp

import (
	"time"

	"github.com/google/uuid"
)

// UserXP represents user XP and leveling data
type UserXP struct {
	ID             uuid.UUID `json:"id" db:"id"`
	UserID         uuid.UUID `json:"user_id" db:"user_id"`
	TotalXP        int       `json:"total_xp" db:"total_xp"`
	Level          int       `json:"level" db:"level"`
	CurrentLevelXP int       `json:"current_level_xp" db:"current_level_xp"`
	NextLevelXP   int       `json:"next_level_xp" db:"next_level_xp"`
	UpdatedAt     time.Time `json:"updated_at" db:"updated_at"`
}

// AddXPRequest represents a request to add XP
type AddXPRequest struct {
	Amount int    `json:"amount" validate:"required,gt=0"`
	Reason string `json:"reason,omitempty"`
}

// LeaderboardEntry represents a leaderboard entry
type LeaderboardEntry struct {
	UserID  uuid.UUID `json:"user_id"`
	Name    string    `json:"name"`
	TotalXP int       `json:"total_xp"`
	Level   int       `json:"level"`
	Rank    int       `json:"rank"`
}
