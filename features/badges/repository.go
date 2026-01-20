package badges

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

func (r *Repository) GetByUserID(userID uuid.UUID) ([]*Badge, error) {
	if r.db == nil {
		return nil, errors.ErrDatabase
	}
	query := `SELECT id, user_id, badge_id, name, description, icon, unlocked_at, xp_reward, created_at FROM badges WHERE user_id = $1 ORDER BY unlocked_at DESC`
	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, errors.WrapError(err, errors.ErrDatabase)
	}
	defer rows.Close()
	var badges []*Badge
	for rows.Next() {
		b := &Badge{}
		if err := rows.Scan(&b.ID, &b.UserID, &b.BadgeID, &b.Name, &b.Description, &b.Icon, &b.UnlockedAt, &b.XPReward, &b.CreatedAt); err != nil {
			return nil, errors.WrapError(err, errors.ErrDatabase)
		}
		badges = append(badges, b)
	}
	return badges, nil
}

func (r *Repository) GetByUserIDAndBadgeID(userID uuid.UUID, badgeID string) (*Badge, error) {
	if r.db == nil {
		return nil, errors.ErrDatabase
	}
	b := &Badge{}
	query := `SELECT id, user_id, badge_id, name, description, icon, unlocked_at, xp_reward, created_at FROM badges WHERE user_id = $1 AND badge_id = $2`
	err := r.db.QueryRow(query, userID, badgeID).Scan(&b.ID, &b.UserID, &b.BadgeID, &b.Name, &b.Description, &b.Icon, &b.UnlockedAt, &b.XPReward, &b.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.ErrNotFound
		}
		return nil, errors.WrapError(err, errors.ErrDatabase)
	}
	return b, nil
}

func (r *Repository) Create(badge *Badge) error {
	if r.db == nil {
		return errors.ErrDatabase
	}
	query := `INSERT INTO badges (id, user_id, badge_id, name, description, icon, unlocked_at, xp_reward, created_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) ON CONFLICT (user_id, badge_id) DO NOTHING RETURNING id, user_id, badge_id, name, description, icon, unlocked_at, xp_reward, created_at`
	now := time.Now()
	err := r.db.QueryRow(query, badge.ID, badge.UserID, badge.BadgeID, badge.Name, badge.Description, badge.Icon, badge.UnlockedAt, badge.XPReward, now).Scan(&badge.ID, &badge.UserID, &badge.BadgeID, &badge.Name, &badge.Description, &badge.Icon, &badge.UnlockedAt, &badge.XPReward, &badge.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return errors.ErrAlreadyExists
		}
		return errors.WrapError(err, errors.ErrDatabase)
	}
	return nil
}
