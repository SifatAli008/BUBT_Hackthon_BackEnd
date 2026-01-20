package leaderboard

import (
	"github.com/google/uuid"
)

type Service struct {
	repo *Repository
}

func NewService() *Service {
	return &Service{repo: NewRepository()}
}

func (s *Service) GetLeaderboard(leaderboardType string) (*Leaderboard, error) {
	if leaderboardType == "" {
		leaderboardType = "top-sharers"
	}
	return s.repo.GetByType(leaderboardType)
}

func (s *Service) GetImpact() (*Impact, error) {
	return s.repo.GetImpact()
}

func (s *Service) GetPersonalImpact(userID uuid.UUID) (*Impact, error) {
	return s.repo.GetPersonalImpact(userID)
}
