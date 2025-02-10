package unit

import (
	"net/http"
	"net/http/httptest"
	"testing"

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
	handler := SetupTest()

	router := gin.Default()
	router.GET("/:shortID", handler.RedirectShortURL)

	return router
}
