package mock

import (
	"context"
	"errors"

	"github.com/Noviiich/golang-url-shortener/internal/core/domain"
)

type MockStatsRepo struct {
	Stats []domain.Stats
}

func NewMockStatsRepo() *MockStatsRepo {
	return &MockStatsRepo{Stats: MockStatsData}
}

func (m *MockStatsRepo) All(ctx context.Context) ([]domain.Stats, error) {
	return m.Stats, nil
}
func (m *MockStatsRepo) Get(ctx context.Context, id string) (*domain.Stats, error) {
	for _, stat := range m.Stats {
		if stat.Id == id {
			return &stat, nil
		}
	}
	return &domain.Stats{}, errors.New("link not found")
}
func (m *MockStatsRepo) Create(ctx context.Context, stats *domain.Stats) error {
	m.Stats = append(m.Stats, *stats)
	return nil
}
func (m *MockStatsRepo) Delete(ctx context.Context, id string) error {
	for i, stats := range m.Stats {
		if stats.Id == id {
			m.Stats = append(m.Stats[:i], m.Stats[i+1:]...)
			return nil
		}
	}
	return errors.New("link not found")
}

func (m *MockStatsRepo) GetStatsByLinkID(ctx context.Context, id string) ([]domain.Stats, error) {
	var stats []domain.Stats
	for _, stat := range m.Stats {
		if stat.Id == id {
			stats = append(stats, stat)
		}
	}

	return stats, nil
}
