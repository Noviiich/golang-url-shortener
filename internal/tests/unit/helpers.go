package unit

import (
	"context"

	"github.com/Noviiich/golang-url-shortener/internal/adapters/mock"
	"github.com/Noviiich/golang-url-shortener/internal/core/domain"
)

// func SetupTest() *handler.URLHandler {
// 	gin.SetMode(gin.ReleaseMode)

// 	mockDB := mock.NewMockRepository()
// 	cache := mock.NewMockRedisCache()
// 	fillCache(cache, mockDB.Links)
// 	service := service.NewLinkService(mockDB, cache)
// 	handler := handler.New(service)

// 	return handler
// }

func FillCache(cache *mock.MockRedisCache, links []domain.Link) error {
	for _, link := range links {
		_, err := cache.Set(context.Background(), link.Id, link.OriginalURL)
		if err != nil {
			return err
		}
	}
	return nil
}
