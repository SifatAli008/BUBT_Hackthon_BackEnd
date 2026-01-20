package profiles

import (
	"time"

	"github.com/google/uuid"
)

type CommunityProfile struct {
	ID                  uuid.UUID  `json:"id" db:"id"`
	UserID              uuid.UUID  `json:"user_id" db:"user_id"`
	Username            string     `json:"username" db:"username"`
	AvatarURL           string     `json:"avatar_url,omitempty" db:"avatar_url"`
	CommunityRole       string     `json:"community_role" db:"community_role"`
	Bio                 string     `json:"bio,omitempty" db:"bio"`
	PreferredItems      []string   `json:"preferred_items,omitempty" db:"preferred_items"`
	AvoidItems          []string   `json:"avoid_items,omitempty" db:"avoid_items"`
	DietaryRestrictions []string   `json:"dietary_restrictions,omitempty" db:"dietary_restrictions"`
	Allergens           []string   `json:"allergens,omitempty" db:"allergens"`
	AcceptsHotMeals     bool       `json:"accepts_hot_meals" db:"accepts_hot_meals"`
	DistancePreference  string     `json:"distance_preference,omitempty" db:"distance_preference"`
	Visibility          string     `json:"visibility" db:"visibility"`
	NotificationsEnabled bool      `json:"notifications_enabled" db:"notifications_enabled"`
	NotifyOnClaim       bool       `json:"notify_on_claim" db:"notify_on_claim"`
	NotifyOnMessages    bool       `json:"notify_on_messages" db:"notify_on_messages"`
	CreatedAt           time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt           time.Time `json:"updated_at" db:"updated_at"`
}

type CreateProfileRequest struct {
	Username            string   `json:"username" validate:"required,min=3,max=255"`
	AvatarURL           string   `json:"avatar_url,omitempty"`
	CommunityRole       string   `json:"community_role,omitempty"`
	Bio                 string   `json:"bio,omitempty"`
	PreferredItems      []string `json:"preferred_items,omitempty"`
	AvoidItems          []string `json:"avoid_items,omitempty"`
	DietaryRestrictions []string `json:"dietary_restrictions,omitempty"`
	Allergens           []string `json:"allergens,omitempty"`
	AcceptsHotMeals     bool     `json:"accepts_hot_meals,omitempty"`
	DistancePreference  string   `json:"distance_preference,omitempty"`
	Visibility          string   `json:"visibility,omitempty"`
	NotificationsEnabled *bool   `json:"notifications_enabled,omitempty"`
	NotifyOnClaim       *bool    `json:"notify_on_claim,omitempty"`
	NotifyOnMessages    *bool    `json:"notify_on_messages,omitempty"`
}

type UpdateProfileRequest struct {
	AvatarURL           string   `json:"avatar_url,omitempty"`
	CommunityRole       string   `json:"community_role,omitempty"`
	Bio                 string   `json:"bio,omitempty"`
	PreferredItems      []string `json:"preferred_items,omitempty"`
	AvoidItems          []string `json:"avoid_items,omitempty"`
	DietaryRestrictions []string `json:"dietary_restrictions,omitempty"`
	Allergens           []string `json:"allergens,omitempty"`
	AcceptsHotMeals     *bool    `json:"accepts_hot_meals,omitempty"`
	DistancePreference  string   `json:"distance_preference,omitempty"`
	Visibility          string   `json:"visibility,omitempty"`
	NotificationsEnabled *bool   `json:"notifications_enabled,omitempty"`
	NotifyOnClaim       *bool    `json:"notify_on_claim,omitempty"`
	NotifyOnMessages    *bool    `json:"notify_on_messages,omitempty"`
}
