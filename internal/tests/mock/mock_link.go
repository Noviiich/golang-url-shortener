package mock

import (
	"context"
	"errors"

	"github.com/Noviiich/golang-url-shortener/internal/core/domain"
)

type MockRepository struct {
	Links []domain.Link
	Stats []domain.Stats
}

func NewMockRepository() *MockRepository {
	return &MockRepository{
		Links: MockLinkData,
		Stats: MockStatsData,
	}
}

func (m *MockRepository) All(ctx context.Context) ([]domain.Link, error) {
	var links []domain.Link
	links = append(links, m.Links...)
	return links, nil
}
func (m *MockRepository) Get(ctx context.Context, shortID string) (*domain.Link, error) {
	for _, link := range m.Links {
		if link.Id == shortID {
			return &link, nil
		}
	}
	return &domain.Link{}, errors.New("link not found")
}
func (m *MockRepository) Create(ctx context.Context, link *domain.Link) error {
	for _, l := range m.Links {
		if l.Id == link.Id {
			return errors.New("link already exists")
		}
	}
	m.Links = append(m.Links, *link)
	return nil
}
func (m *MockRepository) Delete(ctx context.Context, shortID string) error {
	for i, link := range m.Links {
		if link.Id == shortID {
			m.Links = append(m.Links[:i], m.Links[i+1:]...)
			return nil
		}
	}
	return errors.New("link not found")
}
