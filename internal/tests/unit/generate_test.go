package unit

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Noviiich/golang-url-shortener/internal/adapters/handler"
	"github.com/Noviiich/golang-url-shortener/internal/core/service"
	"github.com/Noviiich/golang-url-shortener/internal/tests/mock"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestGenerateUnit(t *testing.T) {
	apiHandler := setupTestGenerate()

	tests := []struct {
		longURL            string
		expectedStatusCode int
		expectedBody       string
	}{
		{
			longURL:            "https://example.com/abcdefg",
			expectedStatusCode: http.StatusCreated,
			expectedBody:       "",
		},
		{
			longURL:            "",
			expectedStatusCode: http.StatusBadRequest,
			expectedBody:       "URL не должен быть пустым",
		},
		{
			longURL:            "invalid",
			expectedStatusCode: http.StatusBadRequest,
			expectedBody:       "URL должен быть больше 15 символов",
		},
	}

	for _, test := range tests {
		t.Run(test.longURL, func(t *testing.T) {
			body := `{"long": "` + test.longURL + `"}`
			req, err := http.NewRequest("POST", "/generate", strings.NewReader(body))
			assert.NoError(t, err)

			res := httptest.NewRecorder()
			apiHandler.ServeHTTP(res, req)

			assert.Equal(t, test.expectedStatusCode, res.Code)
			if test.expectedStatusCode != http.StatusCreated {
				assert.Contains(t, res.Body.String(), test.expectedBody)
			}
		})

	}
}

func setupTestGenerate() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)

	mockLink := mock.NewMockRepository()
	mockStats := mock.NewMockStatsRepo()

	cache := mock.NewMockRedisCache()
	FillCache(cache, mockLink.Links)

	linkService := service.NewLinkService(mockLink, cache)
	statsService := service.NewStatsService(mockStats, cache)

	apiHandler := handler.NewGenerateFunctionHandler(linkService, statsService)

	router := gin.Default()
	router.POST("/generate", apiHandler.CreateShortLink)

	return router
}
