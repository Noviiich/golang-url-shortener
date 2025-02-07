package main

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Noviiich/golang-url-shortener/internal/adapters/handler"
	"github.com/Noviiich/golang-url-shortener/internal/adapters/mock"
	"github.com/Noviiich/golang-url-shortener/internal/core/service"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func setupTest(body string) (*httptest.ResponseRecorder, error) {
	gin.SetMode(gin.ReleaseMode)

	mock := mock.NewMockRepository()
	service := service.NewURLService(mock)
	handler := handler.NewURLHandler(service)

	router := gin.Default()
	router.POST("/genarate", handler.CreateShortLink)

	req, err := http.NewRequest("POST", "/genarate", bytes.NewReader([]byte(body)))
	response := httptest.NewRecorder()
	router.ServeHTTP(response, req)

	return response, err
}

func TestCreateShortLink_Success(t *testing.T) {
	body := `{"long": "https://example.com/abcdefg"}`
	response, err := setupTest(body)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, response.Code)

	// var link model.Link
	// err := json.Unmarshal([]byte(w.Body.String()), &link)                     потом добавить, чтобы возвращался ответ
	// assert.NoError(t, err)
	// assert.Equal(t, "https://example.com/abcdefg", link.OriginalURL)
}

func TestCreateShortLink_EmptyString(t *testing.T) {
	body := `{"long":""}`
	response, err := setupTest(body)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, response.Code)
}

func TestCreateShoreLink_InvalidURL(t *testing.T) {
	body := `{"long": "invalid"}`
	response, err := setupTest(body)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, response.Code)
}
