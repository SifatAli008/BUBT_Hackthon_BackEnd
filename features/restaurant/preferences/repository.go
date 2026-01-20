package preferences

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

func (r *Repository) GetByUserID(userID uuid.UUID) (*RestaurantPreferences, error) {
	if r.db == nil {
		return nil, errors.ErrDatabase
	}
	prefs := &RestaurantPreferences{}
	query := `SELECT id, user_id, cuisine_type, operating_hours, donation_preferences, created_at, updated_at FROM restaurant_preferences WHERE user_id = $1`
	err := r.db.QueryRow(query, userID).Scan(&prefs.ID, &prefs.UserID, &prefs.CuisineType, &prefs.OperatingHours, pq.Array(&prefs.DonationPreferences), &prefs.CreatedAt, &prefs.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.ErrNotFound
		}
		return nil, errors.WrapError(err, errors.ErrDatabase)
	}
	return prefs, nil
}

func (r *Repository) CreateOrUpdate(prefs *RestaurantPreferences) error {
	if r.db == nil {
		return errors.ErrDatabase
	}
	query := `INSERT INTO restaurant_preferences (id, user_id, cuisine_type, operating_hours, donation_preferences, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7) ON CONFLICT (user_id) DO UPDATE SET cuisine_type=EXCLUDED.cuisine_type, operating_hours=EXCLUDED.operating_hours, donation_preferences=EXCLUDED.donation_preferences, updated_at=EXCLUDED.updated_at RETURNING id, user_id, cuisine_type, operating_hours, donation_preferences, created_at, updated_at`
	now := time.Now()
	return r.db.QueryRow(query, prefs.ID, prefs.UserID, prefs.CuisineType, prefs.OperatingHours, pq.Array(prefs.DonationPreferences), now, now).Scan(&prefs.ID, &prefs.UserID, &prefs.CuisineType, &prefs.OperatingHours, pq.Array(&prefs.DonationPreferences), &prefs.CreatedAt, &prefs.UpdatedAt)
}
