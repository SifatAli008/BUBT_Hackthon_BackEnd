package nutrition

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

func (r *Repository) GetByUserIDAndDateRange(userID uuid.UUID, startDate, endDate time.Time) ([]*NutritionData, error) {
	if r.db == nil {
		return nil, errors.ErrDatabase
	}
	query := `SELECT id, user_id, date, calories, protein, carbs, fats, fiber, sugar, sodium, vitamin_a, vitamin_b, vitamin_c, vitamin_d, iron, calcium, nutrition_score, created_at, updated_at FROM nutrition_data WHERE user_id = $1 AND date BETWEEN $2 AND $3 ORDER BY date DESC`
	rows, err := r.db.Query(query, userID, startDate, endDate)
	if err != nil {
		return nil, errors.WrapError(err, errors.ErrDatabase)
	}
	defer rows.Close()
	var data []*NutritionData
	for rows.Next() {
		d := &NutritionData{}
		if err := rows.Scan(&d.ID, &d.UserID, &d.Date, &d.Calories, &d.Protein, &d.Carbs, &d.Fats, &d.Fiber, &d.Sugar, &d.Sodium, &d.VitaminA, &d.VitaminB, &d.VitaminC, &d.VitaminD, &d.Iron, &d.Calcium, &d.NutritionScore, &d.CreatedAt, &d.UpdatedAt); err != nil {
			return nil, errors.WrapError(err, errors.ErrDatabase)
		}
		data = append(data, d)
	}
	return data, nil
}

func (r *Repository) GetByUserIDAndDate(userID uuid.UUID, date time.Time) (*NutritionData, error) {
	if r.db == nil {
		return nil, errors.ErrDatabase
	}
	d := &NutritionData{}
	query := `SELECT id, user_id, date, calories, protein, carbs, fats, fiber, sugar, sodium, vitamin_a, vitamin_b, vitamin_c, vitamin_d, iron, calcium, nutrition_score, created_at, updated_at FROM nutrition_data WHERE user_id = $1 AND date = $2`
	err := r.db.QueryRow(query, userID, date).Scan(&d.ID, &d.UserID, &d.Date, &d.Calories, &d.Protein, &d.Carbs, &d.Fats, &d.Fiber, &d.Sugar, &d.Sodium, &d.VitaminA, &d.VitaminB, &d.VitaminC, &d.VitaminD, &d.Iron, &d.Calcium, &d.NutritionScore, &d.CreatedAt, &d.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.ErrNotFound
		}
		return nil, errors.WrapError(err, errors.ErrDatabase)
	}
	return d, nil
}

func (r *Repository) GetByID(id uuid.UUID) (*NutritionData, error) {
	if r.db == nil {
		return nil, errors.ErrDatabase
	}
	d := &NutritionData{}
	query := `SELECT id, user_id, date, calories, protein, carbs, fats, fiber, sugar, sodium, vitamin_a, vitamin_b, vitamin_c, vitamin_d, iron, calcium, nutrition_score, created_at, updated_at FROM nutrition_data WHERE id = $1`
	err := r.db.QueryRow(query, id).Scan(&d.ID, &d.UserID, &d.Date, &d.Calories, &d.Protein, &d.Carbs, &d.Fats, &d.Fiber, &d.Sugar, &d.Sodium, &d.VitaminA, &d.VitaminB, &d.VitaminC, &d.VitaminD, &d.Iron, &d.Calcium, &d.NutritionScore, &d.CreatedAt, &d.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.ErrNotFound
		}
		return nil, errors.WrapError(err, errors.ErrDatabase)
	}
	return d, nil
}

func (r *Repository) Create(d *NutritionData) error {
	if r.db == nil {
		return errors.ErrDatabase
	}
	query := `INSERT INTO nutrition_data (id, user_id, date, calories, protein, carbs, fats, fiber, sugar, sodium, vitamin_a, vitamin_b, vitamin_c, vitamin_d, iron, calcium, nutrition_score, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19) ON CONFLICT (user_id, date) DO UPDATE SET calories=EXCLUDED.calories, protein=EXCLUDED.protein, carbs=EXCLUDED.carbs, fats=EXCLUDED.fats, fiber=EXCLUDED.fiber, sugar=EXCLUDED.sugar, sodium=EXCLUDED.sodium, vitamin_a=EXCLUDED.vitamin_a, vitamin_b=EXCLUDED.vitamin_b, vitamin_c=EXCLUDED.vitamin_c, vitamin_d=EXCLUDED.vitamin_d, iron=EXCLUDED.iron, calcium=EXCLUDED.calcium, nutrition_score=EXCLUDED.nutrition_score, updated_at=EXCLUDED.updated_at RETURNING id, user_id, date, calories, protein, carbs, fats, fiber, sugar, sodium, vitamin_a, vitamin_b, vitamin_c, vitamin_d, iron, calcium, nutrition_score, created_at, updated_at`
	now := time.Now()
	return r.db.QueryRow(query, d.ID, d.UserID, d.Date, d.Calories, d.Protein, d.Carbs, d.Fats, d.Fiber, d.Sugar, d.Sodium, d.VitaminA, d.VitaminB, d.VitaminC, d.VitaminD, d.Iron, d.Calcium, d.NutritionScore, now, now).Scan(&d.ID, &d.UserID, &d.Date, &d.Calories, &d.Protein, &d.Carbs, &d.Fats, &d.Fiber, &d.Sugar, &d.Sodium, &d.VitaminA, &d.VitaminB, &d.VitaminC, &d.VitaminD, &d.Iron, &d.Calcium, &d.NutritionScore, &d.CreatedAt, &d.UpdatedAt)
}

func (r *Repository) Update(d *NutritionData) error {
	if r.db == nil {
		return errors.ErrDatabase
	}
	query := `UPDATE nutrition_data SET calories=$1, protein=$2, carbs=$3, fats=$4, fiber=$5, sugar=$6, sodium=$7, vitamin_a=$8, vitamin_b=$9, vitamin_c=$10, vitamin_d=$11, iron=$12, calcium=$13, nutrition_score=$14, updated_at=$15 WHERE id=$16 RETURNING id, user_id, date, calories, protein, carbs, fats, fiber, sugar, sodium, vitamin_a, vitamin_b, vitamin_c, vitamin_d, iron, calcium, nutrition_score, created_at, updated_at`
	return r.db.QueryRow(query, d.Calories, d.Protein, d.Carbs, d.Fats, d.Fiber, d.Sugar, d.Sodium, d.VitaminA, d.VitaminB, d.VitaminC, d.VitaminD, d.Iron, d.Calcium, d.NutritionScore, time.Now(), d.ID).Scan(&d.ID, &d.UserID, &d.Date, &d.Calories, &d.Protein, &d.Carbs, &d.Fats, &d.Fiber, &d.Sugar, &d.Sodium, &d.VitaminA, &d.VitaminB, &d.VitaminC, &d.VitaminD, &d.Iron, &d.Calcium, &d.NutritionScore, &d.CreatedAt, &d.UpdatedAt)
}

func (r *Repository) GetStats(userID uuid.UUID, days int) (*NutritionStats, error) {
	if r.db == nil {
		return nil, errors.ErrDatabase
	}
	stats := &NutritionStats{}
	endDate := time.Now()
	startDate := endDate.AddDate(0, 0, -days)
	query := `SELECT COALESCE(AVG(calories), 0), COALESCE(AVG(protein), 0), COALESCE(AVG(carbs), 0), COALESCE(AVG(fats), 0), COUNT(*), COALESCE(AVG(nutrition_score), 0) FROM nutrition_data WHERE user_id = $1 AND date BETWEEN $2 AND $3`
	var avgScore float64
	err := r.db.QueryRow(query, userID, startDate, endDate).Scan(&stats.AvgCalories, &stats.AvgProtein, &stats.AvgCarbs, &stats.AvgFats, &stats.TotalDays, &avgScore)
	if err != nil {
		return nil, errors.WrapError(err, errors.ErrDatabase)
	}
	if avgScore > 0 {
		score := int(avgScore)
		stats.AvgNutritionScore = &score
	}
	return stats, nil
}
