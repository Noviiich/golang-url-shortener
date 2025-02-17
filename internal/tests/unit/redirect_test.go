package unit

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Noviiich/golang-url-shortener/internal/adapters/handler"
	"github.com/Noviiich/golang-url-shortener/internal/core/service"
	"github.com/Noviiich/golang-url-shortener/internal/tests/mock"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestRedirectUnit(t *testing.T) {
	router := setupTestRedirect()

	tests := []struct {
		shortID          string
		expectStatusCode int
		expectLocation   string
		expectBody       string
	}{
		{
			shortID:          "testid1",
			expectStatusCode: http.StatusFound,
			expectLocation:   "https://example.com/link1",
			expectBody:       "",
		},
		{
			shortID:          "testid2",
			expectStatusCode: http.StatusFound,
			expectLocation:   "https://example.com/link2",
			expectBody:       "",
		},
		{
			shortID:          "testid3",
			expectStatusCode: http.StatusFound,
			expectLocation:   "https://example.com/link3",
			expectBody:       "",
		},
		{
			shortID:          "nonexistenid",
			expectStatusCode: 404,
			expectLocation:   "",
			expectBody:       "Link not found",
		},
	}

	for _, test := range tests {
		t.Run(test.shortID, func(t *testing.T) {
			req, err := http.NewRequest("GET", "/"+test.shortID, nil)
			assert.NoError(t, err)

			res := httptest.NewRecorder()
			router.ServeHTTP(res, req)

			assert.Equal(t, test.expectStatusCode, res.Code)
			assert.Equal(t, test.expectLocation, res.Header().Get("Location"))
		})
	}
}

func setupTestRedirect() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)

	mockLink := mock.NewMockRepository()
	mockStats := mock.NewMockStatsRepo()

	cache := mock.NewMockRedisCache()
	FillCache(cache, mockLink.Links)

	linkService := service.NewLinkService(mockLink, cache)
	statsService := service.NewStatsService(mockStats, cache)

	apiHandler := handler.NewRedirectFunctionHandler(linkService, statsService)

	router := gin.Default()
	router.GET("/:shortID", apiHandler.RedirectShortURL)

	return router
}
