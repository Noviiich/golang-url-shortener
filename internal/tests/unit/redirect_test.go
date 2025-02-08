package main

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Noviiich/golang-url-shortener/internal/adapters/handler"
	"github.com/Noviiich/golang-url-shortener/internal/adapters/mock"
	"github.com/Noviiich/golang-url-shortener/internal/core/model"
	"github.com/Noviiich/golang-url-shortener/internal/core/service"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func setupTestRedirect(shortID string) (*httptest.ResponseRecorder, error) {
	gin.SetMode(gin.ReleaseMode)
	mockDB := mock.NewMockRepository()
	cache := mock.NewMockRedisCache()
	fillCache(cache, mockDB.Links)
	service := service.NewURLService(mockDB, cache)
	handler := handler.NewURLHandler(service)

	router := gin.Default()
	router.GET("/:shortID", handler.RedirectShortURL)

	req, err := http.NewRequest("GET", "/"+shortID, nil)
	response := httptest.NewRecorder()
	router.ServeHTTP(response, req)

	return response, err
}

func TestRedirectShortURL_Success(t *testing.T) {
	shortID := "testid1"
	response, err := setupTestRedirect(shortID)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusFound, response.Code)
}

func TestRedirectShortURL_NotFound(t *testing.T) {
	shortID := "nonexistentid"
	response, err := setupTestRedirect(shortID)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusNotFound, response.Code)
}

func fillCache(cache *mock.MockRedisCache, links map[string]model.Link) error {
	for _, link := range links {
		_, err := cache.Set(context.Background(), link.ShortID, link.OriginalURL)
		if err != nil {
			return err
		}
	}
	return nil
}
