package handler

import (
	"github.com/Noviiich/golang-url-shortener/internal/core/service"
)

type URLHandler struct {
	service *service.URLService
}

func NewURLHandler(s *service.URLService) *URLHandler {
	return &URLHandler{service: s}
}
