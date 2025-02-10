package service

import (
	"github.com/Noviiich/golang-url-shortener/internal/core/ports"
)

type StatsService struct {
	db    ports.StatsPort
	cache ports.Cache
}

func NewStatsService(db ports.StatsPort, cache ports.Cache) *StatsService {
	return &StatsService{
		db:    db,
		cache: cache,
	}
}

// func (s *StatsService) All(ctx context.Context) ([]domain.Stats, error) {
// 	links, err := s.db.All(ctx)
// 	if err != nil {
// 		return nil, fmt.Errorf("ошибка при получении всех коротких urls: %w", err)
// 	}
// 	return links, nil
// }

// func (s *StatsService) GetOriginalURL(ctx context.Context, shortID string) (*string, error) {
// 	//link, err := s.db.Get(ctx, shortID)
// 	link, err := s.cache.Get(ctx, shortID)
// 	if err != nil {
// 		return nil, fmt.Errorf("ошибка при получении короткого url for indentifier '%s': %w", shortID, err)
// 	}
// 	return &link, nil
// }

// func (s *StatsService) Create(ctx context.Context, link *domain.Stats) (string, error) {
// 	key, err := s.cache.Set(ctx, link.ShortID, link.OriginalURL)
// 	if err != nil {
// 		return "", fmt.Errorf("failed to set short URL for identifier '%s': %w", link.ShortID, err)
// 	}
// 	err = s.db.Create(ctx, link)
// 	if err != nil {
// 		return "", fmt.Errorf("ошибка создания короткого url: %w", err)
// 	}
// 	return key, nil
// }

// func (s *StatsService) Delete(ctx context.Context, shortID string) error {
// 	err := s.db.Delete(ctx, shortID)
// 	if err != nil {
// 		return fmt.Errorf("ошибка удаления короткого url for indentifier '%s': %w", shortID, err)
// 	}
// 	if err := s.cache.Delete(ctx, shortID); err != nil {
// 		return fmt.Errorf("failed to delete short URL for identifier '%s': %w", shortID, err)
// 	}
// 	return nil
// }
