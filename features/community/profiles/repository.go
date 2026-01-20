package profiles

import (
	"database/sql"
	"foodlink_backend/database"
	"foodlink_backend/errors"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

type Repository struct {
	db *sql.DB
}

func NewRepository() *Repository {
	return &Repository{db: database.GetDB()}
}

func (r *Repository) GetByUserID(userID uuid.UUID) (*CommunityProfile, error) {
	if r.db == nil {
		return nil, errors.ErrDatabase
	}
	profile := &CommunityProfile{}
	query := `SELECT id, user_id, username, avatar_url, community_role, bio, preferred_items, avoid_items, dietary_restrictions, allergens, accepts_hot_meals, distance_preference, visibility, notifications_enabled, notify_on_claim, notify_on_messages, created_at, updated_at FROM community_profiles WHERE user_id = $1`
	err := r.db.QueryRow(query, userID).Scan(&profile.ID, &profile.UserID, &profile.Username, &profile.AvatarURL, &profile.CommunityRole, &profile.Bio, pq.Array(&profile.PreferredItems), pq.Array(&profile.AvoidItems), pq.Array(&profile.DietaryRestrictions), pq.Array(&profile.Allergens), &profile.AcceptsHotMeals, &profile.DistancePreference, &profile.Visibility, &profile.NotificationsEnabled, &profile.NotifyOnClaim, &profile.NotifyOnMessages, &profile.CreatedAt, &profile.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.ErrNotFound
		}
		return nil, errors.WrapError(err, errors.ErrDatabase)
	}
	return profile, nil
}

func (r *Repository) GetByUsername(username string) (*CommunityProfile, error) {
	if r.db == nil {
		return nil, errors.ErrDatabase
	}
	profile := &CommunityProfile{}
	query := `SELECT id, user_id, username, avatar_url, community_role, bio, preferred_items, avoid_items, dietary_restrictions, allergens, accepts_hot_meals, distance_preference, visibility, notifications_enabled, notify_on_claim, notify_on_messages, created_at, updated_at FROM community_profiles WHERE username = $1`
	err := r.db.QueryRow(query, username).Scan(&profile.ID, &profile.UserID, &profile.Username, &profile.AvatarURL, &profile.CommunityRole, &profile.Bio, pq.Array(&profile.PreferredItems), pq.Array(&profile.AvoidItems), pq.Array(&profile.DietaryRestrictions), pq.Array(&profile.Allergens), &profile.AcceptsHotMeals, &profile.DistancePreference, &profile.Visibility, &profile.NotificationsEnabled, &profile.NotifyOnClaim, &profile.NotifyOnMessages, &profile.CreatedAt, &profile.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.ErrNotFound
		}
		return nil, errors.WrapError(err, errors.ErrDatabase)
	}
	return profile, nil
}

func (r *Repository) Create(profile *CommunityProfile) error {
	if r.db == nil {
		return errors.ErrDatabase
	}
	query := `INSERT INTO community_profiles (id, user_id, username, avatar_url, community_role, bio, preferred_items, avoid_items, dietary_restrictions, allergens, accepts_hot_meals, distance_preference, visibility, notifications_enabled, notify_on_claim, notify_on_messages, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18) RETURNING id, user_id, username, avatar_url, community_role, bio, preferred_items, avoid_items, dietary_restrictions, allergens, accepts_hot_meals, distance_preference, visibility, notifications_enabled, notify_on_claim, notify_on_messages, created_at, updated_at`
	now := time.Now()
	return r.db.QueryRow(query, profile.ID, profile.UserID, profile.Username, profile.AvatarURL, profile.CommunityRole, profile.Bio, pq.Array(profile.PreferredItems), pq.Array(profile.AvoidItems), pq.Array(profile.DietaryRestrictions), pq.Array(profile.Allergens), profile.AcceptsHotMeals, profile.DistancePreference, profile.Visibility, profile.NotificationsEnabled, profile.NotifyOnClaim, profile.NotifyOnMessages, now, now).Scan(&profile.ID, &profile.UserID, &profile.Username, &profile.AvatarURL, &profile.CommunityRole, &profile.Bio, pq.Array(&profile.PreferredItems), pq.Array(&profile.AvoidItems), pq.Array(&profile.DietaryRestrictions), pq.Array(&profile.Allergens), &profile.AcceptsHotMeals, &profile.DistancePreference, &profile.Visibility, &profile.NotificationsEnabled, &profile.NotifyOnClaim, &profile.NotifyOnMessages, &profile.CreatedAt, &profile.UpdatedAt)
}

func (r *Repository) Update(profile *CommunityProfile) error {
	if r.db == nil {
		return errors.ErrDatabase
	}
	query := `UPDATE community_profiles SET avatar_url=$1, community_role=$2, bio=$3, preferred_items=$4, avoid_items=$5, dietary_restrictions=$6, allergens=$7, accepts_hot_meals=$8, distance_preference=$9, visibility=$10, notifications_enabled=$11, notify_on_claim=$12, notify_on_messages=$13, updated_at=$14 WHERE user_id=$15 RETURNING id, user_id, username, avatar_url, community_role, bio, preferred_items, avoid_items, dietary_restrictions, allergens, accepts_hot_meals, distance_preference, visibility, notifications_enabled, notify_on_claim, notify_on_messages, created_at, updated_at`
	return r.db.QueryRow(query, profile.AvatarURL, profile.CommunityRole, profile.Bio, pq.Array(profile.PreferredItems), pq.Array(profile.AvoidItems), pq.Array(profile.DietaryRestrictions), pq.Array(profile.Allergens), profile.AcceptsHotMeals, profile.DistancePreference, profile.Visibility, profile.NotificationsEnabled, profile.NotifyOnClaim, profile.NotifyOnMessages, time.Now(), profile.UserID).Scan(&profile.ID, &profile.UserID, &profile.Username, &profile.AvatarURL, &profile.CommunityRole, &profile.Bio, pq.Array(&profile.PreferredItems), pq.Array(&profile.AvoidItems), pq.Array(&profile.DietaryRestrictions), pq.Array(&profile.Allergens), &profile.AcceptsHotMeals, &profile.DistancePreference, &profile.Visibility, &profile.NotificationsEnabled, &profile.NotifyOnClaim, &profile.NotifyOnMessages, &profile.CreatedAt, &profile.UpdatedAt)
}
