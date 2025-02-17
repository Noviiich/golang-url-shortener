package service

import (
	"context"
	"fmt"

	"github.com/Noviiich/golang-url-shortener/internal/core/domain"
	"github.com/Noviiich/golang-url-shortener/internal/core/ports"
)

type LinkService struct {
	port  ports.LinkPort
	cache ports.Cache
}

func NewLinkService(port ports.LinkPort, cache ports.Cache) *LinkService {
	return &LinkService{
		port:  port,
		cache: cache,
	}
}

func (s *LinkService) All(ctx context.Context) ([]domain.Link, error) {
	links, err := s.port.All(ctx)
	if err != nil {
		return nil, fmt.Errorf("ошибка при получении всех коротких urls: %w", err)
	}
	return links, nil
}

func (s *LinkService) GetOriginalURL(ctx context.Context, shortID string) (*string, error) {
	linkCache, err := s.cache.Get(ctx, shortID)
	if err == nil {
		fmt.Printf("Данные взяты из кэша: %s", linkCache)
		return &linkCache, nil
	}
	link, err := s.port.Get(ctx, shortID)
	if err != nil {
		return nil, fmt.Errorf("ошибка при получении короткого url for indentifier '%s': %w", shortID, err)
	}
	s.cache.Set(ctx, link.Id, link.OriginalURL)
	return &link.OriginalURL, nil
}

func (s *LinkService) Create(ctx context.Context, link *domain.Link) error {
	_, err := s.cache.Set(ctx, link.Id, link.OriginalURL)
	if err != nil {
		return fmt.Errorf("failed to set cache for identifier '%s': %w", link.Id, err)
	}
	err = s.port.Create(ctx, link)
	if err != nil {
		return fmt.Errorf("ошибка создания короткого url: %w", err)
	}
	return nil
}

func (s *LinkService) Delete(ctx context.Context, shortID string) error {
	err := s.port.Delete(ctx, shortID)
	if err != nil {
		return fmt.Errorf("ошибка удаления короткого url for indentifier '%s': %w", shortID, err)
	}
	if err := s.cache.Delete(ctx, shortID); err != nil {
		return fmt.Errorf("failed to delete short URL for identifier '%s': %w", shortID, err)
	}
	return nil
}
