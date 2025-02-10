package unit

import (
	"context"

	"github.com/Noviiich/golang-url-shortener/internal/adapters/handler"
	"github.com/Noviiich/golang-url-shortener/internal/adapters/mock"
	"github.com/Noviiich/golang-url-shortener/internal/core/model"
	"github.com/Noviiich/golang-url-shortener/internal/core/service"
	"github.com/gin-gonic/gin"
)

func SetupTest() *handler.URLHandler {
	gin.SetMode(gin.ReleaseMode)

	mockDB := mock.NewMockRepository()
	cache := mock.NewMockRedisCache()
	fillCache(cache, mockDB.Links)
	service := service.NewURLService(mockDB, cache)
	handler := handler.NewURLHandler(service)

	return handler
}

func fillCache(cache *mock.MockRedisCache, links []model.Link) error {
	for _, link := range links {
		_, err := cache.Set(context.Background(), link.ShortID, link.OriginalURL)
		if err != nil {
			return err
		}
	}
	return nil
}
