package domain

import (
	"context"
	"fmt"

	"github.com/Noviiich/golang-url-shortener/types"
)

type Link struct {
	db types.DB
}

func NewLinkDomain(d types.DB) *Link {
	return &Link{db: d}
}

func (s *Link) All(ctx context.Context) ([]types.Link, error) {
	links, err := s.db.All(ctx)
	if err != nil {
		return nil, fmt.Errorf("ошибка при получении всех коротких urls: %w", err)
	}
	return links, nil
}

func (s *Link) Get(ctx context.Context, short string) (types.Link, error) {
	link, err := s.db.Get(ctx, short)
	if err != nil {
		return types.Link{}, fmt.Errorf("ошибка при получении короткого url for indentifier '%s': %w", short, err)
	}
	return link, nil
}

func (s *Link) Create(ctx context.Context, short types.Link) error {
	err := s.db.Create(ctx, short)
	if err != nil {
		return fmt.Errorf("ошибка создания короткого url: %w", err)
	}
	return nil
}

func (s *Link) Delete(ctx context.Context, short string) error {
	err := s.db.Delete(ctx, short)
	if err != nil {
		return fmt.Errorf("ошибка удаления короткого url for indentifier '%s': %w", short, err)
	}
	return nil
}
