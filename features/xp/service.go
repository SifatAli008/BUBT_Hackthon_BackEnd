package xp

import (
	"foodlink_backend/errors"

	"github.com/google/uuid"
)

type Service struct {
	repo *Repository
}

func NewService() *Service {
	return &Service{repo: NewRepository()}
}

func (s *Service) GetByUserID(userID uuid.UUID) (*UserXP, error) {
	xp, err := s.repo.GetByUserID(userID)
	if err != nil && err == errors.ErrNotFound {
		// Create default XP record if doesn't exist
		xp = &UserXP{
			ID:             uuid.New(),
			UserID:         userID,
			TotalXP:        0,
			Level:          1,
			CurrentLevelXP: 0,
			NextLevelXP:   100,
		}
		if err := s.repo.CreateOrUpdate(xp); err != nil {
			return nil, err
		}
		return xp, nil
	}
	return xp, err
}

func (s *Service) AddXP(userID uuid.UUID, amount int) (*UserXP, error) {
	xp, err := s.GetByUserID(userID)
	if err != nil {
		return nil, err
	}
	
	xp.TotalXP += amount
	xp.CurrentLevelXP += amount
	
	// Level up logic
	for xp.CurrentLevelXP >= xp.NextLevelXP {
		xp.CurrentLevelXP -= xp.NextLevelXP
		xp.Level++
		xp.NextLevelXP = s.calculateNextLevelXP(xp.Level)
	}
	
	if err := s.repo.CreateOrUpdate(xp); err != nil {
		return nil, err
	}
	
	return xp, nil
}

func (s *Service) calculateNextLevelXP(level int) int {
	// Exponential leveling: 100 * level^1.5
	baseXP := 100.0
	nextXP := baseXP * float64(level) * 1.5
	return int(nextXP)
}

func (s *Service) GetLeaderboard(limit int) ([]*LeaderboardEntry, error) {
	return s.repo.GetLeaderboard(limit)
}
