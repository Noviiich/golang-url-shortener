package handler

import (
	"github.com/Noviiich/golang-url-shortener/internal/core/service"
)

type URLHandler struct {
	service *service.LinkService
}

func NewURLHandler(s *service.LinkService) *URLHandler {
	return &URLHandler{service: s}
}
