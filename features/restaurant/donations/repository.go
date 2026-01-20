package donations

import (
	"database/sql"
	"encoding/json"
	"foodlink_backend/database"
	"foodlink_backend/errors"
	"time"

	"github.com/google/uuid"
)

type Repository struct {
	db *sql.DB
}

func NewRepository() *Repository {
	return &Repository{db: database.GetDB()}
}

func (r *Repository) GetAllByUserID(userID uuid.UUID) ([]*DonationLog, error) {
	if r.db == nil {
		return nil, errors.ErrDatabase
	}
	query := `SELECT id, user_id, date, recipient_type, recipient_name, items, quantity, unit, meals_provided, co2_saved_kg, notes, created_at FROM restaurant_donation_logs WHERE user_id = $1 ORDER BY date DESC, created_at DESC`
	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, errors.WrapError(err, errors.ErrDatabase)
	}
	defer rows.Close()
	var logs []*DonationLog
	for rows.Next() {
		log := &DonationLog{}
		if err := rows.Scan(&log.ID, &log.UserID, &log.Date, &log.RecipientType, &log.RecipientName, &log.Items, &log.Quantity, &log.Unit, &log.MealsProvided, &log.CO2SavedKg, &log.Notes, &log.CreatedAt); err != nil {
			return nil, errors.WrapError(err, errors.ErrDatabase)
		}
		logs = append(logs, log)
	}
	return logs, nil
}

func (r *Repository) Create(log *DonationLog) error {
	if r.db == nil {
		return errors.ErrDatabase
	}
	query := `INSERT INTO restaurant_donation_logs (id, user_id, date, recipient_type, recipient_name, items, quantity, unit, meals_provided, co2_saved_kg, notes, created_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12) RETURNING id, user_id, date, recipient_type, recipient_name, items, quantity, unit, meals_provided, co2_saved_kg, notes, created_at`
	return r.db.QueryRow(query, log.ID, log.UserID, log.Date, log.RecipientType, log.RecipientName, log.Items, log.Quantity, log.Unit, log.MealsProvided, log.CO2SavedKg, log.Notes, time.Now()).Scan(&log.ID, &log.UserID, &log.Date, &log.RecipientType, &log.RecipientName, &log.Items, &log.Quantity, &log.Unit, &log.MealsProvided, &log.CO2SavedKg, &log.Notes, &log.CreatedAt)
}

func (r *Repository) GetImpactByUserID(userID uuid.UUID) (*ImpactMetrics, error) {
	if r.db == nil {
		return nil, errors.ErrDatabase
	}
	metrics := &ImpactMetrics{}
	var weeklyTrendJSON, monthlyTrendJSON, categoryBreakdownJSON []byte
	query := `SELECT id, user_id, waste_prevented_kg, surplus_donation_rate, water_saved_liters, co2_prevented_kg, sustainability_score, weekly_trend, monthly_trend, category_breakdown, updated_at FROM restaurant_impact_metrics WHERE user_id = $1`
	err := r.db.QueryRow(query, userID).Scan(&metrics.ID, &metrics.UserID, &metrics.WastePreventedKg, &metrics.SurplusDonationRate, &metrics.WaterSavedLiters, &metrics.CO2PreventedKg, &metrics.SustainabilityScore, &weeklyTrendJSON, &monthlyTrendJSON, &categoryBreakdownJSON, &metrics.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.ErrNotFound
		}
		return nil, errors.WrapError(err, errors.ErrDatabase)
	}
	if len(weeklyTrendJSON) > 0 {
		json.Unmarshal(weeklyTrendJSON, &metrics.WeeklyTrend)
	}
	if len(monthlyTrendJSON) > 0 {
		json.Unmarshal(monthlyTrendJSON, &metrics.MonthlyTrend)
	}
	if len(categoryBreakdownJSON) > 0 {
		json.Unmarshal(categoryBreakdownJSON, &metrics.CategoryBreakdown)
	}
	return metrics, nil
}
