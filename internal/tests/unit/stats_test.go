package unit

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Noviiich/golang-url-shortener/internal/adapters/handler"
	"github.com/Noviiich/golang-url-shortener/internal/core/domain"
	"github.com/Noviiich/golang-url-shortener/internal/core/service"
	"github.com/Noviiich/golang-url-shortener/internal/tests/mock"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestStatsUnit(t *testing.T) {
	cache := mock.NewMockRedisCache()

	mockLinkRepo := mock.NewMockRepository()
	mockStatsRepo := mock.NewMockStatsRepo()
	FillCache(cache, mockLinkRepo.Links)

	linkService := service.NewLinkService(mockLinkRepo, cache)
	statsService := service.NewStatsService(mockStatsRepo, cache)

	handler := handler.NewStatsFunctionHandler(linkService, statsService)

	router := gin.Default()
	router.GET("/stats", handler.Stats)
	t.Run("Stats unit test", func(t *testing.T) {
		req, err := http.NewRequest("GET", "/stats", nil)
		assert.NoError(t, err)

		res := httptest.NewRecorder()
		router.ServeHTTP(res, req)

		var links []domain.Link
		err = json.Unmarshal([]byte(res.Body.Bytes()), &links)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, res.Code)
		assert.Equal(t, len(links), 3)
	})
}
