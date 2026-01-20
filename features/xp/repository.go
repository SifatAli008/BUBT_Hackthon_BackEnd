package xp

import (
	"database/sql"
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

func (r *Repository) GetByUserID(userID uuid.UUID) (*UserXP, error) {
	if r.db == nil {
		return nil, errors.ErrDatabase
	}
	xp := &UserXP{}
	query := `SELECT id, user_id, total_xp, level, current_level_xp, next_level_xp, updated_at FROM user_xp WHERE user_id = $1`
	err := r.db.QueryRow(query, userID).Scan(&xp.ID, &xp.UserID, &xp.TotalXP, &xp.Level, &xp.CurrentLevelXP, &xp.NextLevelXP, &xp.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.ErrNotFound
		}
		return nil, errors.WrapError(err, errors.ErrDatabase)
	}
	return xp, nil
}

func (r *Repository) CreateOrUpdate(xp *UserXP) error {
	if r.db == nil {
		return errors.ErrDatabase
	}
	query := `
		INSERT INTO user_xp (id, user_id, total_xp, level, current_level_xp, next_level_xp, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		ON CONFLICT (user_id) DO UPDATE SET
			total_xp = EXCLUDED.total_xp,
			level = EXCLUDED.level,
			current_level_xp = EXCLUDED.current_level_xp,
			next_level_xp = EXCLUDED.next_level_xp,
			updated_at = EXCLUDED.updated_at
		RETURNING id, user_id, total_xp, level, current_level_xp, next_level_xp, updated_at
	`
	err := r.db.QueryRow(query, xp.ID, xp.UserID, xp.TotalXP, xp.Level, xp.CurrentLevelXP, xp.NextLevelXP, time.Now()).Scan(&xp.ID, &xp.UserID, &xp.TotalXP, &xp.Level, &xp.CurrentLevelXP, &xp.NextLevelXP, &xp.UpdatedAt)
	if err != nil {
		return errors.WrapError(err, errors.ErrDatabase)
	}
	return nil
}

func (r *Repository) GetLeaderboard(limit int) ([]*LeaderboardEntry, error) {
	if r.db == nil {
		return nil, errors.ErrDatabase
	}
	if limit <= 0 {
		limit = 10
	}
	query := `
		SELECT ux.user_id, u.name, ux.total_xp, ux.level,
			ROW_NUMBER() OVER (ORDER BY ux.total_xp DESC) as rank
		FROM user_xp ux
		JOIN users u ON ux.user_id = u.id
		ORDER BY ux.total_xp DESC
		LIMIT $1
	`
	rows, err := r.db.Query(query, limit)
	if err != nil {
		return nil, errors.WrapError(err, errors.ErrDatabase)
	}
	defer rows.Close()
	var entries []*LeaderboardEntry
	for rows.Next() {
		e := &LeaderboardEntry{}
		if err := rows.Scan(&e.UserID, &e.Name, &e.TotalXP, &e.Level, &e.Rank); err != nil {
			return nil, errors.WrapError(err, errors.ErrDatabase)
		}
		entries = append(entries, e)
	}
	return entries, nil
}
