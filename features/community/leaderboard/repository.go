package leaderboard

import (
	"database/sql"
	"encoding/json"
	"foodlink_backend/database"
	"foodlink_backend/errors"

	"github.com/google/uuid"
)

type Repository struct {
	db *sql.DB
}

func NewRepository() *Repository {
	return &Repository{db: database.GetDB()}
}

func (r *Repository) GetByType(leaderboardType string) (*Leaderboard, error) {
	if r.db == nil {
		return nil, errors.ErrDatabase
	}
	lb := &Leaderboard{}
	var entriesJSON []byte
	query := `SELECT id, type, entries, updated_at FROM community_leaderboard WHERE type = $1 ORDER BY updated_at DESC LIMIT 1`
	err := r.db.QueryRow(query, leaderboardType).Scan(&lb.ID, &lb.Type, &entriesJSON, &lb.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.ErrNotFound
		}
		return nil, errors.WrapError(err, errors.ErrDatabase)
	}
	if len(entriesJSON) > 0 {
		json.Unmarshal(entriesJSON, &lb.Entries)
	}
	return lb, nil
}

func (r *Repository) GetImpact() (*Impact, error) {
	if r.db == nil {
		return nil, errors.ErrDatabase
	}
	impact := &Impact{}
	var weeklyTrendJSON, personalContributionJSON []byte
	query := `SELECT id, total_surplus_kg, donations, co2_prevented_kg, water_saved_liters, meals_provided, weekly_trend, personal_contribution, updated_at FROM community_impact ORDER BY updated_at DESC LIMIT 1`
	err := r.db.QueryRow(query).Scan(&impact.ID, &impact.TotalSurplusKg, &impact.Donations, &impact.CO2PreventedKg, &impact.WaterSavedLiters, &impact.MealsProvided, &weeklyTrendJSON, &personalContributionJSON, &impact.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.ErrNotFound
		}
		return nil, errors.WrapError(err, errors.ErrDatabase)
	}
	if len(weeklyTrendJSON) > 0 {
		json.Unmarshal(weeklyTrendJSON, &impact.WeeklyTrend)
	}
	if len(personalContributionJSON) > 0 {
		json.Unmarshal(personalContributionJSON, &impact.PersonalContribution)
	}
	return impact, nil
}

func (r *Repository) GetPersonalImpact(userID uuid.UUID) (*Impact, error) {
	impact := &Impact{}
	impact.PersonalContribution = JSONB{
		"surplus_shared": 0,
		"leftovers_shared": 0,
		"volunteer_hours": 0,
		"meals_provided": 0,
	}
	return impact, nil
}
