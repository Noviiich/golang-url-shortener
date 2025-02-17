package handler

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/Noviiich/golang-url-shortener/internal/core/domain"
	"github.com/Noviiich/golang-url-shortener/internal/core/service"
	"github.com/gin-gonic/gin"
)

type RedirectFunctionHandler struct {
	linkService  *service.LinkService
	statsService *service.StatsService
}

func NewRedirectFunctionHandler(l *service.LinkService, s *service.StatsService) *RedirectFunctionHandler {
	return &RedirectFunctionHandler{linkService: l, statsService: s}
}

func (h *RedirectFunctionHandler) RedirectShortURL(c *gin.Context) {
	shortID := c.Param("shortID")
	link, err := h.linkService.GetOriginalURL(context.Background(), shortID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	fmt.Println("handler link", *link)

	if err := h.statsService.Create(context.Background(), &domain.Stats{
		Platform:  domain.PlatformTwitter,
		Id:        shortID,
		CreatedAt: time.Now(),
	}); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Redirect(http.StatusFound, *link)
}
