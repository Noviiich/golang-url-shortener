package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Noviiich/golang-url-shortener/internal/adapters/handler"
	"github.com/Noviiich/golang-url-shortener/internal/adapters/mock"
	"github.com/Noviiich/golang-url-shortener/internal/core/service"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func setupTest(shortID string) (*httptest.ResponseRecorder, error) {
	gin.SetMode(gin.ReleaseMode)
	mockRep := mock.NewMockRepository()
	service := service.NewURLService(mockRep)
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
	response, err := setupTest(shortID)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusFound, response.Code)
}

func TestRedirectShortURL_NotFound(t *testing.T) {
	shortID := "nonexistentid"
	response, err := setupTest(shortID)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusNotFound, response.Code)
}
