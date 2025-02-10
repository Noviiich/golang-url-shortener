package service

import (
	"context"
	"fmt"

	"github.com/Noviiich/golang-url-shortener/internal/core/domain"
	"github.com/Noviiich/golang-url-shortener/internal/core/ports"
)

type LinkService struct {
	db    ports.LinkPort
	cache ports.Cache
}

func NewLinkService(db ports.LinkPort, cache ports.Cache) *LinkService {
	return &LinkService{
		db:    db,
		cache: cache,
	}
}

func (s *LinkService) All(ctx context.Context) ([]domain.Link, error) {
	links, err := s.db.All(ctx)
	if err != nil {
		return nil, fmt.Errorf("ошибка при получении всех коротких urls: %w", err)
	}
	return links, nil
}

func (s *LinkService) GetOriginalURL(ctx context.Context, shortID string) (*string, error) {
	//link, err := s.db.Get(ctx, shortID)
	link, err := s.cache.Get(ctx, shortID)
	if err != nil {
		return nil, fmt.Errorf("ошибка при получении короткого url for indentifier '%s': %w", shortID, err)
	}
	return &link, nil
}

func (s *LinkService) Create(ctx context.Context, link *domain.Link) error {
	_, err := s.cache.Set(ctx, link.ShortID, link.OriginalURL)
	if err != nil {
		return fmt.Errorf("failed to set short URL for identifier '%s': %w", link.ShortID, err)
	}
	err = s.db.Create(ctx, link)
	if err != nil {
		return fmt.Errorf("ошибка создания короткого url: %w", err)
	}
	return nil
}

func (s *LinkService) Delete(ctx context.Context, shortID string) error {
	err := s.db.Delete(ctx, shortID)
	if err != nil {
		return fmt.Errorf("ошибка удаления короткого url for indentifier '%s': %w", shortID, err)
	}
	if err := s.cache.Delete(ctx, shortID); err != nil {
		return fmt.Errorf("failed to delete short URL for identifier '%s': %w", shortID, err)
	}
	return nil
}
