package mock

import (
	"context"
	"errors"

	"github.com/Noviiich/golang-url-shortener/internal/core/model"
)

type MockRepository struct {
	Links map[string]model.Link
}

func NewMockRepository() *MockRepository {
	return &MockRepository{
		Links: map[string]model.Link{
			"testid1": {ShortId: "testid1", OriginalURL: "https://example.com/link1"},
			"testid2": {ShortId: "testid1", OriginalURL: "https://example.com/link1"},
			"testid3": {ShortId: "testid1", OriginalURL: "https://example.com/link1"},
		},
	}
}

func (m *MockRepository) All(ctx context.Context) ([]model.Link, error) {
	var links []model.Link
	for _, link := range m.Links {
		links = append(links, link)
	}
	return links, nil
}
func (m *MockRepository) Get(ctx context.Context, shortID string) (*model.Link, error) {
	if link, ok := m.Links[shortID]; ok {
		return &link, nil
	}
	return nil, errors.New("link not found")
}
func (m *MockRepository) Create(ctx context.Context, link *model.Link) error {
	if _, ok := m.Links[link.ShortId]; ok {
		return errors.New("link already exists")
	}
	m.Links[link.ShortId] = *link
	return nil
}
func (m *MockRepository) Delete(ctx context.Context, shortID string) error {
	if _, ok := m.Links[shortID]; !ok {
		return errors.New("link not found")
	}
	delete(m.Links, shortID)
	return nil
}
