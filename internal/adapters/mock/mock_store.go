package mock

import (
	"context"
	"errors"

	"github.com/Noviiich/golang-url-shortener/internal/core/model"
)

type MockRepository struct {
	Links []model.Link
}

func NewMockRepository() *MockRepository {
	return &MockRepository{
		Links: []model.Link{
			{ShortID: "testid1", OriginalURL: "https://example.com/link1"},
			{ShortID: "testid2", OriginalURL: "https://example.com/link2"},
			{ShortID: "testid3", OriginalURL: "https://example.com/link3"},
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
	for _, link := range m.Links {
		if link.ShortID == shortID {
			return &link, nil
		}
	}
	return &model.Link{}, errors.New("link not found")
}
func (m *MockRepository) Create(ctx context.Context, link *model.Link) error {
	for _, l := range m.Links {
		if l.ShortID == link.ShortID {
			return errors.New("link already exists")
		}
	}
	m.Links = append(m.Links, *link)
	return nil
}
func (m *MockRepository) Delete(ctx context.Context, shortID string) error {
	for i, link := range m.Links {
		if link.ShortID == shortID {
			m.Links = append(m.Links[:i], m.Links[i+1:]...)
			return nil
		}
	}
	return errors.New("link not found")
}
