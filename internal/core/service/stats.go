package service

import (
	"context"
	"fmt"

	"github.com/Noviiich/golang-url-shortener/internal/core/domain"
	"github.com/Noviiich/golang-url-shortener/internal/core/ports"
)

type StatsService struct {
	port  ports.StatsPort
	cache ports.Cache
}

func NewStatsService(port ports.StatsPort, cache ports.Cache) *StatsService {
	return &StatsService{
		port:  port,
		cache: cache,
	}
}

func (s *StatsService) All(ctx context.Context) ([]domain.Stats, error) {
	links, err := s.port.All(ctx)
	if err != nil {
		return nil, fmt.Errorf("ошибка при получении всех коротких urls: %w", err)
	}
	return links, nil
}

func (s *StatsService) Get(ctx context.Context, linkID string) (*domain.Stats, error) {
	link, err := s.port.Get(ctx, linkID)
	if err != nil {
		return nil, fmt.Errorf("ошибка при получении короткого url for indentifier '%s': %w", linkID, err)
	}
	return link, nil
}

func (s *StatsService) Create(ctx context.Context, link *domain.Stats) error {
	err := s.port.Create(ctx, link)
	if err != nil {
		return fmt.Errorf("ошибка создания короткого url: %w", err)
	}
	return nil
}

func (s *StatsService) Delete(ctx context.Context, shortID string) error {
	err := s.port.Delete(ctx, shortID)
	if err != nil {
		return fmt.Errorf("ошибка удаления статистики по indentifier '%s': %w", shortID, err)
	}
	if err := s.cache.Delete(ctx, shortID); err != nil {
		return fmt.Errorf("failed to delete stats for identifier '%s': %w", shortID, err)
	}
	return nil
}

func (service *StatsService) GetStatsByLinkID(ctx context.Context, linkID string) ([]domain.Stats, error) {
	stats, err := service.port.GetStatsByLinkID(ctx, linkID)
	if err != nil {
		return []domain.Stats{}, fmt.Errorf("failed to get stats for identifier '%s': %w", linkID, err)
	}
	return stats, nil
}
